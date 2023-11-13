package sapi

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/slate"
)

// ----------------------------------------------------------------------------
// defs
// ----------------------------------------------------------------------------

const (
	// RestLogMwContainerID defines the default id used to register
	// the application log middleware and related services.
	RestLogMwContainerID = RestContainerID + ".log.mw"

	// RestLogMwEnvID defines the log middleware module base
	// environment variable name.
	RestLogMwEnvID = RestEnvID + "_LOG_MW"
)

var (
	// RestLogMwRequestChannel defines the channel id to be used when
	// the log middleware sends the request logging signal to the logger
	// instance.
	RestLogMwRequestChannel = slate.EnvString(RestLogMwEnvID+"_REQUEST_CHANNEL", "rest")

	// RestLogMwRequestLevel defines the logging level to be used when
	// the log middleware sends the request logging signal to the logger
	// instance.
	RestLogMwRequestLevel = envToLogLevel(RestLogMwEnvID+"_REQUEST_LEVEL", slate.DEBUG)

	// RestLogMwRequestMessage defines the request event logging message to
	// be used when the log middleware sends the logging signal to the logger
	// instance.
	RestLogMwRequestMessage = slate.EnvString(RestLogMwEnvID+"_REQUEST_MESSAGE", "Request")

	// RestLogMwResponseChannel defines the channel id to be used when the
	// log middleware sends the response logging signal to the logger instance.
	RestLogMwResponseChannel = slate.EnvString(RestLogMwEnvID+"_RESPONSE_CHANNEL", "rest")

	// RestLogMwResponseLevel defines the logging level to be used when the
	// log middleware sends the response logging signal to the logger instance.
	RestLogMwResponseLevel = envToLogLevel(RestLogMwEnvID+"_RESPONSE_LEVEL", slate.INFO)

	// RestLogMwResponseMessage defines the response event logging message
	// to be used when the log middleware sends the logging signal to the
	// logger instance.
	RestLogMwResponseMessage = slate.EnvString(RestLogMwEnvID+"_RESPONSE_MESSAGE", "Response")
)

func envToLogLevel(ev string, def slate.LogLevel) slate.LogLevel {
	v, ok := slate.LogLevelMap[strings.ToLower(ev)]
	if !ok {
		return def
	}
	return v
}

// ----------------------------------------------------------------------------
// Rest Log Middleware Request Reader
// ----------------------------------------------------------------------------

// RestLogMwRequestReader defines the function used by the middleware that compose the
// logging request context object.
type RestLogMwRequestReader func(ctx *gin.Context) (slate.LogContext, error)

// NewRestLogMwRequestReader is the default function used to parse the request
// context information.
func NewRestLogMwRequestReader() RestLogMwRequestReader {
	return func(
		ctx *gin.Context,
	) (slate.LogContext, error) {
		// check the context argument reference
		if ctx == nil {
			return nil, errNilPointer("ctx")
		}
		// obtain the request parameters
		params := slate.LogContext{}
		for p, v := range ctx.Request.URL.Query() {
			if len(v) == 1 {
				params[p] = v[0]
			} else {
				params[p] = v
			}
		}
		// return the default request information
		return slate.LogContext{
			"headers": getRequestHeaders(ctx.Request),
			"method":  ctx.Request.Method,
			"path":    ctx.Request.URL.Path,
			"params":  params,
			"body":    getRequestBody(ctx.Request),
		}, nil
	}
}

func getRequestHeaders(request *http.Request) slate.LogContext {
	// try to flat single entry header fields
	headers := slate.LogContext{}
	for index, header := range request.Header {
		if len(header) == 1 {
			headers[index] = header[0]
		} else {
			headers[index] = header
		}
	}
	return headers
}

func getRequestBody(request *http.Request) string {
	// obtain the request body (content destructible action)
	var bodyBytes []byte
	if request.Body != nil {
		bodyBytes, _ = io.ReadAll(request.Body)
	}
	// reassign the request body with a memory buffer
	request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}

// ----------------------------------------------------------------------------
// Rest Log Middleware Request Reader JSON Decorator
// ----------------------------------------------------------------------------

// NewRestLogMwRequestReaderJSONDecorator will instantiate a new request
// event context reader JSON decorator used to parse the request body as a JSON
// and add the parsed content into the logging data.
func NewRestLogMwRequestReaderJSONDecorator(
	reader RestLogMwRequestReader,
	model interface{},
) (RestLogMwRequestReader, error) {
	// check the reader argument reference
	if reader == nil {
		return nil, errNilPointer("reader")
	}
	// return the decorated request reader method
	return func(
		ctx *gin.Context,
	) (slate.LogContext, error) {
		// check the context argument reference
		if ctx == nil {
			return nil, errNilPointer("ctx")
		}
		// read the logging request data from the context
		data, e := reader(ctx)
		if e != nil {
			return nil, e
		}
		// try to unmarshall the request body content if the request
		// is in JSON format, and store it in the data map on the
		// bodyJson field
		contentType := strings.ToLower(ctx.Request.Header.Get("Content-Type"))
		if strings.HasPrefix(contentType, gin.MIMEJSON) {
			if e = json.Unmarshal([]byte(data["body"].(string)), &model); e == nil {
				data["bodyJson"] = model
			}
		}
		// return the request information
		return data, nil
	}, nil
}

