package sapi

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/slate"
)

func Test_envToLogLevel(t *testing.T) {
	t.Run("existing parsing", func(t *testing.T) {
		scenarios := []struct {
			input    string
			def      slate.LogLevel
			expected slate.LogLevel
		}{
			{ // fatal
				input:    "fatal",
				def:      slate.DEBUG,
				expected: slate.FATAL,
			},
			{ // FATAL
				input:    "FATAL",
				def:      slate.DEBUG,
				expected: slate.FATAL,
			},
			{ // error
				input:    "error",
				def:      slate.DEBUG,
				expected: slate.ERROR,
			},
			{ // ERROR
				input:    "ERROR",
				def:      slate.DEBUG,
				expected: slate.ERROR,
			},
			{ // warning
				input:    "warning",
				def:      slate.DEBUG,
				expected: slate.WARNING,
			},
			{ // WARNING
				input:    "WARNING",
				def:      slate.DEBUG,
				expected: slate.WARNING,
			},
			{ // notice
				input:    "notice",
				def:      slate.DEBUG,
				expected: slate.NOTICE,
			},
			{ // NOTICE
				input:    "NOTICE",
				def:      slate.DEBUG,
				expected: slate.NOTICE,
			},
			{ // info
				input:    "info",
				def:      slate.DEBUG,
				expected: slate.INFO,
			},
			{ // INFO
				input:    "INFO",
				def:      slate.DEBUG,
				expected: slate.INFO,
			},
			{ // debug
				input:    "debug",
				def:      slate.DEBUG,
				expected: slate.DEBUG,
			},
			{ // DEBUG
				input:    "DEBUG",
				def:      slate.DEBUG,
				expected: slate.DEBUG,
			},
			{ // unknown -> return default
				input:    "unknown",
				def:      slate.INFO,
				expected: slate.INFO,
			},
		}

		for _, scenario := range scenarios {
			if chk := envToLogLevel(scenario.input, scenario.def); chk != scenario.expected {
				t.Errorf("parsed to  (%v) when expecting (%v)", chk, scenario.expected)
			}
		}
	})
}

