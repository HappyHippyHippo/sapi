package sapi

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/slate"
)

func Test_validation_err(t *testing.T) {
	t.Run("errValidationTranslatorNotFound", func(t *testing.T) {
		arg := "dummy argument"
		context := map[string]interface{}{"field": "value"}
		message := "dummy argument : validation translator not found"

		t.Run("creation without context", func(t *testing.T) {
			if e := errValidationTranslatorNotFound(arg); !errors.Is(e, ErrValidationTranslatorNotFound) {
				t.Errorf("error not a instance of ErrValidationTranslatorNotFound")
			} else if e.Error() != message {
				t.Errorf("error message (%v) not same as expected (%v)", e, message)
			} else {
				var te *slate.Error
				if !errors.As(e, &te) {
					t.Errorf("didn't returned a slate error instance")
				}
			}
		})

		t.Run("creation with context", func(t *testing.T) {
			if e := errValidationTranslatorNotFound(arg, context); !errors.Is(e, ErrValidationTranslatorNotFound) {
				t.Errorf("error not a instance of ErrValidationTranslatorNotFound")
			} else if e.Error() != message {
				t.Errorf("error message (%v) not same as expected (%v)", e, message)
			} else {
				var te *slate.Error
				if !errors.As(e, &te) {
					t.Errorf("didn't returned a slate error instance")
				}
			}
		})
	})
}

func Test_ValidationUniversalTranslator(t *testing.T) {
	t.Run("NewValidationUniversalTranslator", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			if sut := NewValidationUniversalTranslator(); sut == nil {
				t.Errorf("didn't create the desired universal translator")
			}
		})
	})
}

func Test_ValidationTranslator(t *testing.T) {
	t.Run("NewValidationTranslator", func(t *testing.T) {
		t.Run("locale not found", func(t *testing.T) {
			prev := ValidationLocale
			ValidationLocale = "invalid"
			defer func() { ValidationLocale = prev }()

			translator := NewValidationUniversalTranslator()
			sut, e := NewValidationTranslator(translator)
			switch {
			case sut != nil:
				t.Errorf("returned an unexpected valid reference")
			case e == nil:
				t.Errorf("didn't returned the expected error")
			}
		})

		t.Run("valid creation", func(t *testing.T) {
			translator := NewValidationUniversalTranslator()
			sut, e := NewValidationTranslator(translator)
			switch {
			case sut == nil:
				t.Errorf("didn't returned the expected reference to the translator")
			case e != nil:
				t.Errorf("returned the unexpected error : %v", e)
			}
		})
	})
}

