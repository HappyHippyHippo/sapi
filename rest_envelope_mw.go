package sapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/slate"
)

// ----------------------------------------------------------------------------
// defs
// ----------------------------------------------------------------------------

const (
	// RestEnvelopeMwContainerID defines the default id used to register
	// the application envelope middleware and related services.
	RestEnvelopeMwContainerID = RestContainerID + ".envelope.mw"

	// RestEnvelopeMwEnvID defines the envelope middleware module
	// base environment variable name.
	RestEnvelopeMwEnvID = RestEnvID + "_ENVELOPE_MW"
)

var (
	// RestEnvelopeMwConfigPathServiceID defines the config path that used to store the
	// application service identifier.
	RestEnvelopeMwConfigPathServiceID = slate.EnvString(RestEnvelopeMwEnvID+"_CONFIG_PATH_SERVICE_ID", "slate.api.rest.service.id")

	// RestEnvelopeMwConfigPathFormatAcceptList defines the config path that used toLogAcceptListErrorMessage
	// store the application accepted mime types formats.
	RestEnvelopeMwConfigPathFormatAcceptList = slate.EnvString(RestEnvelopeMwEnvID+"_CONFIG_PATH_FORMAT_ACCEPT_LIST", "slate.api.rest.accept")

	// RestEnvelopeMwConfigPathEndpointID defines the format of the configuration
	// path where the endpoint identification number can be retrieved.
	RestEnvelopeMwConfigPathEndpointID = slate.EnvString(RestEnvelopeMwEnvID+"_CONFIG_PATH_ENDPOINT_ID", "slate.api.rest.endpoints.%s.id")

	// RestEnvelopeMwLogLevel @todo doc
	RestEnvelopeMwLogLevel = slate.EnvString(RestEnvelopeMwEnvID+"_LOG_LEVEL", "error")

	// RestEnvelopeMwLogChannel @todo doc
	RestEnvelopeMwLogChannel = slate.EnvString(RestEnvelopeMwEnvID+"_LOG_CHANNEL", "rest")

	// RestEnvelopeMwLogServiceErrorMessage @todo doc
	RestEnvelopeMwLogServiceErrorMessage = slate.EnvString(RestEnvelopeMwEnvID+"_LOG_SERVICE_ERROR_MESSAGE", "Invalid service id")

	// RestEnvelopeMwLogAcceptListErrorMessage @todo doc
	RestEnvelopeMwLogAcceptListErrorMessage = slate.EnvString(RestEnvelopeMwEnvID+"_LOG_ACCEPT_LIST_ERROR_MESSAGE", "Invalid accept list")

	// RestEnvelopeMwLogEndpointErrorMessage @todo doc
	RestEnvelopeMwLogEndpointErrorMessage = slate.EnvString(RestEnvelopeMwEnvID+"_LOG_ENDPOINT_ERROR_MESSAGE", "Invalid endpoint id")

	// RestEnvelopeMwContextField @todo doc
	RestEnvelopeMwContextField = slate.EnvString(RestEnvelopeMwEnvID+"_CONTEXT_FIELD", "sapi_response")
)

// ----------------------------------------------------------------------------
// Rest Envelope Middleware Context Handlers
// ----------------------------------------------------------------------------

// RestSetResponse @todo doc
func RestSetResponse(
	ctx *gin.Context,
	response interface{},
) *gin.Context {
	ctx.Set(RestEnvelopeMwContextField, response)
	return ctx
}

func restGetResponse(
	ctx *gin.Context,
) (interface{}, bool) {
	return ctx.Get(RestEnvelopeMwContextField)
}

// ----------------------------------------------------------------------------
// Rest Envelope Middleware Generator
// ----------------------------------------------------------------------------

// RestEnvelopeMwGenerator @todo doc
type RestEnvelopeMwGenerator func(string) (RestMiddleware, error)