func Test_RestLogMwRequestReader(t *testing.T) {
	t.Run("NewRestLogMwRequestReader", func(t *testing.T) {
		t.Run("nil writer", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if _, e := NewRestLogMwRequestReader()(nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("valid request", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			method := "method"
			path := "/resource"
			uri := path + "?param1=value1&param2=value2&param2=value3"
			reqURL, _ := url.Parse("http://domain" + uri)
			headers := map[string][]string{"header1": {"value1", "value2"}, "header2": {"value3"}}
			expHeaders := slate.LogContext{"header1": []string{"value1", "value2"}, "header2": "value3"}
			expParams := slate.LogContext{
				"param1": "value1",
				"param2": []string{"value2", "value3"},
			}
			jsonBody := map[string]interface{}{"field": "value"}
			rawBody, _ := json.Marshal(jsonBody)

			body := NewMockReader(ctrl)
			gomock.InOrder(
				body.
					EXPECT().
					Read(gomock.Any()).
					DoAndReturn(func(p []byte) (int, error) {
						copy(p, rawBody)
						return len(rawBody), nil
					}),
				body.EXPECT().Read(gomock.Any()).Return(0, io.EOF),
			)

			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Method = method
			ctx.Request.URL = reqURL
			ctx.Request.Header = headers
			ctx.Request.Body = body

			data, _ := NewRestLogMwRequestReader()(ctx)

			t.Run("retrieve the request method", func(t *testing.T) {
				if value := data["method"]; value != method {
					t.Errorf("stored the (%s) method value when expecting (%v)", value, method)
				}
			})

			t.Run("retrieve the request path", func(t *testing.T) {
				if value := data["path"]; value != path {
					t.Errorf("stored the (%s) path value when expecting (%v)", value, path)
				}
			})

			t.Run("retrieve the request params", func(t *testing.T) {
				if value := data["params"]; !reflect.DeepEqual(value, expParams) {
					t.Errorf("stored the (%s) params value when expecting (%v)", value, expParams)
				}
			})

			t.Run("retrieve the request headers", func(t *testing.T) {
				if value := data["headers"]; !reflect.DeepEqual(value, expHeaders) {
					t.Errorf("stored the (%v) headers when expecting (%v)", value, expHeaders)
				}
			})

			t.Run("retrieve the request body", func(t *testing.T) {
				if value := data["body"]; !reflect.DeepEqual(value, string(rawBody)) {
					t.Errorf("stored the (%v) body when expecting (%v)", value, string(rawBody))
				}
			})
		})
	})
}

func Test_RestLogMwRequestReaderJSONDecorator(t *testing.T) {
	t.Run("NewRestLogMwRequestReaderJSONDecorator", func(t *testing.T) {
		t.Run("nil reader", func(t *testing.T) {
			if _, e := NewRestLogMwRequestReaderJSONDecorator(nil, nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("nil context", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reader := func(_ *gin.Context) (slate.LogContext, error) { return nil, nil }
			decorator, _ := NewRestLogMwRequestReaderJSONDecorator(reader, nil)

			result, e := decorator(nil)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("base reader error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := fmt.Errorf("error message")
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			reader := func(_ *gin.Context) (slate.LogContext, error) { return nil, expected }
			decorator, _ := NewRestLogMwRequestReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !reflect.DeepEqual(e, expected):
				t.Errorf("(%v) when expecting (%v)", e, expected)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("empty content-type does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{"body": `{"field":"value"}`}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			reader := func(_ *gin.Context) (slate.LogContext, error) { return data, nil }
			decorator, _ := NewRestLogMwRequestReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyJson"]; ok {
					t.Error("added the bodyJson field")
				}
			}
		})

		t.Run("non-json content-type does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{"body": `{"field":"value"}`}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Content-Type", gin.MIMEXML)
			reader := func(_ *gin.Context) (slate.LogContext, error) { return data, nil }
			decorator, _ := NewRestLogMwRequestReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyJson"]; ok {
					t.Error("added the bodyJson field")
				}
			}
		})

		t.Run("invalid json content does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{"body": "field"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Content-Type", gin.MIMEJSON)
			ctx.Request.Header.Add("Content-Type", gin.MIMEXML)
			reader := func(_ *gin.Context) (slate.LogContext, error) { return data, nil }
			decorator, _ := NewRestLogMwRequestReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyJson"]; ok {
					t.Error("added the bodyJson field")
				}
			}
		})

		t.Run("correctly add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{"body": `{"field":"value"}`}
			expected := map[string]interface{}{"field": "value"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Content-Type", gin.MIMEJSON)
			reader := func(_ *gin.Context) (slate.LogContext, error) { return data, nil }
			decorator, _ := NewRestLogMwRequestReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if body, ok := result["bodyJson"]; !ok {
					t.Error("didn't added the bodyJson field")
				} else if !reflect.DeepEqual(body, expected) {
					t.Errorf("(%v) when expecting (%v)", body, expected)
				}
			}
		})
	})
}

func Test_RestLogMwRequestReaderXMLDecorator(t *testing.T) {
	t.Run("NewRestLogMwRequestReaderXMLDecorator", func(t *testing.T) {
		t.Run("nil reader", func(t *testing.T) {
			if _, e := NewRestLogMwRequestReaderXMLDecorator(nil, nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})
		t.Run("nil context", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reader := func(ctx *gin.Context) (slate.LogContext, error) { return nil, nil }
			decorator, _ := NewRestLogMwRequestReaderXMLDecorator(reader, nil)

			result, e := decorator(nil)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("base reader error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := fmt.Errorf("error message")
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			reader := func(ctx *gin.Context) (slate.LogContext, error) { return nil, expected }
			decorator, _ := NewRestLogMwRequestReaderXMLDecorator(reader, nil)

			result, e := decorator(ctx)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !reflect.DeepEqual(e, expected):
				t.Errorf("(%v) when expecting (%v)", e, expected)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("empty content-type does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message><field>value</field></message>"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			reader := func(ctx *gin.Context) (slate.LogContext, error) { return data, nil }
			decorator, _ := NewRestLogMwRequestReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyXml"]; ok {
					t.Error("added the bodyXml field")
				}
			}
		})

		t.Run("non-xml content-type does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message><field>value</field></message>"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Content-Type", gin.MIMEJSON)
			reader := func(ctx *gin.Context) (slate.LogContext, error) { return data, nil }
			decorator, _ := NewRestLogMwRequestReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyXml"]; ok {
					t.Error("added the bodyXml field")
				}
			}
		})

		t.Run("invalid xml content does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message field value /field /message>"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Content-Type", gin.MIMEXML)
			ctx.Request.Header.Add("Content-Type", gin.MIMEXML)
			reader := func(ctx *gin.Context) (slate.LogContext, error) { return data, nil }
			decorator, _ := NewRestLogMwRequestReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyXml"]; ok {
					t.Error("added the bodyXml field")
				}
			}
		})

		t.Run("correctly add decorated field for application/xml", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message><field>value</field></message>"}
			expected := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{XMLName: xml.Name{Local: "message"}, Field: "value"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Content-Type", gin.MIMEXML)
			reader := func(ctx *gin.Context) (slate.LogContext, error) { return data, nil }
			decorator, _ := NewRestLogMwRequestReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if body, ok := result["bodyXml"]; !ok {
					t.Error("didn't added the bodyXml field")
				} else if !reflect.DeepEqual(body, &expected) {
					t.Errorf("(%v) when expecting (%v)", body, &expected)
				}
			}
		})

		t.Run("correctly add decorated field for text/xml", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message><field>value</field></message>"}
			expected := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{XMLName: xml.Name{Local: "message"}, Field: "value"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Content-Type", gin.MIMEXML2)
			reader := func(ctx *gin.Context) (slate.LogContext, error) { return data, nil }
			decorator, _ := NewRestLogMwRequestReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if body, ok := result["bodyXml"]; !ok {
					t.Error("didn't added the bodyXml field")
				} else if !reflect.DeepEqual(body, &expected) {
					t.Errorf("(%v) when expecting (%v)", body, &expected)
				}
			}
		})
	})
}

