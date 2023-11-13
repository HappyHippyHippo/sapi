package sapi

import (
	"fmt"
	"html/template"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/slate"
)

// ----------------------------------------------------------------------------
// defs
// ----------------------------------------------------------------------------

const (
	// RestContainerID defines a base id of all other rest
	// module instances registered in the application container.
	RestContainerID = slate.ContainerID + ".rest"

	// RestEndpointRegisterTag @todo doc
	RestEndpointRegisterTag = RestContainerID + ".endpoints"

	// RestAllEndpointRegistersContainerID @todo doc
	RestAllEndpointRegistersContainerID = RestEndpointRegisterTag + ".all"

	// RestProcessContainerID @todo doc.
	RestProcessContainerID = RestContainerID + ".process"

	// RestLoaderContainerID @todo doc.
	RestLoaderContainerID = RestContainerID + ".loader"

	// RestEnvID defines the rest module base environment variable name.
	RestEnvID = slate.EnvID + "_REST"
)

var (
	// RestConfigPath defines the configuration location where is
	// defined the REST service configuration.
	RestConfigPath = slate.EnvString(RestEnvID+"_CONFIG_PATH", "slate.api.rest")

	// RestWatchdogName defines the default REST service watchdog name.
	RestWatchdogName = slate.EnvString(RestEnvID+"_WATCHDOG_NAME", "rest")

	// RestPort defines the default rest service port.
	RestPort = slate.EnvInt(RestEnvID+"_PORT", 80)

	// RestLogChannel defines the default logging channel.
	RestLogChannel = slate.EnvString(RestEnvID+"_LOG_CHANNEL", "rest")

	// RestLogLevel defines the default logging level.
	RestLogLevel = slate.EnvString(RestEnvID+"_LOG_LEVEL", "info")

	// RestLogStartMessage defines the default service start logging message.
	RestLogStartMessage = slate.EnvString(RestEnvID+"_LOG_START_MESSAGE", "[service:rest] service starting ...")

	// RestLogErrorMessage defines the default service error logging message.
	RestLogErrorMessage = slate.EnvString(RestEnvID+"_LOG_ERROR_MESSAGE", "[service:rest] service error")

	// RestLogEndMessage defines the default service end logging message.
	RestLogEndMessage = slate.EnvString(RestEnvID+"_LOG_END_MESSAGE", "[service:rest] service terminated")
)

// ----------------------------------------------------------------------------
// Rest Engine
// ----------------------------------------------------------------------------

// RestEngine interface for the gin-gonic engine object.
type RestEngine interface {
	gin.IRoutes
	gin.IRouter

	Handler() http.Handler
	Delims(left string, right string) *gin.Engine
	SecureJsonPrefix(prefix string) *gin.Engine
	LoadHTMLGlob(pattern string)
	LoadHTMLFiles(files ...string)
	SetHTMLTemplate(templ *template.Template)
	SetFuncMap(funcMap template.FuncMap)
	NoRoute(handlers ...gin.HandlerFunc)
	NoMethod(handlers ...gin.HandlerFunc)
	Use(middleware ...gin.HandlerFunc) gin.IRoutes
	Routes() (routes gin.RoutesInfo)
	Run(addr ...string) (err error)
	SetTrustedProxies(trustedProxies []string) error
	RunTLS(addr string, certFile string, keyFile string) (err error)
	RunUnix(file string) (err error)
	RunFd(fd int) (err error)
	RunListener(listener net.Listener) (err error)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	HandleContext(c *gin.Context)
}

var _ RestEngine = &gin.Engine{}

// ----------------------------------------------------------------------------
// Rest Endpoint Register
// ----------------------------------------------------------------------------

// RestEndpointRegister defines an interface to an instance that
// is able to register endpoints to the REST engine/service
type RestEndpointRegister interface {
	Reg(engine RestEngine) error
}

// ----------------------------------------------------------------------------
// Rest Middleware
// ----------------------------------------------------------------------------

// RestMiddleware defines a type of data that represents
// a rest method middleware function.
type RestMiddleware func(gin.HandlerFunc) gin.HandlerFunc

// ----------------------------------------------------------------------------
// Rest Process
// ----------------------------------------------------------------------------

// RestProcess defines the REST watchdog process instance.
type RestProcess struct {
	slate.WatchdogProcess
}

var _ slate.WatchdogProcessor = &RestProcess{}