// NewRestEnvelopeMwGenerator returns a middleware generator function
// based on the application configuration. This middleware generator function
// should be called with the corresponding endpoint name, so it can generate
// the appropriate middleware function.
func NewRestEnvelopeMwGenerator(
	config *slate.Config,
	logger *slate.Log,
) (RestEnvelopeMwGenerator, error) {
	// check the config argument reference
	if config == nil {
		return nil, errNilPointer("config")
	}
	// check the logger argument reference
	if logger == nil {
		return nil, errNilPointer("logger")
	}
	// logging adapter
	log := func(msg string, ctx slate.LogContext) error {
		logLevel, ok := slate.LogLevelMap[RestEnvelopeMwLogLevel]
		if !ok {
			logLevel = slate.ERROR
		}
		return logger.Signal(RestEnvelopeMwLogChannel, logLevel, msg, ctx)
	}
	// retrieve the service id from the configuration
	service, e := config.Int(RestEnvelopeMwConfigPathServiceID, 0)
	if e != nil {
		_ = log(RestEnvelopeMwLogServiceErrorMessage, slate.LogContext{"error": e})
		return nil, e
	}
	// add a config observer for the service ID
	_ = config.AddObserver(RestEnvelopeMwConfigPathServiceID, func(old interface{}, new interface{}) {
		// new value type check for integer
		tnew, ok := new.(int)
		if !ok {
			_ = log(RestEnvelopeMwLogServiceErrorMessage, slate.LogContext{"value": new})
			return
		}
		service = tnew
	})
	// retrieve the service REST accepted format list
	acceptedList, e := config.List(RestEnvelopeMwConfigPathFormatAcceptList)
	if e != nil {
		_ = log(RestEnvelopeMwLogAcceptListErrorMessage, slate.LogContext{"error": e})
		return nil, e
	}
	// parse the list retrieved from the configuration
	var accepted []string
	for _, v := range acceptedList {
		if tv, ok := v.(string); ok {
			accepted = append(accepted, tv)
		}
	}
	// add a config observer for the REST accepted format list
	_ = config.AddObserver(RestEnvelopeMwConfigPathFormatAcceptList, func(old interface{}, new interface{}) {
		// new value type check for an array
		tnew, ok := new.([]interface{})
		if !ok {
			_ = log(RestEnvelopeMwLogAcceptListErrorMessage, slate.LogContext{"list": new})
			return
		}
		// iterate through all the array elements
		accepted = []string{}
		for _, v := range tnew {
			// type check for a string
			if tv, ok := v.(string); !ok {
				_ = log(RestEnvelopeMwLogAcceptListErrorMessage, slate.LogContext{"value": v})
			} else {
				// add the iterated element to the accepted format list
				accepted = append(accepted, tv)
			}
		}
	})
	// return the middleware generator
	return func(
		id string,
	) (RestMiddleware, error) {
		// retrieve the endpoint id integer value from the configuration
		configPathEndpointID := fmt.Sprintf(RestEnvelopeMwConfigPathEndpointID, id)
		endpoint, e := config.Int(configPathEndpointID, 0)
		if e != nil {
			_ = log(RestEnvelopeMwLogEndpointErrorMessage, slate.LogContext{"error": e})
			return nil, e
		}
		// add a config observer for the endpoint id integer value
		_ = config.AddObserver(configPathEndpointID, func(old interface{}, new interface{}) {
			// new value type check for integer
			tnew, ok := new.(int)
			if !ok {
				_ = log(RestEnvelopeMwLogEndpointErrorMessage, slate.LogContext{"value": new})
				return
			}
			endpoint = tnew
		})
		// return the generated middleware function
		return func(
			next gin.HandlerFunc,
		) gin.HandlerFunc {
			// return the middleware handler function
			return func(
				ctx *gin.Context,
			) {
				// declare the result parsing method
				parse := func(val interface{}) {
					var response *Envelope
					// type check the value to be enveloped
					switch v := val.(type) {
					case *Envelope:
						// just set the result as the envelope reference
						response = v
					case error:
						// set the result as a new envelope with an
						// internal server error with the given error as the
						// error message
						response =
							NewEnvelope(http.StatusInternalServerError, nil).
								AddError(NewEnvelopeStatusError(0, v.Error()))
					default:
						// set the result as a new envelope with an
						// internal server error with a generic error message
						response =
							NewEnvelope(http.StatusInternalServerError, nil).
								AddError(NewEnvelopeStatusError(0, "internal server error"))
					}
					// try to negotiate the response format with the defined
					// accepted format mime types giving the response envelope
					// as the content data of the response
					ctx.Negotiate(
						response.GetStatusCode(),
						gin.Negotiate{
							Offered: accepted,
							Data:    response.SetService(service).SetEndpoint(endpoint),
						},
					)
				}
				// always try to fallback retrieve any error to be parsed
				// and result in a proper envelope
				defer func() {
					if e := recover(); e != nil {
						parse(e)
					}
				}()
				// execute the middleware stored execution method
				next(ctx)
				// check if the response as been stored in the context to be
				// correctly parsed
				if response, exists := restGetResponse(ctx); exists {
					parse(response)
				}
			}
		}, nil
	}, nil
}

// ----------------------------------------------------------------------------
// Rest Envelope Middleware Service Register
// ----------------------------------------------------------------------------

// RestEnvelopeMwServiceRegister defines the default envelope provider to be used on
// the application initialization to register the file system adapter service.
type RestEnvelopeMwServiceRegister struct {
	slate.ServiceRegister
}

var _ slate.ServiceProvider = &RestEnvelopeMwServiceRegister{}

// Provide will add to the container a new file system adapter instance.
func (RestEnvelopeMwServiceRegister) Provide(
	container *slate.ServiceContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	_ = container.Add(RestEnvelopeMwContainerID, NewRestEnvelopeMwGenerator)
	return nil
}