func Test_restLogMwResponseWriter(t *testing.T) {
	t.Run("newRestLogMwResponseWriter", func(t *testing.T) {
		t.Run("error when missing writer", func(t *testing.T) {
			if _, e := newRestLogMwResponseWriter(nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("new log response writer", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			writer := NewMockResponseWriter(ctrl)

			if value, e := newRestLogMwResponseWriter(writer); e != nil {
				t.Errorf("return the (%v) error", e)
			} else if value == nil {
				t.Error("didn't returned a valid reference")
			}
		})
	})

	t.Run("Write", func(t *testing.T) {
		t.Run("write to buffer and underlying writer", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			b := []byte{12, 34, 56}
			ginWriter := NewMockResponseWriter(ctrl)
			ginWriter.EXPECT().Write(b).Times(1)
			writer, _ := newRestLogMwResponseWriter(ginWriter)
			writer.body = &bytes.Buffer{}
			_, _ = writer.Write(b)

			if !reflect.DeepEqual(writer.body.Bytes(), b) {
				t.Errorf("written (%v) bytes on buffer", writer.body)
			}
		})
	})

	t.Run("body", func(t *testing.T) {
		t.Run("write to buffer and underlying writer", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			b := []byte{12, 34, 56}
			ginWriter := NewMockResponseWriter(ctrl)
			writer, _ := newRestLogMwResponseWriter(ginWriter)
			writer.body = bytes.NewBuffer(b)

			if !reflect.DeepEqual(writer.Body(), b) {
				t.Errorf("written (%v) bytes on buffer", writer.body)
			}
		})
	})
}

func Test_RestLogMwResponseReader(t *testing.T) {
	t.Run("NewRestLogMwResponseReader", func(t *testing.T) {
		t.Run("nil writer", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			statusCode := 200
			if _, e := NewRestLogMwResponseReader()(nil, nil, statusCode); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("don't store the body if status code is the expected", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			statusCode := 200
			headers := map[string][]string{"header1": {"value1", "value2"}, "header2": {"value3"}}
			expHeaders := slate.LogContext{"header1": []string{"value1", "value2"}, "header2": "value3"}
			writer := NewMockResponseWriter(ctrl)
			writer.EXPECT().Status().Return(statusCode).Times(1)
			writer.EXPECT().Header().Return(headers).Times(1)

			if data, e := NewRestLogMwResponseReader()(nil, writer, statusCode); e != nil {
				t.Errorf("returned the unextected (%v) error", e)
			} else if value := data["status"]; value != statusCode {
				t.Errorf("stored the (%s) status value", value)
			} else if value := data["headers"]; !reflect.DeepEqual(value, expHeaders) {
				t.Errorf("stored the (%v) headers", value)
			} else if value, exists := data["body"]; exists {
				t.Errorf("stored the (%v) body", value)
			}
		})

		t.Run("store the body if status code is different then the expected", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			statusCode := 200
			headers := map[string][]string{"header1": {"value1", "value2"}, "header2": {"value3"}}
			expHeaders := slate.LogContext{"header1": []string{"value1", "value2"}, "header2": "value3"}
			jsonBody := map[string]interface{}{"field": "value"}
			rawBody, _ := json.Marshal(jsonBody)
			writer := NewMockResponseWriter(ctrl)
			writer.EXPECT().Status().Return(statusCode).Times(1)
			writer.EXPECT().Header().Return(headers).Times(1)
			writer.EXPECT().Body().Return(rawBody).Times(1)

			if data, e := NewRestLogMwResponseReader()(nil, writer, statusCode+1); e != nil {
				t.Errorf("returned the unextected (%v) error", e)
			} else if value := data["status"]; value != statusCode {
				t.Errorf("stored the (%s) status value", value)
			} else if value := data["headers"]; !reflect.DeepEqual(value, expHeaders) {
				t.Errorf("stored the (%v) headers", value)
			} else if value := data["body"]; !reflect.DeepEqual(value, string(rawBody)) {
				t.Errorf("stored the (%v) body", value)
			}
		})
	})
}