// NewRestProcess will try to instantiate an REST watchdog process.
func NewRestProcess(
	config *slate.Config,
	logger *slate.Log,
	engine RestEngine,
) (*RestProcess, error) {
	// check the config reference
	if config == nil {
		return nil, errNilPointer("cfg")
	}
	// check the log reference
	if logger == nil {
		return nil, errNilPointer("logger")
	}
	// check the engine reference
	if engine == nil {
		return nil, errNilPointer("engine")
	}
	// retrieve the rest configuration
	c, e := config.Partial(RestConfigPath, slate.ConfigPartial{})
	if e != nil {
		return nil, e
	}
	// parse the retrieved configuration
	wc := struct {
		Watchdog string
		Port     int
		Log      struct {
			Level   string
			Channel string
			Message struct {
				Start string
				Error string
				End   string
			}
		}
	}{
		Watchdog: RestWatchdogName,
		Port:     RestPort,
		Log: struct {
			Level   string
			Channel string
			Message struct {
				Start string
				Error string
				End   string
			}
		}{
			Level:   RestLogLevel,
			Channel: RestLogChannel,
			Message: struct {
				Start string
				Error string
				End   string
			}{
				Start: RestLogStartMessage,
				Error: RestLogErrorMessage,
				End:   RestLogEndMessage,
			},
		},
	}
	if _, e := c.Populate("", &wc); e != nil {
		return nil, e
	}
	// validate the logging level read from config
	logLevel, ok := slate.LogLevelMap[wc.Log.Level]
	if !ok {
		return nil, errConversion(wc.Log.Level, "log.Level")
	}
	// generate the watchdog process instance
	proc, _ := slate.NewWatchdogProcess(wc.Watchdog, func() error {
		_ = logger.Signal(wc.Log.Channel, logLevel, wc.Log.Message.Start, slate.LogContext{"port": wc.Port})

		if e := engine.Run(fmt.Sprintf(":%d", wc.Port)); e != nil {
			_ = logger.Signal(wc.Log.Channel, slate.FATAL, wc.Log.Message.Error, slate.LogContext{"error": e.Error()})
			return e
		}
		_ = logger.Signal(wc.Log.Channel, logLevel, wc.Log.Message.End)
		return nil
	})
	// return a locally defined instance of the watchdog process
	return &RestProcess{
		WatchdogProcess: *proc,
	}, nil
}

// ----------------------------------------------------------------------------
// Rest Loader
// ----------------------------------------------------------------------------

// RestLoader @todo doc
type RestLoader struct {
	engine    RestEngine
	registers []RestEndpointRegister
}

// NewRestLoader @todo doc
func NewRestLoader(
	engine RestEngine,
	registers []RestEndpointRegister,
) (*RestLoader, error) {
	// check engine argument reference
	if engine == nil {
		return nil, errNilPointer("engine")
	}
	// return the new loader instance
	return &RestLoader{
		engine:    engine,
		registers: registers,
	}, nil
}

// Load @todo doc
func (l *RestLoader) Load() error {
	// load all the registers
	for _, r := range l.registers {
		if e := r.Reg(l.engine); e != nil {
			return e
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
// Rest Service Provider
// ----------------------------------------------------------------------------

// RestServiceRegister defines the REST services provider instance.
type RestServiceRegister struct {
	slate.ServiceRegister
}

var _ slate.ServiceProvider = &RestServiceRegister{}

// NewRestServiceRegister will generate a new registry instance
func NewRestServiceRegister(
	app ...*slate.App,
) *RestServiceRegister {
	return &RestServiceRegister{
		ServiceRegister: *slate.NewServiceRegister(app...),
	}
}

// Provide will register the REST section instances in the
// application container.
func (sr RestServiceRegister) Provide(
	container *slate.ServiceContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	_ = container.Add(RestAllEndpointRegistersContainerID, sr.getEndpointRegisters(container))
	_ = container.Add(RestContainerID, func() RestEngine { return gin.New() })
	_ = container.Add(RestProcessContainerID, NewRestProcess, slate.WatchdogProcessTag)
	_ = container.Add(RestLoaderContainerID, NewRestLoader)
	return nil
}

// Boot will start the REST engine with the defined controllers.
func (sr RestServiceRegister) Boot(
	container *slate.ServiceContainer,
) (e error) {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	// execute the loader action
	loader, e := sr.getLoader(container)
	if e != nil {
		return e
	}
	return loader.Load()
}

func (RestServiceRegister) getLoader(
	container *slate.ServiceContainer,
) (*RestLoader, error) {
	// retrieve the loader entry
	entry, e := container.Get(RestLoaderContainerID)
	if e != nil {
		return nil, e
	}
	// validate the retrieved entry type
	if instance, ok := entry.(*RestLoader); ok {
		return instance, nil
	}
	return nil, errConversion(entry, "*RestLoader")
}

func (RestServiceRegister) getEndpointRegisters(
	container *slate.ServiceContainer,
) func() []RestEndpointRegister {
	return func() []RestEndpointRegister {
		// retrieve all the endpoint registers
		var registers []RestEndpointRegister
		entries, _ := container.Tag(RestEndpointRegisterTag)
		for _, entry := range entries {
			// type check the retrieved service
			register, ok := entry.(RestEndpointRegister)
			if ok {
				registers = append(registers, register)
			}
		}
		return registers
	}
}
