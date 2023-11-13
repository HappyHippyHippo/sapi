package sapi

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/slate"
)

func Test_RestSetResponse(t *testing.T) {
	t.Run("store response data", func(t *testing.T) {
		ctx := &gin.Context{}
		data := "string message"
		if chk := RestSetResponse(ctx, data); chk != ctx {
			t.Errorf("didn't returned the passed context")
		} else if msg, ok := ctx.Get(RestEnvelopeMwContextField); !ok {
			t.Errorf("didn't stored the data in the expected context variable")
		} else if msg != data {
			t.Errorf("didn't stored the correct data")
		}
	})
}

func Test_RestGetResponse(t *testing.T) {
	t.Run("store response data", func(t *testing.T) {
		ctx := &gin.Context{}
		data := "string message"
		_ = RestSetResponse(ctx, data)

		if msg, ok := restGetResponse(ctx); !ok {
			t.Errorf("didn't stored the data in the expected context variable")
		} else if msg != data {
			t.Errorf("didn't stored the correct data")
		}
	})
}

func Test_NewRestEnvelopeMwGenerator(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		generator, e := NewRestEnvelopeMwGenerator(nil, slate.NewLog())
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrNilPointer):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
		}
	})

	t.Run("nil logger", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		generator, e := NewRestEnvelopeMwGenerator(slate.NewConfig(), nil)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrNilPointer):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
		}
	})

	t.Run("error getting the service id from config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				RestEnvelopeMwLogChannel,
				slate.ERROR,
				RestEnvelopeMwLogServiceErrorMessage,
				slate.LogContext{"error": errConversion("invalid", "int")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("default to error level logging on invalid log level", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogChannel
		RestEnvelopeMwLogChannel = "invalid"
		defer func() { RestEnvelopeMwLogChannel = prev }()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				RestEnvelopeMwLogChannel,
				slate.ERROR,
				RestEnvelopeMwLogServiceErrorMessage,
				slate.LogContext{"error": errConversion("invalid", "int")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log for environment defined channel when getting the service id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogChannel
		RestEnvelopeMwLogChannel = "new channel"
		defer func() { RestEnvelopeMwLogChannel = prev }()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				RestEnvelopeMwLogChannel,
				slate.ERROR,
				RestEnvelopeMwLogServiceErrorMessage,
				slate.LogContext{"error": errConversion("invalid", "int")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log with environment defined level when getting the service id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogLevel
		RestEnvelopeMwLogLevel = "warning"
		defer func() { RestEnvelopeMwLogLevel = prev }()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				RestEnvelopeMwLogChannel,
				slate.WARNING,
				RestEnvelopeMwLogServiceErrorMessage,
				slate.LogContext{"error": errConversion("invalid", "int")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log with error level when unrecognizable level while getting the service id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogLevel
		RestEnvelopeMwLogLevel = "unknown"
		defer func() { RestEnvelopeMwLogLevel = prev }()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				RestEnvelopeMwLogChannel,
				slate.ERROR,
				RestEnvelopeMwLogServiceErrorMessage,
				slate.LogContext{"error": errConversion("invalid", "int")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log with environment defined message when getting the service id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogServiceErrorMessage
		RestEnvelopeMwLogServiceErrorMessage = "test"
		defer func() { RestEnvelopeMwLogServiceErrorMessage = prev }()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				RestEnvelopeMwLogChannel,
				slate.ERROR,
				"test",
				slate.LogContext{"error": errConversion("invalid", "int")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("error getting the service accept list from config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				RestEnvelopeMwLogChannel,
				slate.ERROR,
				RestEnvelopeMwLogAcceptListErrorMessage,
				slate.LogContext{"error": errConversion("invalid", "[]interface{}")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log for environment defined channel when getting the service accept list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogChannel
		RestEnvelopeMwLogChannel = "test"
		defer func() { RestEnvelopeMwLogChannel = prev }()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				"test",
				slate.ERROR,
				RestEnvelopeMwLogAcceptListErrorMessage,
				slate.LogContext{"error": errConversion("invalid", "[]interface{}")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log with environment defined level when getting the service accept list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogLevel
		RestEnvelopeMwLogLevel = "warning"
		defer func() { RestEnvelopeMwLogLevel = prev }()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				RestEnvelopeMwLogChannel,
				slate.WARNING,
				RestEnvelopeMwLogAcceptListErrorMessage,
				slate.LogContext{"error": errConversion("invalid", "[]interface{}")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log with environment defined message when getting the service accept list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogAcceptListErrorMessage
		RestEnvelopeMwLogAcceptListErrorMessage = "test"
		defer func() { RestEnvelopeMwLogAcceptListErrorMessage = prev }()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal(
				RestEnvelopeMwLogChannel,
				slate.ERROR,
				"test",
				slate.LogContext{"error": errConversion("invalid", "[]interface{}")}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator != nil:
			t.Error("returned a valid reference")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("valid generator instantiation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logger := slate.NewLog()

		generator, e := NewRestEnvelopeMwGenerator(config, logger)
		switch {
		case generator == nil:
			t.Error("didn't returned a valid reference")
		case e != nil:
			t.Errorf("unexpected (%v) error", e)
		}
	})

	t.Run("error while retrieving endpoint path when generating middleware", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", "string")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "Invalid endpoint id", gomock.Any()).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)

		mw, e := generator(endpoint)
		switch {
		case mw != nil:
			t.Error("unexpected valid reference to a middleware")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log for environment defined channel when retrieving endpoint path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogChannel
		RestEnvelopeMwLogChannel = "test"
		defer func() { RestEnvelopeMwLogChannel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", "string")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("test", slate.ERROR, "Invalid endpoint id", gomock.Any()).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)

		mw, e := generator(endpoint)
		switch {
		case mw != nil:
			t.Error("unexpected valid reference to a middleware")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log for environment level channel when retrieving endpoint path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogLevel
		RestEnvelopeMwLogLevel = "debug"
		defer func() { RestEnvelopeMwLogLevel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", "string")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.DEBUG, "Invalid endpoint id", gomock.Any()).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)

		mw, e := generator(endpoint)
		switch {
		case mw != nil:
			t.Error("unexpected valid reference to a middleware")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("log with environment defined message when retrieving endpoint path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogEndpointErrorMessage
		RestEnvelopeMwLogEndpointErrorMessage = "test"
		defer func() { RestEnvelopeMwLogEndpointErrorMessage = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", "string")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "test", gomock.Any()).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)

		mw, e := generator(endpoint)
		switch {
		case mw != nil:
			t.Error("unexpected valid reference to a middleware")
		case e == nil:
			t.Error("didn't returned the expected error")
		case !errors.Is(e, slate.ErrConversion):
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
		}
	})

	t.Run("valid middleware creation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 123)
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)

		mw, e := generator(endpoint)
		switch {
		case mw == nil:
			t.Error("didn't returned a valid reference")
		case e != nil:
			t.Errorf("unexpected (%v) error", e)
		}
	})

	t.Run("calling the generated handler calls the given original handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 123)
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		calls := 0
		handler := mw(func(*gin.Context) {
			calls++
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)

		handler(ctx)

		if calls != 1 {
			t.Errorf("didn't called the original underlying handler")
		}
	})

	t.Run("parse data envelope stored in the response field of context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 123)
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		handler := mw(func(ctx *gin.Context) {
			RestSetResponse(ctx, NewEnvelope(200, []string{"data1", "data2"}, nil))
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":true,"error":[]},"data":["data1","data2"]}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("parse error stored in the response field of context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		handler := mw(func(ctx *gin.Context) {
			RestSetResponse(ctx, fmt.Errorf("error message"))
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"error message"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("parse invalid stored in the response field of context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		handler := mw(func(ctx *gin.Context) {
			RestSetResponse(ctx, "string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("parse panic error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		handler := mw(func(ctx *gin.Context) {
			panic(fmt.Errorf("error message"))
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"error message"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("panic non-error value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered observer update the service id value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.service.id", 321)
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		handler := mw(func(ctx *gin.Context) {
			RestSetResponse(ctx, fmt.Errorf("error message"))
		})

		_ = config.AddSupplier("id2", 1, newSource)

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:321.e:456.c:0","message":"error message"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered service id observer log on invalid new service id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "Invalid service id", slate.LogContext{"value": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered service id observer log on invalid new service id with environment defined channel", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogChannel
		RestEnvelopeMwLogChannel = "test"
		defer func() { RestEnvelopeMwLogChannel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("test", slate.ERROR, "Invalid service id", slate.LogContext{"value": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered service id observer log on invalid new service id with environment defined level", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogLevel
		RestEnvelopeMwLogLevel = "debug"
		defer func() { RestEnvelopeMwLogLevel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.DEBUG, "Invalid service id", slate.LogContext{"value": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered service id observer log on invalid new service id with environment defined message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogServiceErrorMessage
		RestEnvelopeMwLogServiceErrorMessage = "test"
		defer func() { RestEnvelopeMwLogServiceErrorMessage = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.service.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "test", slate.LogContext{"value": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered observer update the accept formats value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"text/xml"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered accept format observer log on invalid new format list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.accept", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "Invalid accept list", slate.LogContext{"list": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered accept format observer log on invalid new format list with environment defined channel", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogChannel
		RestEnvelopeMwLogChannel = "test"
		defer func() { RestEnvelopeMwLogChannel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.accept", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("test", slate.ERROR, "Invalid accept list", slate.LogContext{"list": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered accept format observer log on invalid new format list with environment defined level", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogLevel
		RestEnvelopeMwLogLevel = "debug"
		defer func() { RestEnvelopeMwLogLevel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.accept", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.DEBUG, "Invalid accept list", slate.LogContext{"list": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered accept format observer log on invalid new format list with environment defined message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogAcceptListErrorMessage
		RestEnvelopeMwLogAcceptListErrorMessage = "test"
		defer func() { RestEnvelopeMwLogAcceptListErrorMessage = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.accept", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "test", slate.LogContext{"list": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered accept format observer log on invalid new format list entry", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.accept", []interface{}{"application/json", 123})
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "Invalid accept list", slate.LogContext{"value": 123}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered accept format observer log on invalid new format list entry with environment defined channel", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogChannel
		RestEnvelopeMwLogChannel = "test"
		defer func() { RestEnvelopeMwLogChannel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.accept", []interface{}{"application/json", 123})
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("test", slate.ERROR, "Invalid accept list", slate.LogContext{"value": 123}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered accept format observer log on invalid new format list entry with environment defined level", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogLevel
		RestEnvelopeMwLogLevel = "debug"
		defer func() { RestEnvelopeMwLogLevel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.accept", []interface{}{"application/json", 123})
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.DEBUG, "Invalid accept list", slate.LogContext{"value": 123}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered accept format observer log on invalid new format list entry with environment defined message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogAcceptListErrorMessage
		RestEnvelopeMwLogAcceptListErrorMessage = "test"
		defer func() { RestEnvelopeMwLogAcceptListErrorMessage = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.accept", []interface{}{"application/json", 123})
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "test", slate.LogContext{"value": 123}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered endpoint id observer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.endpoints.index.id", 654)
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:654.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered endpoint id observer log on invalid new id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.endpoints.index.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "Invalid endpoint id", slate.LogContext{"value": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered endpoint id observer log on invalid new id with environment defined channel", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogChannel
		RestEnvelopeMwLogChannel = "test"
		defer func() { RestEnvelopeMwLogChannel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.endpoints.index.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("test", slate.ERROR, "Invalid endpoint id", slate.LogContext{"value": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered endpoint id observer log on invalid new id with environment defined level", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogLevel
		RestEnvelopeMwLogLevel = "debug"
		defer func() { RestEnvelopeMwLogLevel = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.endpoints.index.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.DEBUG, "Invalid endpoint id", slate.LogContext{"value": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})

	t.Run("registered endpoint id observer log on invalid new id with environment defined message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		prev := RestEnvelopeMwLogEndpointErrorMessage
		RestEnvelopeMwLogEndpointErrorMessage = "test"
		defer func() { RestEnvelopeMwLogEndpointErrorMessage = prev }()

		endpoint := "index"
		partial := slate.ConfigPartial{}
		_, _ = partial.Set("slate.api.rest.service.id", 123)
		_, _ = partial.Set("slate.api.rest.accept", []interface{}{"application/json"})
		_, _ = partial.Set("slate.api.rest.endpoints.index.id", 456)
		newPartial := slate.ConfigPartial{}
		_, _ = newPartial.Set("slate.api.rest.endpoints.index.id", "invalid")
		supplier := NewMockConfigSupplier(ctrl)
		supplier.EXPECT().Get("").Return(partial, nil).AnyTimes()
		newSource := NewMockConfigSupplier(ctrl)
		newSource.EXPECT().Get("").Return(newPartial, nil).Times(1)
		config := slate.NewConfig()
		_ = config.AddSupplier("id1", 0, supplier)
		logWriter := NewMockLogWriter(ctrl)
		logWriter.
			EXPECT().
			Signal("rest", slate.ERROR, "test", slate.LogContext{"value": "invalid"}).
			Return(nil).
			Times(1)
		logger := slate.NewLog()
		_ = logger.AddWriter("id", logWriter)
		generator, _ := NewRestEnvelopeMwGenerator(config, logger)
		mw, _ := generator(endpoint)

		_ = config.AddSupplier("id2", 1, newSource)

		handler := mw(func(ctx *gin.Context) {
			panic("string message")
		})

		gin.SetMode(gin.ReleaseMode)
		writer := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = &http.Request{}
		handler(ctx)

		expected := `{"status":{"success":false,"error":[{"code":"s:123.e:456.c:0","message":"internal server error"}]}}`

		if check := writer.Body.String(); check != expected {
			t.Errorf("(%v) when expecting (%v)", check, expected)
		}
	})
}

func Test_RestEnvelopeMwServiceRegister(t *testing.T) {
	t.Run("NewRestEnvelopeMwServiceRegister", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			if NewRestEnvelopeMwServiceRegister() == nil {
				t.Error("didn't returned a valid reference")
			}
		})

		t.Run("create with app reference", func(t *testing.T) {
			app := slate.NewApp()
			if sut := NewRestEnvelopeMwServiceRegister(app); sut == nil {
				t.Error("didn't returned a valid reference")
			} else if sut.App != app {
				t.Error("didn't stored the app reference")
			}
		})
	})

	t.Run("Provide", func(t *testing.T) {
		t.Run("nil container", func(t *testing.T) {
			if e := NewRestEnvelopeMwServiceRegister().Provide(nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expected (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("register components", func(t *testing.T) {
			container := slate.NewServiceContainer()
			sut := NewRestEnvelopeMwServiceRegister()

			e := sut.Provide(container)
			switch {
			case e != nil:
				t.Errorf("unexpected (%v) error", e)
			case !container.Has(RestEnvelopeMwContainerID):
				t.Errorf("no envelope middleware generator : %v", sut)
			}
		})

		t.Run("retrieving the envelope middleware generator", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := slate.NewServiceContainer()
			_ = slate.NewFileSystemServiceRegister().Provide(container)
			_ = slate.NewConfigServiceRegister().Provide(container)
			_ = slate.NewLogServiceRegister().Provide(container)
			_ = NewRestEnvelopeMwServiceRegister().Provide(container)

			cfg := slate.ConfigPartial{}
			_, _ = cfg.Set("slate.api.rest.accept", []interface{}{"application/json"})
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(cfg, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			_ = container.Add(slate.ConfigContainerID, func() *slate.Config {
				return config
			})

			sut, e := container.Get(RestEnvelopeMwContainerID)
			switch {
			case e != nil:
				t.Errorf("unexpected error (%v)", e)
			case sut == nil:
				t.Error("didn't returned a reference to service")
			default:
				switch sut.(type) {
				case RestEnvelopeMwGenerator:
				default:
					t.Error("didn't returned the envelope middleware generator")
				}
			}
		})
	})
}
