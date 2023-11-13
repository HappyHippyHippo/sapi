package sapi

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/slate"
)

func Test_RestProcess(t *testing.T) {
	t.Run("NewRestProcess", func(t *testing.T) {
		t.Run("nil config", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := slate.NewLog()
			engine := NewMockRestEngine(ctrl)

			sut, e := NewRestProcess(nil, logger, engine)
			switch {
			case sut != nil:
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

			config := slate.NewConfig()
			engine := NewMockRestEngine(ctrl)

			sut, e := NewRestProcess(config, nil, engine)
			switch {
			case sut != nil:
				t.Error("returned a valid reference")
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("nil engine", func(t *testing.T) {
			config := slate.NewConfig()
			logger := slate.NewLog()

			sut, e := NewRestProcess(config, logger, nil)
			switch {
			case sut != nil:
				t.Error("returned a valid reference")
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("error while retrieving configuration", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			partial := slate.ConfigPartial{}
			_, _ = partial.Set("slate.api.rest", "string")
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logger := slate.NewLog()
			engine := NewMockRestEngine(ctrl)

			sut, e := NewRestProcess(config, logger, engine)
			switch {
			case sut != nil:
				t.Error("returned a valid reference")
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrConversion):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
			}
		})

		t.Run("error while retrieving configuration from env path", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			prev := RestConfigPath
			RestConfigPath = "test"
			defer func() { RestConfigPath = prev }()

			partial := slate.ConfigPartial{"test": "string"}
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logger := slate.NewLog()
			engine := NewMockRestEngine(ctrl)

			sut, e := NewRestProcess(config, logger, engine)
			switch {
			case sut != nil:
				t.Error("returned a valid reference")
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrConversion):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
			}
		})

		t.Run("error while populating configuration", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			partial := slate.ConfigPartial{}
			_, _ = partial.Set("slate.api.rest.watchdog", 123)
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logger := slate.NewLog()
			engine := NewMockRestEngine(ctrl)

			sut, e := NewRestProcess(config, logger, engine)
			switch {
			case sut != nil:
				t.Error("returned a valid reference")
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrConversion):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
			}
		})

		t.Run("invalid log level", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			partial := slate.ConfigPartial{}
			_, _ = partial.Set("slate.api.rest.log.level", "invalid")
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logger := slate.NewLog()
			engine := NewMockRestEngine(ctrl)

			sut, e := NewRestProcess(config, logger, engine)
			switch {
			case sut != nil:
				t.Error("returned a valid reference")
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrConversion):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
			}
		})

		t.Run("successful process creation", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			partial := slate.ConfigPartial{}
			_, _ = partial.Set("slate.api.rest", slate.ConfigPartial{})
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logger := slate.NewLog()
			engine := NewMockRestEngine(ctrl)

			sut, e := NewRestProcess(config, logger, engine)
			switch {
			case sut == nil:
				t.Error("didn't returned the expected valid reference")
			case sut.Service() != RestWatchdogName:
				t.Error("didn't returned the expected valid reference")
			case e != nil:
				t.Errorf("unexpected (%v) error", e)
			}
		})

		t.Run("successful process creation with name from env", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			prev := RestWatchdogName
			RestWatchdogName = "test"
			defer func() { RestWatchdogName = prev }()

			partial := slate.ConfigPartial{}
			_, _ = partial.Set("slate.api.rest", slate.ConfigPartial{})
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logger := slate.NewLog()
			engine := NewMockRestEngine(ctrl)

			sut, e := NewRestProcess(config, logger, engine)
			switch {
			case sut == nil:
				t.Error("didn't returned the expected valid reference")
			case sut.Service() != RestWatchdogName:
				t.Error("didn't returned the expected valid reference")
			case e != nil:
				t.Errorf("unexpected (%v) error", e)
			}
		})

		t.Run("successful process run", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			partial := slate.ConfigPartial{}
			_, _ = partial.Set("slate.api.rest", slate.ConfigPartial{})
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logWriter := NewMockLogWriter(ctrl)
			gomock.InOrder(
				logWriter.
					EXPECT().
					Signal("rest", slate.INFO, "[service:rest] service starting ...", slate.LogContext{"port": 80}).
					Return(nil),
				logWriter.
					EXPECT().
					Signal("rest", slate.INFO, "[service:rest] service terminated").
					Return(nil),
			)
			logger := slate.NewLog()
			_ = logger.AddWriter("id", logWriter)
			engine := NewMockRestEngine(ctrl)
			engine.EXPECT().Run(":80").Return(nil).Times(1)

			sut, _ := NewRestProcess(config, logger, engine)
			if e := sut.Runner()(); e != nil {
				t.Errorf("unexpected (%v) error", e)
			}
		})

		t.Run("successful process run with config values", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			name := "watchdog name"
			port := 1234
			logLevel := slate.FATAL
			logChannel := "test channel"
			logStartMessage := "start message"
			logEndMessage := "end message"

			partial := slate.ConfigPartial{}
			_, _ = partial.Set("slate.api.rest.watchdog", name)
			_, _ = partial.Set("slate.api.rest.port", port)
			_, _ = partial.Set("slate.api.rest.log.level", slate.LogLevelMapName[logLevel])
			_, _ = partial.Set("slate.api.rest.log.channel", logChannel)
			_, _ = partial.Set("slate.api.rest.log.message.start", logStartMessage)
			_, _ = partial.Set("slate.api.rest.log.message.end", logEndMessage)
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logWriter := NewMockLogWriter(ctrl)
			gomock.InOrder(
				logWriter.
					EXPECT().
					Signal(logChannel, logLevel, logStartMessage, slate.LogContext{"port": port}).
					Return(nil),
				logWriter.
					EXPECT().
					Signal(logChannel, logLevel, logEndMessage).
					Return(nil),
			)
			logger := slate.NewLog()
			_ = logger.AddWriter("id", logWriter)
			engine := NewMockRestEngine(ctrl)
			engine.EXPECT().Run(":1234").Return(nil).Times(1)

			sut, _ := NewRestProcess(config, logger, engine)
			if chk := sut.Service(); chk != name {
				t.Errorf("(%v) when expected (%v)", chk, name)
			} else if e := sut.Runner()(); e != nil {
				t.Errorf("unexpected (%v) error", e)
			}
		})

		t.Run("failure when running process", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			errorMessage := "error message"
			expected := fmt.Errorf("%s", errorMessage)

			partial := slate.ConfigPartial{}
			_, _ = partial.Set("slate.api.rest", slate.ConfigPartial{})
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logWriter := NewMockLogWriter(ctrl)
			gomock.InOrder(
				logWriter.
					EXPECT().
					Signal("rest", slate.INFO, "[service:rest] service starting ...", slate.LogContext{"port": 80}).
					Return(nil),
				logWriter.
					EXPECT().
					Signal("rest", slate.FATAL, "[service:rest] service error", slate.LogContext{"error": errorMessage}).
					Return(nil),
			)
			logger := slate.NewLog()
			_ = logger.AddWriter("id", logWriter)
			engine := NewMockRestEngine(ctrl)
			engine.EXPECT().Run(":80").Return(expected).Times(1)

			sut, _ := NewRestProcess(config, logger, engine)
			e := sut.Runner()()
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, expected):
				t.Errorf("(%v) when expecting (%v)", e, expected)
			}
		})

		t.Run("failure when running process with values from config", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			errorMessage := "error message"
			expected := fmt.Errorf("%s", errorMessage)

			name := "watchdog name"
			port := 1234
			logLevel := slate.FATAL
			logChannel := "test channel"
			logStartMessage := "start message"
			logErrorMessage := "error message"

			partial := slate.ConfigPartial{}
			_, _ = partial.Set("slate.api.rest.watchdog", name)
			_, _ = partial.Set("slate.api.rest.port", port)
			_, _ = partial.Set("slate.api.rest.log.level", slate.LogLevelMapName[logLevel])
			_, _ = partial.Set("slate.api.rest.log.channel", logChannel)
			_, _ = partial.Set("slate.api.rest.log.message.start", logStartMessage)
			_, _ = partial.Set("slate.api.rest.log.message.error", logErrorMessage)
			supplier := NewMockConfigSupplier(ctrl)
			supplier.EXPECT().Get("").Return(partial, nil).Times(1)
			config := slate.NewConfig()
			_ = config.AddSupplier("id", 0, supplier)
			logWriter := NewMockLogWriter(ctrl)
			gomock.InOrder(
				logWriter.
					EXPECT().
					Signal(logChannel, logLevel, logStartMessage, slate.LogContext{"port": port}).
					Return(nil),
				logWriter.
					EXPECT().
					Signal(logChannel, logLevel, logErrorMessage, slate.LogContext{"error": errorMessage}).
					Return(nil),
			)
			logger := slate.NewLog()
			_ = logger.AddWriter("id", logWriter)
			engine := NewMockRestEngine(ctrl)
			engine.EXPECT().Run(":1234").Return(expected).Times(1)

			sut, _ := NewRestProcess(config, logger, engine)
			e := sut.Runner()()
			switch {
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, expected):
				t.Errorf("(%v) when expecting (%v)", e, expected)
			}
		})
	})
}