// ----------------------------------------------------------------------------
// Rest Log Middleware Request Reader XML Decorator
// ----------------------------------------------------------------------------

// NewRestLogMwRequestReaderXMLDecorator will instantiate a new request
// event context reader XML decorator used to parse the request body as an XML
// and add the parsed content into the logging data.
func NewRestLogMwRequestReaderXMLDecorator(
	reader RestLogMwRequestReader,
	model interface{},
) (RestLogMwRequestReader, error) {
	// check the reader argument reference
	if reader == nil {
		return nil, errNilPointer("reader")
	}
	// return the decorated request reader method
	return func(
		ctx *gin.Context,
	) (slate.LogContext, error) {
		// check the context argument reference
		if ctx == nil {
			return nil, errNilPointer("ctx")
		}
		// read the logging request data from the context
		data, err := reader(ctx)
		if err != nil {
			return nil, err
		}
		// try to unmarshall the request body content if the request
		// is in XML format, and store it in the data map on the
		// bodyXml field
		contentType := strings.ToLower(ctx.Request.Header.Get("Content-Type"))
		if strings.HasPrefix(contentType, gin.MIMEXML) || strings.HasPrefix(contentType, gin.MIMEXML2) {
			if err = xml.Unmarshal([]byte(data["body"].(string)), &model); err == nil {
				data["bodyXml"] = model
			}
		}
		// return the request information
		return data, nil
	}, nil
}

// ----------------------------------------------------------------------------
// Rest Log Middleware Response Writer
// ----------------------------------------------------------------------------

type bodyHolder interface {
	Body() []byte
}

type restLogMwResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var _ gin.ResponseWriter = &restLogMwResponseWriter{}

func newRestLogMwResponseWriter(
	w gin.ResponseWriter,
) (*restLogMwResponseWriter, error) {
	// check the writer argument reference
	if w == nil {
		return nil, errNilPointer("writer")
	}
	// return a new decorated writer instance
	return &restLogMwResponseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}, nil
}