func Test_RestLogMwResponseReaderJSONDecorator(t *testing.T) {
	t.Run("NewRestLogMwResponseReaderJSONDecorator", func(t *testing.T) {
		t.Run("nil reader", func(t *testing.T) {
			if _, e := NewRestLogMwResponseReaderJSONDecorator(nil, nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("nil context", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return nil, nil
			}
			decorator, _ := NewRestLogMwResponseReaderJSONDecorator(reader, nil)

			result, e := decorator(nil, writer, 0)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("nil writer", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := &gin.Context{}
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return nil, nil
			}
			decorator, _ := NewRestLogMwResponseReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx, nil, 0)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("base reader error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := fmt.Errorf("error message")
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return nil, expected
			}
			decorator, _ := NewRestLogMwResponseReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !reflect.DeepEqual(e, expected):
				t.Errorf("(%v) when expecting (%v)", e, expected)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("missing body does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyJson"]; ok {
					t.Error("added the bodyJson field")
				}
			}
		})

		t.Run("empty accept does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{"body": `{"field":"value"}`}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyJson"]; ok {
					t.Error("added the bodyJson field")
				}
			}
		})

		t.Run("non-json accept does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{"body": `{"field":"value"}`}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Accept", gin.MIMEXML)
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyJson"]; ok {
					t.Error("added the bodyJson field")
				}
			}
		})

		t.Run("invalid json content does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{"body": "{field value}"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Accept", gin.MIMEJSON)
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyJson"]; ok {
					t.Error("added the bodyJson field")
				}
			}
		})

		t.Run("correctly add decorated field for application/json", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{"body": `{"field":"value"}`}
			expected := map[string]interface{}{"field": "value"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Accept", gin.MIMEJSON)
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if body, ok := result["bodyJson"]; !ok {
					t.Error("didn't added the bodyJson field")
				} else if !reflect.DeepEqual(body, expected) {
					t.Errorf("(%v) when expecting (%v)", body, expected)
				}
			}
		})

		t.Run("correctly add decorated field for 'any mime type'", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{"body": `{"field":"value"}`}
			expected := map[string]interface{}{"field": "value"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Accept", "*/*")
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderJSONDecorator(reader, nil)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if body, ok := result["bodyJson"]; !ok {
					t.Error("didn't added the bodyJson field")
				} else if !reflect.DeepEqual(body, expected) {
					t.Errorf("(%v) when expecting (%v)", body, expected)
				}
			}
		})
	})
}

func Test_RestLogMwResponseReaderXMLDecorator(t *testing.T) {
	t.Run("NewRestLogMwResponseReaderXMLDecorator", func(t *testing.T) {
		t.Run("nil reader", func(t *testing.T) {
			if _, e := NewRestLogMwResponseReaderXMLDecorator(nil, nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("nil context", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return nil, nil
			}
			decorator, _ := NewRestLogMwResponseReaderXMLDecorator(reader, nil)

			result, e := decorator(nil, writer, 0)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("nil writer", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := &gin.Context{}
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return nil, nil
			}
			decorator, _ := NewRestLogMwResponseReaderXMLDecorator(reader, nil)

			result, e := decorator(ctx, nil, 0)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("base reader error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := fmt.Errorf("error message")
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return nil, expected
			}
			decorator, _ := NewRestLogMwResponseReaderXMLDecorator(reader, nil)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !reflect.DeepEqual(e, expected):
				t.Errorf("(%v) when expecting (%v)", e, expected)
			case result != nil:
				t.Errorf("returned the unexpeted context data : %v", result)
			}
		})

		t.Run("missing body does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := slate.LogContext{}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderXMLDecorator(reader, nil)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyXml"]; ok {
					t.Error("added the bodyXml field")
				}
			}
		})

		t.Run("empty accept does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message><field>value</field></message>"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyXml"]; ok {
					t.Error("added the bodyXml field")
				}
			}
		})

		t.Run("non-xml accept does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message><field>value</field></message>"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Accept", gin.MIMEJSON)
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyXml"]; ok {
					t.Error("added the bodyXml field")
				}
			}
		})

		t.Run("invalid xml content does not add decorated field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message> field value /field /message>"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Accept", gin.MIMEXML)
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if _, ok := result["bodyXml"]; ok {
					t.Error("added the bodyXml field")
				}
			}
		})

		t.Run("correctly add decorated field for application/xml", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message><field>value</field></message>"}
			expected := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{XMLName: xml.Name{Local: "message"}, Field: "value"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Accept", gin.MIMEXML)
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if body, ok := result["bodyXml"]; !ok {
					t.Error("didn't added the bodyXml field")
				} else if !reflect.DeepEqual(body, &expected) {
					t.Errorf("(%v) when expecting (%v)", body, &expected)
				}
			}
		})

		t.Run("correctly add decorated field for text/xml", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			model := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{}
			data := slate.LogContext{"body": "<message><field>value</field></message>"}
			expected := struct {
				XMLName xml.Name `xml:"message"`
				Field   string   `xml:"field"`
			}{XMLName: xml.Name{Local: "message"}, Field: "value"}
			ctx := &gin.Context{}
			ctx.Request = &http.Request{}
			ctx.Request.Header = http.Header{}
			ctx.Request.Header.Add("Accept", gin.MIMEXML2)
			writer := NewMockResponseWriter(ctrl)
			reader := func(_ *gin.Context, _ gin.ResponseWriter, _ int) (slate.LogContext, error) {
				return data, nil
			}
			decorator, _ := NewRestLogMwResponseReaderXMLDecorator(reader, &model)

			result, e := decorator(ctx, writer, 0)
			switch {
			case e != nil:
				t.Errorf("returned the unexpected (%v) error", e)
			case result == nil:
				t.Error("didn't returned the expected context data")
			default:
				if body, ok := result["bodyXml"]; !ok {
					t.Error("didn't added the bodyXml field")
				} else if !reflect.DeepEqual(body, &expected) {
					t.Errorf("(%v) when expecting (%v)", body, &expected)
				}
			}
		})
	})
}

func Test_RestLogMwGenerator(t *testing.T) {
	t.Run("NewRestLogMwGenerator", func(t *testing.T) {
		t.Run("nil logger", func(t *testing.T) {
			reqReader := func(ctx *gin.Context) (slate.LogContext, error) { return nil, nil }
			resReader := func(ctx *gin.Context, writer gin.ResponseWriter, statusCode int) (slate.LogContext, error) {
				return nil, nil
			}

			generator, e := NewRestLogMwGenerator(nil, reqReader, resReader)
			switch {
			case e == nil:
				t.Errorf("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			case generator != nil:
				t.Error("unexpected valid middleware generator reference")
			}
		})

		t.Run("nil request reader", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			resReader := func(ctx *gin.Context, writer gin.ResponseWriter, statusCode int) (slate.LogContext, error) {
				return nil, nil
			}

			generator, e := NewRestLogMwGenerator(slate.NewLog(), nil, resReader)
			switch {
			case e == nil:
				t.Errorf("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			case generator != nil:
				t.Error("unexpected valid middleware generator reference")
			}
		})

		t.Run("nil response reader", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reqReader := func(ctx *gin.Context) (slate.LogContext, error) { return nil, nil }

			generator, e := NewRestLogMwGenerator(slate.NewLog(), reqReader, nil)
			switch {
			case e == nil:
				t.Errorf("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			case generator != nil:
				t.Error("unexpected valid middleware generator reference")
			}
		})

		t.Run("valid middleware generator", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reqReader := func(ctx *gin.Context) (slate.LogContext, error) { return nil, nil }
			resReader := func(ctx *gin.Context, writer gin.ResponseWriter, statusCode int) (slate.LogContext, error) {
				return nil, nil
			}

			generator, e := NewRestLogMwGenerator(slate.NewLog(), reqReader, resReader)
			switch {
			case e != nil:
				t.Errorf("unexpected (%v) error", e)
			case generator == nil:
				t.Error("didn't returned the expected middleware generator reference")
			}
		})

		t.Run("correctly call next handler", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			RestLogMwRequestChannel = "channel.request"
			RestLogMwRequestLevel = slate.WARNING
			RestLogMwResponseChannel = "channel.response"
			RestLogMwResponseLevel = slate.ERROR
			defer func() {
				RestLogMwRequestChannel = "Request"
				RestLogMwRequestLevel = slate.DEBUG
				RestLogMwResponseChannel = "Response"
				RestLogMwResponseLevel = slate.INFO
			}()

			statusCode := 123
			writer := NewMockResponseWriter(ctrl)
			ctx := &gin.Context{}
			ctx.Writer = writer
			callCount := 0
			var next gin.HandlerFunc = func(context *gin.Context) {
				if context != ctx {
					t.Errorf("handler called with unexpected context instance")
					return
				}
				callCount++
			}
			req := slate.LogContext{"type": "request"}
			res := slate.LogContext{"type": "response"}
			logWriter := NewMockLogWriter(ctrl)
			gomock.InOrder(
				logWriter.
					EXPECT().
					Signal(
						RestLogMwRequestChannel,
						RestLogMwRequestLevel,
						RestLogMwRequestMessage,
						slate.LogContext{"request": req},
					),
				logWriter.
					EXPECT().
					Signal(
						RestLogMwResponseChannel,
						RestLogMwResponseLevel,
						RestLogMwResponseMessage,
						slate.LogContext{"request": req, "response": res, "duration": int64(0)},
					),
			)
			logger := slate.NewLog()
			_ = logger.AddWriter("id", logWriter)
			requestReader := func(context *gin.Context) (slate.LogContext, error) {
				if context != ctx {
					t.Errorf("handler called with unexpected context instance")
				}
				return req, nil
			}
			responseReader := func(context *gin.Context, _ gin.ResponseWriter, sc int) (slate.LogContext, error) {
				if context != ctx {
					t.Errorf("handler called with unexpected context instance")
				}
				if sc != statusCode {
					t.Errorf("handler called with unexpected status code")
				}
				return res, nil
			}

			generator, _ := NewRestLogMwGenerator(logger, requestReader, responseReader)
			mw := generator(statusCode)
			handler := mw(next)
			handler(ctx)

			if callCount != 1 {
				t.Errorf("didn't called the next handler")
			}
		})
	})
}

func Test_RestLogMwServiceRegister(t *testing.T) {
	t.Run("NewRestLogMwServiceRegister", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			if NewRestLogMwServiceRegister() == nil {
				t.Error("didn't returned a valid reference")
			}
		})

		t.Run("create with app reference", func(t *testing.T) {
			app := slate.NewApp()
			if sut := NewRestLogMwServiceRegister(app); sut == nil {
				t.Error("didn't returned a valid reference")
			} else if sut.App != app {
				t.Error("didn't stored the app reference")
			}
		})
	})

	t.Run("Provide", func(t *testing.T) {
		t.Run("nil container", func(t *testing.T) {
			if e := NewRestLogMwServiceRegister().Provide(nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expected (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("register components", func(t *testing.T) {
			container := slate.NewServiceContainer()
			sut := NewRestLogMwServiceRegister()

			e := sut.Provide(container)
			switch {
			case e != nil:
				t.Errorf("unexpected (%v) error", e)
			case !container.Has(RestLogMwContainerID):
				t.Errorf("no log middleware generator : %v", sut)
			}
		})

		t.Run("retrieving the log middleware generator", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := slate.NewServiceContainer()
			_ = slate.NewFileSystemServiceRegister().Provide(container)
			_ = slate.NewConfigServiceRegister().Provide(container)
			_ = slate.NewLogServiceRegister().Provide(container)
			_ = NewRestLogMwServiceRegister().Provide(container)

			_ = container.Add("request.reader", NewRestLogMwRequestReader)
			_ = container.Add("response.reader", NewRestLogMwResponseReader)

			sut, e := container.Get(RestLogMwContainerID)
			switch {
			case e != nil:
				t.Errorf("unexpected error (%v)", e)
			case sut == nil:
				t.Error("didn't returned a reference to service")
			default:
				switch sut.(type) {
				case RestLogMwGenerator:
				default:
					t.Error("didn't returned the log middleware generator")
				}
			}
		})
	})
}