func Test_RestLoader(t *testing.T) {
	t.Run("NewRestLoader", func(t *testing.T) {
		t.Run("nil loader", func(t *testing.T) {
			sut, e := NewRestLoader(nil, nil)
			switch {
			case sut != nil:
				t.Error("returned a valid reference")
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("valid initialization", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			engine := NewMockRestEngine(ctrl)
			sut, e := NewRestLoader(engine, nil)

			switch {
			case e != nil:
				t.Errorf("return the unexpected error (%v)", e)
			case sut == nil:
				t.Error("didn't return the expected validation instance")
			}
		})

		t.Run("storing of registers", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			engine := NewMockRestEngine(ctrl)
			register := NewMockRestEndpointRegister(ctrl)
			sut, e := NewRestLoader(engine, []RestEndpointRegister{register})

			switch {
			case e != nil:
				t.Errorf("return the unexpected error (%v)", e)
			case sut == nil:
				t.Error("didn't return the expected validation instance")
			case len(sut.registers) != 1:
				t.Errorf("didn't stored any endpoint register")
			case sut.registers[0] != register:
				t.Errorf("didn't stored the expected endpoint register")
			}
		})
	})

	t.Run("Load", func(t *testing.T) {
		t.Run("error on registration", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := fmt.Errorf("error message")
			engine := NewMockRestEngine(ctrl)
			register1 := NewMockRestEndpointRegister(ctrl)
			register1.EXPECT().Reg(engine).Return(nil).Times(1)
			register2 := NewMockRestEndpointRegister(ctrl)
			register2.EXPECT().Reg(engine).Return(expected).Times(1)
			sut, _ := NewRestLoader(engine, []RestEndpointRegister{register1, register2})

			if e := sut.Load(); e == nil {
				t.Errorf("didn't return the expected error")
			} else if !errors.Is(e, expected) {
				t.Errorf("(%v) when expecting (%v)", e, expected)
			}
		})

		t.Run("valid run", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			engine := NewMockRestEngine(ctrl)
			register1 := NewMockRestEndpointRegister(ctrl)
			register1.EXPECT().Reg(engine).Return(nil).Times(1)
			register2 := NewMockRestEndpointRegister(ctrl)
			register2.EXPECT().Reg(engine).Return(nil).Times(1)
			sut, _ := NewRestLoader(engine, []RestEndpointRegister{register1, register2})

			if e := sut.Load(); e != nil {
				t.Errorf("return the unexpected error (%v)", e)
			}
		})
	})
}