func Test_ValidationParser(t *testing.T) {
	t.Run("NewValidationParser", func(t *testing.T) {
		t.Run("nil translator", func(t *testing.T) {
			parser, e := NewValidationParser(nil)
			switch {
			case parser != nil:
				t.Error("returned a valid reference")
			case e == nil:
				t.Error("didn't returned the expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("new parser", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			translator := NewMockTranslator(ctrl)

			p, e := NewValidationParser(translator)
			switch {
			case p == nil:
				t.Error("didn't returned a valid reference")
			case e != nil:
				t.Errorf("return the (%v) error", e)
			case p.translator != translator:
				t.Error("didn't stored the translator reference")
			}
		})
	})

	t.Run("Parse", func(t *testing.T) {
		t.Run("nil value", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			translator := NewMockTranslator(ctrl)
			sut, _ := NewValidationParser(translator)

			resp, e := sut.Parse(nil, []validator.FieldError{})
			switch {
			case resp != nil:
				t.Error("unexpectedly valid instance of a response")
			case e == nil:
				t.Error("didn't returned an expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Error("not nil pointer error")
			}
		})

		t.Run("no-op on nil error list", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			value := struct{ Message string }{Message: "message"}
			translator := NewMockTranslator(ctrl)
			sut, _ := NewValidationParser(translator)

			resp, e := sut.Parse(value, nil)
			switch {
			case e != nil:
				t.Errorf("unexpected ('%v') error", e)
			case resp != nil:
				t.Error("unexpectedly valid instance of a response")
			}
		})

		t.Run("no-op if error list is empty", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			value := struct{ Message string }{Message: "message"}
			translator := NewMockTranslator(ctrl)
			sut, _ := NewValidationParser(translator)

			resp, e := sut.Parse(value, []validator.FieldError{})
			switch {
			case e != nil:
				t.Errorf("unexpected ('%v') error", e)
			case resp != nil:
				t.Error("unexpectedly valid instance of a response")
			}
		})

		t.Run("invalid nil field error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := struct {
				Field string `vparam:"string"`
			}{Field: "message"}
			translator := NewMockTranslator(ctrl)

			sut, _ := NewValidationParser(translator)

			resp, e := sut.Parse(data, []validator.FieldError{nil})
			switch {
			case resp != nil:
				t.Error("unexpectedly valid instance of a response")
			case e == nil:
				t.Error("didn't returned an expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Error("not nil pointer error")
			}
		})

		t.Run("error retrieving field/param value", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := struct {
				Field string `vparam:"string"`
			}{Field: "message"}
			expected := `strconv.Atoi: parsing "string": invalid syntax`
			translator := NewMockTranslator(ctrl)
			fieldError := NewMockFieldError(ctrl)
			fieldError.EXPECT().StructField().Return("Field").Times(1)

			sut, _ := NewValidationParser(translator)

			resp, e := sut.Parse(data, []validator.FieldError{fieldError})
			switch {
			case resp != nil:
				t.Error("unexpectedly valid instance of a response")
			case e == nil:
				t.Error("didn't returned an expected error")
			case e.Error() != expected:
				t.Error("not nil pointer error")
			}
		})

		t.Run("generating error for a non-tagged field", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := struct {
				Field int `validate:"gt=0"`
			}{Field: 0}
			errMsg := "error message"
			expected := "c:89"
			translator := NewMockTranslator(ctrl)
			fieldError := NewMockFieldError(ctrl)
			fieldError.EXPECT().StructField().Return("Field").Times(1)
			fieldError.EXPECT().Translate(translator).Return(errMsg).Times(1)
			fieldError.EXPECT().Tag().Return("gt").Times(1)

			sut, _ := NewValidationParser(translator)

			resp, e := sut.Parse(data, []validator.FieldError{fieldError})
			switch {
			case e != nil:
				t.Errorf("return the unexpected error (%v)", e)
			case len(resp.Status.Errors) == 0:
				t.Errorf("didn't stored the expected error")
			case resp.Status.Errors[0].GetCode() != expected:
				t.Errorf("(%v) when expecting (%s)", resp.Status.Errors[0].GetCode(), expected)
			case resp.Status.Errors[0].GetMessage() != errMsg:
				t.Errorf("(%v) when expecting (%v)", resp.Status.Errors[0].GetMessage(), errMsg)
			}
		})

		t.Run("generating error for an unrecognized error tag", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := struct {
				Field int `validate:"gt=0"`
			}{Field: 0}
			errMsg := "error message"
			expected := "c:0"
			translator := NewMockTranslator(ctrl)
			fieldError := NewMockFieldError(ctrl)
			fieldError.EXPECT().StructField().Return("Field").Times(1)
			fieldError.EXPECT().Translate(translator).Return(errMsg).Times(1)
			fieldError.EXPECT().Tag().Return("unrecognized").Times(1)

			sut, _ := NewValidationParser(translator)

			resp, e := sut.Parse(data, []validator.FieldError{fieldError})
			switch {
			case e != nil:
				t.Errorf("return the unexpected error (%v)", e)
			case resp.Status.Errors[0].GetCode() != expected:
				t.Errorf("(%v) when expecting (%s)", resp.Status.Errors[0].GetCode(), expected)
			case resp.Status.Errors[0].GetMessage() != errMsg:
				t.Errorf("(%v) when expecting (%v)", resp.Status.Errors[0].GetMessage(), errMsg)
			}
		})

		t.Run("generating error with all information", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := struct {
				Field int `validate:"gt=0" vparam:"10"`
			}{Field: 0}
			expected := "p:10.c:89"
			errMsg := "error message"
			translator := NewMockTranslator(ctrl)
			fieldError := NewMockFieldError(ctrl)
			fieldError.EXPECT().StructField().Return("Field").Times(1)
			fieldError.EXPECT().Translate(translator).Return(errMsg).Times(1)
			fieldError.EXPECT().Tag().Return("gt").Times(1)

			sut, _ := NewValidationParser(translator)

			resp, e := sut.Parse(data, []validator.FieldError{fieldError})
			switch {
			case e != nil:
				t.Errorf("return the unexpected error (%v)", e)
			case resp.Status.Errors[0].GetCode() != expected:
				t.Errorf("(%v) when expecting (%s)", resp.Status.Errors[0].GetCode(), expected)
			case resp.Status.Errors[0].GetMessage() != errMsg:
				t.Errorf("(%v) when expecting (%v)", resp.Status.Errors[0].GetMessage(), errMsg)
			}
		})
	})

	t.Run("AddError", func(t *testing.T) {
		t.Run("adding a new error mapping value", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mappedErrorName := "unrecognized"
			mappedErrorCode := 10000
			expected := fmt.Sprintf("c:%d", mappedErrorCode)
			data := struct {
				Field int `validate:"gt=0"`
			}{Field: 0}
			errMsg := "error message"
			translator := NewMockTranslator(ctrl)
			fieldError := NewMockFieldError(ctrl)
			fieldError.EXPECT().StructField().Return("Field").Times(1)
			fieldError.EXPECT().Translate(translator).Return(errMsg).Times(1)
			fieldError.EXPECT().Tag().Return(mappedErrorName).Times(1)

			sut, _ := NewValidationParser(translator)
			sut.AddError(mappedErrorName, mappedErrorCode)

			resp, e := sut.Parse(data, []validator.FieldError{fieldError})
			switch {
			case e != nil:
				t.Errorf("return the unexpected error (%v)", e)
			case resp.Status.Errors[0].GetCode() != expected:
				t.Errorf("(%v) when expecting (%s)", resp.Status.Errors[0].GetCode(), expected)
			case resp.Status.Errors[0].GetMessage() != errMsg:
				t.Errorf("(%v) when expecting (%v)", resp.Status.Errors[0].GetMessage(), errMsg)
			}
		})
	})
}