// Write executes the writing the desired bytes into the underlying writer
// and storing them in the internal buffer.
func (w restLogMwResponseWriter) Write(
	b []byte,
) (int, error) {
	// write the content in the local body copy and
	// in the default response writer
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Body will retrieve the stored bytes given on the previous calls
// to the Write method.
func (w restLogMwResponseWriter) Body() []byte {
	// get the local copy of the response body
	return w.body.Bytes()
}

// ----------------------------------------------------------------------------
// Rest Log Middleware Response Reader
// ----------------------------------------------------------------------------

// RestLogMwResponseReader defines the interface methods of a response
// context reader used to compose the data to be sent to the logger on a
// response event.
type RestLogMwResponseReader func(ctx *gin.Context, writer gin.ResponseWriter, statusCode int) (slate.LogContext, error)

// NewRestLogMwResponseReader @todo doc.
func NewRestLogMwResponseReader() RestLogMwResponseReader {
	return func(
		_ *gin.Context,
		writer gin.ResponseWriter,
		statusCode int,
	) (slate.LogContext, error) {
		// check the writer argument reference
		if writer == nil {
			return nil, errNilPointer("writer")
		}
		// obtain the response status code
		status := writer.Status()
		// store the default logging information
		data := slate.LogContext{
			"status":  status,
			"headers": getResponseHeaders(writer),
		}
		// add the response body to the logging information if the
		// response status code differs from the expected
		if status != statusCode {
			if tw, ok := writer.(bodyHolder); ok {
				data["body"] = string(tw.Body())
			}
		}
		// return the response logging information
		return data, nil
	}

}

func getResponseHeaders(
	response gin.ResponseWriter,
) slate.LogContext {
	// try to flat single entry header fields
	headers := slate.LogContext{}
	for index, header := range response.Header() {
		if len(header) == 1 {
			headers[index] = header[0]
		} else {
			headers[index] = header
		}
	}
	return headers
}

// ----------------------------------------------------------------------------
// Rest Log Middleware Response Reader JSON Decorator
// ----------------------------------------------------------------------------

// NewRestLogMwResponseReaderJSONDecorator will instantiate a new response
// event context reader JSON decorator used to parse the response   body as
// a JSON and add the parsed content into the logging data.
func NewRestLogMwResponseReaderJSONDecorator(
	reader RestLogMwResponseReader,
	model interface{},
) (RestLogMwResponseReader, error) {
	// check the reader argument reference
	if reader == nil {
		return nil, errNilPointer("reader")
	}
	// return the decorated response reader method
	return func(
		ctx *gin.Context,
		writer gin.ResponseWriter,
		statusCode int,
	) (slate.LogContext, error) {
		// check the context argument reference
		if ctx == nil {
			return nil, errNilPointer("ctx")
		}
		// check the writer argument reference
		if writer == nil {
			return nil, errNilPointer("writer")
		}
		// read the logging response data from the context
		data, err := reader(ctx, writer, statusCode)
		if err != nil {
			return nil, err
		}
		// check if there is content in the response body logging data
		// and try to unmarshall it if the response is in JSON to be logged
		// in the bodyJson field
		if body, ok := data["body"]; ok == true {
			accept := strings.ToLower(ctx.Request.Header.Get("Accept"))
			if accept == "*/*" || strings.Contains(accept, gin.MIMEJSON) {
				if err = json.Unmarshal([]byte(body.(string)), &model); err == nil {
					data["bodyJson"] = model
				}
			}
		}
		// return the response information
		return data, nil
	}, nil
}

// ----------------------------------------------------------------------------
// Rest Log Middleware Response Reader XML Decorator
// ----------------------------------------------------------------------------

// NewRestLogMwResponseReaderXMLDecorator will instantiate a new response
// event context reader XML decorator used to parse the response body as an XML
// and add the parsed content into the logging data.
func NewRestLogMwResponseReaderXMLDecorator(
	reader RestLogMwResponseReader,
	model interface{},
) (RestLogMwResponseReader, error) {
	// check the reader argument reference
	if reader == nil {
		return nil, errNilPointer("reader")
	}
	// return the decorated response reader method
	return func(
		ctx *gin.Context,
		writer gin.ResponseWriter,
		statusCode int,
	) (slate.LogContext, error) {
		// check the context argument reference
		if ctx == nil {
			return nil, errNilPointer("ctx")
		}
		// check the writer argument reference
		if writer == nil {
			return nil, errNilPointer("writer")
		}
		// read the logging response data from the context
		data, err := reader(ctx, writer, statusCode)
		if err != nil {
			return nil, err
		}
		// check if there is content in the response body logging data
		// and try to unmarshall it if the response is in XML to be logged
		// in the bodyXml field
		if body, ok := data["body"]; ok == true {
			accept := strings.ToLower(ctx.Request.Header.Get("Accept"))
			if strings.Contains(accept, gin.MIMEXML) || strings.Contains(accept, gin.MIMEXML2) {
				if err = xml.Unmarshal([]byte(body.(string)), &model); err == nil {
					data["bodyXml"] = model
				}
			}
		}
		// return the response information
		return data, nil
	}, nil
}

// ----------------------------------------------------------------------------
// Rest Log Middleware Generator
// ----------------------------------------------------------------------------

// RestLogMwGenerator @todo doc
type RestLogMwGenerator func(statusCode int) RestMiddleware

// NewRestLogMwGenerator @todo doc
func NewRestLogMwGenerator(
	logger *slate.Log,
	requestReader RestLogMwRequestReader,
	responseReader RestLogMwResponseReader,
) (RestLogMwGenerator, error) {
	// check logger argument reference
	if logger == nil {
		return nil, errNilPointer("logger")
	}
	// check request reader argument reference
	if requestReader == nil {
		return nil, errNilPointer("requestReader")
	}
	// check response reader argument reference
	if responseReader == nil {
		return nil, errNilPointer("responseReader")
	}
	// return the middleware generator function
	return func(
		statusCode int,
	) RestMiddleware {
		// return the middleware method with the expected status code
		return func(
			next gin.HandlerFunc,
		) gin.HandlerFunc {
			// return the middleware handler function
			return func(
				ctx *gin.Context,
			) {
				// override the context writer
				w, _ := newRestLogMwResponseWriter(ctx.Writer)
				ctx.Writer = w
				// obtain and log the request content
				req, _ := requestReader(ctx)
				_ = logger.Signal(
					RestLogMwRequestChannel,
					RestLogMwRequestLevel,
					RestLogMwRequestMessage,
					slate.LogContext{
						"request": req,
					},
				)
				// execute the endpoint process and calculate the elapsed
				// time of it
				startTimestamp := time.Now().UnixMilli()
				if next != nil {
					next(ctx)
				}
				duration := time.Now().UnixMilli() - startTimestamp
				// obtain and log the request, response and execution duration
				resp, _ := responseReader(ctx, w, statusCode)
				_ = logger.Signal(
					RestLogMwResponseChannel,
					RestLogMwResponseLevel,
					RestLogMwResponseMessage,
					slate.LogContext{
						"request":  req,
						"response": resp,
						"duration": duration,
					},
				)
			}
		}
	}, nil
}

// ----------------------------------------------------------------------------
// Rest Log Middleware Service Register
// ----------------------------------------------------------------------------

// RestLogMwServiceRegister defines the default envelope provider to be used on
// the application initialization to register the file system adapter service.
type RestLogMwServiceRegister struct {
	slate.ServiceRegister
}

var _ slate.ServiceProvider = &RestLogMwServiceRegister{}

// NewRestLogMwServiceRegister will generate a new registry instance
func NewRestLogMwServiceRegister(
	app ...*slate.App,
) *RestLogMwServiceRegister {
	return &RestLogMwServiceRegister{
		ServiceRegister: *slate.NewServiceRegister(app...),
	}
}

// Provide will add to the container a new file system adapter instance.
func (RestLogMwServiceRegister) Provide(
	container *slate.ServiceContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	_ = container.Add(RestLogMwContainerID, NewRestLogMwGenerator)
	return nil
}
