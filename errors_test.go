package sapi

import (
	"errors"
	"testing"

	"github.com/happyhippyhippo/slate"
)

func Test_err(t *testing.T) {
	t.Run("errNilPointer", func(t *testing.T) {
		arg := "dummy argument"
		context := map[string]interface{}{"field": "value"}
		message := "dummy argument : invalid nil pointer"

		t.Run("creation without context", func(t *testing.T) {
			if e := errNilPointer(arg); !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("error not a instance of ErrNilPointer")
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
			if e := errNilPointer(arg, context); !errors.Is(e, slate.ErrNilPointer) {
				t.Errorf("error not a instance of ErrNilPointer")
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

	t.Run("errConversion", func(t *testing.T) {
		arg := "dummy value"
		typ := "dummy type"
		context := map[string]interface{}{"field": "value"}
		message := "dummy value to dummy type : invalid type conversion"

		t.Run("creation without context", func(t *testing.T) {
			if e := errConversion(arg, typ); !errors.Is(e, slate.ErrConversion) {
				t.Errorf("error not a instance of ErrConversion")
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
			if e := errConversion(arg, typ, context); !errors.Is(e, slate.ErrConversion) {
				t.Errorf("error not a instance of ErrConversion")
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