func Test_RestServiceRegister(t *testing.T) {
	t.Run("NewRestServiceRegister", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			if NewRestServiceRegister() == nil {
				t.Error("didn't returned a valid reference")
			}
		})

		t.Run("create with app reference", func(t *testing.T) {
			app := slate.NewApp()
			if sut := NewRestServiceRegister(app); sut == nil {
				t.Error("didn't returned a valid reference")
			} else if sut.App != app {
				t.Error("didn't stored the app reference")
			}
		})
	})

	t.Run("Provide", func(t *testing.T) {
		t.Run("nil container", func(t *testing.T) {
			if e := NewRestServiceRegister().Provide(nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expected (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("register components", func(t *testing.T) {
			container := slate.NewServiceContainer()
			sut := NewRestServiceRegister()

			e := sut.Provide(container)
			switch {
			case e != nil:
				t.Errorf("unexpected (%v) error", e)
			case !container.Has(RestAllEndpointRegistersContainerID):
				t.Errorf("no list of endpoint registers instance : %v", sut)
			case !container.Has(RestContainerID):
				t.Errorf("no rest engine service : %v", sut)
			case !container.Has(RestProcessContainerID):
				t.Errorf("no rest watchdog process : %v", sut)
			case !container.Has(RestLoaderContainerID):
				t.Errorf("no rest service loader : %v", sut)
			}
		})
	})

	t.Run("retrieving endpoint registers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := slate.NewServiceContainer()
		_ = NewRestServiceRegister().Provide(container)

		register := NewMockRestEndpointRegister(ctrl)
		_ = container.Add("dialect.id", func() RestEndpointRegister {
			return register
		}, RestEndpointRegisterTag)

		creators, e := container.Get(RestAllEndpointRegistersContainerID)
		switch {
		case e != nil:
			t.Errorf("unexpected error (%v)", e)
		case creators == nil:
			t.Error("didn't returned a valid reference")
		default:
			switch s := creators.(type) {
			case []RestEndpointRegister:
				if s[0] != register {
					t.Error("didn't return a endpoint register slice populated with the expected register instance")
				}
			default:
				t.Error("didn't return a endpoint register slice")
			}
		}
	})

	t.Run("retrieving rest engine", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := slate.NewServiceContainer()
		_ = NewRestServiceRegister().Provide(container)

		sut, e := container.Get(RestContainerID)
		switch {
		case e != nil:
			t.Errorf("unexpected error (%v)", e)
		case sut == nil:
			t.Error("didn't returned a reference to service")
		default:
			switch sut.(type) {
			case RestEngine:
			default:
				t.Error("didn't returned the rest engine instance")
			}
		}
	})

	t.Run("retrieving rest process", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := slate.NewServiceContainer()
		_ = slate.NewConfigServiceRegister().Provide(container)
		_ = slate.NewLogServiceRegister().Provide(container)
		_ = NewRestServiceRegister().Provide(container)

		sut, e := container.Get(RestProcessContainerID)
		switch {
		case e != nil:
			t.Errorf("unexpected error (%v)", e)
		case sut == nil:
			t.Error("didn't returned a reference to service")
		default:
			switch sut.(type) {
			case *RestProcess:
			default:
				t.Error("didn't returned the rest process instance")
			}
		}
	})

	t.Run("retrieving rest loader", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := slate.NewServiceContainer()
		_ = slate.NewConfigServiceRegister().Provide(container)
		_ = slate.NewLogServiceRegister().Provide(container)
		_ = NewRestServiceRegister().Provide(container)

		sut, e := container.Get(RestLoaderContainerID)
		switch {
		case e != nil:
			t.Errorf("unexpected error (%v)", e)
		case sut == nil:
			t.Error("didn't returned a reference to service")
		default:
			switch sut.(type) {
			case *RestLoader:
			default:
				t.Error("didn't returned the rest loader instance")
			}
		}
	})

	t.Run("Boot", func(t *testing.T) {
		t.Run("nil container", func(t *testing.T) {
			if e := NewRestServiceRegister(nil).Boot(nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expected (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("error retrieving loader", func(t *testing.T) {
			container := slate.NewServiceContainer()
			_ = slate.NewFileSystemServiceRegister(nil).Provide(container)
			_ = slate.NewConfigServiceRegister(nil).Provide(container)
			sut := NewRestServiceRegister(nil)
			_ = sut.Provide(container)
			_ = container.Add(RestLoaderContainerID, func() (*RestLoader, error) {
				return nil, fmt.Errorf("error message")
			})

			if e := sut.Boot(container); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrServiceContainer) {
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrServiceContainer)
			}
		})

		t.Run("invalid loader", func(t *testing.T) {
			container := slate.NewServiceContainer()
			_ = slate.NewFileSystemServiceRegister(nil).Provide(container)
			_ = slate.NewConfigServiceRegister(nil).Provide(container)
			sut := NewRestServiceRegister(nil)
			_ = sut.Provide(container)
			_ = container.Add(RestLoaderContainerID, func() string {
				return "message"
			})

			if e := sut.Boot(container); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrConversion) {
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrConversion)
			}
		})

		t.Run("valid simple boot with no registers", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := slate.NewServiceContainer()
			_ = slate.NewFileSystemServiceRegister(nil).Provide(container)
			_ = slate.NewConfigServiceRegister(nil).Provide(container)
			sut := NewRestServiceRegister(nil)
			_ = sut.Provide(container)

			if e := sut.Boot(container); e != nil {
				t.Errorf("unexpected error (%v)", e)
			}
		})

		t.Run("valid simple boot with registers", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			container := slate.NewServiceContainer()
			_ = slate.NewFileSystemServiceRegister(nil).Provide(container)
			_ = slate.NewConfigServiceRegister(nil).Provide(container)
			sut := NewRestServiceRegister(nil)
			_ = sut.Provide(container)
			engine := NewMockRestEngine(ctrl)
			reg := NewMockRestEndpointRegister(ctrl)
			reg.EXPECT().Reg(engine).Return(nil).Times(1)
			_ = container.Add(RestContainerID, func() RestEngine {
				return engine
			})
			_ = container.Add("reg.1", func() RestEndpointRegister {
				return reg
			}, RestEndpointRegisterTag)

			if e := sut.Boot(container); e != nil {
				t.Errorf("unexpected error (%v)", e)
			}
		})
	})
}