func Test_Validator(t *testing.T) {
	t.Run("NewValidator", func(t *testing.T) {
		t.Run("nil validate", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			translator := NewMockTranslator(ctrl)
			parser, _ := NewValidationParser(translator)

			check, e := NewValidator(nil, parser)
			switch {
			case check != nil:
				t.Error("return an unexpected valid validator instance")
			case e == nil:
				t.Error("didn't return an expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("nil parser", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			translator := NewMockTranslator(ctrl)

			check, e := NewValidator(translator, nil)
			switch {
			case check != nil:
				t.Error("return an unexpected valid validator instance")
			case e == nil:
				t.Error("didn't return an expected error")
			case !errors.Is(e, slate.ErrNilPointer):
				t.Errorf("(%v) when expecting (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("construct", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			translator := NewMockTranslator(ctrl)
			translator.
				EXPECT().
				Add(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			translator.
				EXPECT().
				AddCardinal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			parser, _ := NewValidationParser(translator)

			if check, e := NewValidator(translator, parser); e != nil {
				t.Errorf("return the unexpected error (%v)", e)
			} else if check == nil {
				t.Error("didn't return the expected validation instance")
			}
		})
	})

	t.Run("call", func(t *testing.T) {
		t.Run("nil data", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			expected := errNilPointer("value")
			translator := NewMockTranslator(ctrl)
			translator.
				EXPECT().
				Add(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			translator.
				EXPECT().
				AddCardinal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			parser, _ := NewValidationParser(translator)
			sut, _ := NewValidator(translator, parser)

			env, e := sut(nil)
			switch {
			case env != nil:
				t.Errorf("return the unexpected envelope (%v)", env)
			case e == nil:
				t.Error("didn't return an expected error")
			case e.Error() != expected.Error():
				t.Errorf("(%v) when expecting (%v)", e, expected)
			}
		})

		t.Run("no error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := struct {
				Field1 int `validate:"gt=0,lte=10" vparam:"1"`
				Field2 int `validate:"gt=10,lte=20" vparam:"2"`
			}{Field1: 1, Field2: 11}
			translator := NewMockTranslator(ctrl)
			translator.
				EXPECT().
				Add(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			translator.
				EXPECT().
				AddCardinal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			parser, _ := NewValidationParser(translator)
			sut, _ := NewValidator(translator, parser)

			if env, e := sut(data); e != nil {
				t.Errorf("unexpected (%v) error", e)
			} else if env != nil {
				t.Errorf("returned the unexpected envelope (%v)", env)
			}
		})

		t.Run("error parsing error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := struct {
				Field1 int `validate:"gt=0,lte=10" vparam:"string"`
				Field2 int `validate:"gt=10,lte=20" vparam:"2"`
			}{Field1: 11, Field2: 11}
			expected := "strconv.Atoi: parsing \"string\": invalid syntax"
			translator := NewMockTranslator(ctrl)
			translator.
				EXPECT().
				Add(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			translator.
				EXPECT().
				AddCardinal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			parser, _ := NewValidationParser(translator)
			sut, _ := NewValidator(translator, parser)

			resp, e := sut(data)
			switch {
			case resp != nil:
				t.Error("unexpected instance of the response envelope")
			case e == nil:
				t.Error("didn't returned the expected error")
			case e.Error() != expected:
				t.Errorf("(%v) when expecting (%v)", e, expected)
			}
		})

		t.Run("parse error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			data := struct {
				Field1 int `validate:"gt=0,lte=10" vparam:"1"`
				Field2 int `validate:"gt=10,lte=20" vparam:"2"`
			}{Field1: 11, Field2: 11}
			errMsg := "error message"
			expected := NewEnvelope(http.StatusBadRequest, nil, nil)
			expected.AddError(NewEnvelopeStatusError(92, errMsg).SetParam(1))
			translator := NewMockTranslator(ctrl)
			translator.
				EXPECT().
				Add(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			translator.
				EXPECT().
				AddCardinal(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()
			translator.
				EXPECT().
				FmtNumber(float64(10), uint64(0)).
				Return("10").
				Times(1)
			translator.
				EXPECT().
				T("lte-number", "Field1", gomock.Any()).
				Return(errMsg, nil).
				Times(1)
			parser, _ := NewValidationParser(translator)
			sut, _ := NewValidator(translator, parser)

			if resp, e := sut(data); e != nil {
				t.Errorf("unexpected (%v) error", e)
			} else if !reflect.DeepEqual(resp, expected) {
				t.Errorf("(%v) when expecting (%v)", resp, expected)
			}
		})
	})
}

func Test_ValidationServiceRegister(t *testing.T) {
	t.Run("NewWatchdogServiceRegister", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			if NewValidationServiceRegister() == nil {
				t.Error("didn't returned a valid reference")
			}
		})

		t.Run("create with app reference", func(t *testing.T) {
			app := slate.NewApp()
			if sut := NewValidationServiceRegister(app); sut == nil {
				t.Error("didn't returned a valid reference")
			} else if sut.App != app {
				t.Error("didn't stored the app reference")
			}
		})
	})

	t.Run("Provide", func(t *testing.T) {
		t.Run("nil container", func(t *testing.T) {
			if e := NewValidationServiceRegister().Provide(nil); e == nil {
				t.Error("didn't returned the expected error")
			} else if !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("(%v) when expected (%v)", e, slate.ErrNilPointer)
			}
		})

		t.Run("register components", func(t *testing.T) {
			container := slate.NewServiceContainer()
			sut := NewValidationServiceRegister()

			e := sut.Provide(container)
			switch {
			case e != nil:
				t.Errorf("unexpected (%v) error", e)
			case !container.Has(ValidationUniversalTranslatorContainerID):
				t.Errorf("no universal translator instance : %v", sut)
			case !container.Has(ValidationTranslatorContainerID):
				t.Errorf("no translator instance : %v", sut)
			case !container.Has(ValidationParserContainerID):
				t.Errorf("no parser instance : %v", sut)
			case !container.Has(ValidationContainerID):
				t.Errorf("no trnalsator creator : %v", sut)
			}
		})
	})

	t.Run("retrieving universal translator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := slate.NewServiceContainer()
		_ = NewValidationServiceRegister().Provide(container)

		sut, e := container.Get(ValidationUniversalTranslatorContainerID)
		switch {
		case e != nil:
			t.Errorf("unexpected error (%v)", e)
		case sut == nil:
			t.Error("didn't returned a reference to service")
		default:
			switch sut.(type) {
			case *ut.UniversalTranslator:
			default:
				t.Error("didn't returned the universal translator")
			}
		}
	})

	t.Run("error retrieving universal translator when retrieving translator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := fmt.Errorf("error message")
		container := slate.NewServiceContainer()
		_ = NewValidationServiceRegister().Provide(container)
		_ = container.Add(ValidationUniversalTranslatorContainerID, func() (*ut.UniversalTranslator, error) {
			return nil, expected
		})

		if _, e := container.Get(ValidationTranslatorContainerID); e == nil {
			t.Error("didn't returned the expected error")
		} else if !errors.Is(e, slate.ErrServiceContainer) {
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrServiceContainer)
		}
	})

	t.Run("retrieving translator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := slate.NewServiceContainer()
		_ = NewValidationServiceRegister().Provide(container)

		sut, e := container.Get(ValidationTranslatorContainerID)
		switch {
		case e != nil:
			t.Errorf("unexpected error (%v)", e)
		case sut == nil:
			t.Error("didn't returned a reference to service")
		default:
			switch sut.(type) {
			case ut.Translator:
			default:
				t.Error("didn't returned the translator")
			}
		}
	})

	t.Run("error retrieving translator when retrieving parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := fmt.Errorf("error message")
		container := slate.NewServiceContainer()
		_ = NewValidationServiceRegister().Provide(container)
		_ = container.Add(ValidationTranslatorContainerID, func() (ut.Translator, error) {
			return nil, expected
		})

		if _, e := container.Get(ValidationParserContainerID); e == nil {
			t.Error("didn't returned the expected error")
		} else if !errors.Is(e, slate.ErrServiceContainer) {
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrServiceContainer)
		}
	})

	t.Run("retrieving parser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := slate.NewServiceContainer()
		_ = NewValidationServiceRegister().Provide(container)

		sut, e := container.Get(ValidationParserContainerID)
		switch {
		case e != nil:
			t.Errorf("unexpected error (%v)", e)
		case sut == nil:
			t.Error("didn't returned a reference to service")
		default:
			switch sut.(type) {
			case *ValidationParser:
			default:
				t.Error("didn't returned the parser")
			}
		}
	})

	t.Run("error retrieving translator when retrieving validator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := fmt.Errorf("error message")
		container := slate.NewServiceContainer()
		_ = NewValidationServiceRegister().Provide(container)
		_ = container.Add(ValidationTranslatorContainerID, func() (ut.Translator, error) {
			return nil, expected
		})

		if _, e := container.Get(ValidationContainerID); e == nil {
			t.Error("didn't returned the expected error")
		} else if !errors.Is(e, slate.ErrServiceContainer) {
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrServiceContainer)
		}
	})

	t.Run("error retrieving parser when retrieving validator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := fmt.Errorf("error message")
		container := slate.NewServiceContainer()
		_ = NewValidationServiceRegister().Provide(container)
		_ = container.Add(ValidationParserContainerID, func() (*ValidationParser, error) {
			return nil, expected
		})

		if _, e := container.Get(ValidationContainerID); e == nil {
			t.Error("didn't returned the expected error")
		} else if !errors.Is(e, slate.ErrServiceContainer) {
			t.Errorf("(%v) when expecting (%v)", e, slate.ErrServiceContainer)
		}
	})

	t.Run("retrieving validator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := slate.NewServiceContainer()
		_ = NewValidationServiceRegister().Provide(container)

		sut, e := container.Get(ValidationContainerID)
		switch {
		case e != nil:
			t.Errorf("unexpected error (%v)", e)
		case sut == nil:
			t.Error("didn't returned a reference to service")
		default:
			switch sut.(type) {
			case Validator:
			default:
				t.Error("didn't returned the validator")
			}
		}
	})
}
