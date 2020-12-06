package validator_test

import (
	"foodmap/internal/infra/errors"
	"foodmap/internal/infra/validator"
	"testing"
)

type Test struct {
	A string `validate:"min=2,max=2"`
	B string `validate:"min=3,max=3"`
	C string `validate:"min=4,max=4"`
}

func TestValidator_Validate(t *testing.T) {
	v := validator.New()
	t.Run("should succeed", func(t *testing.T) {
		err := v.Validate(&Test{
			A: "--",
			B: "---",
			C: "----",
		})
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("should fail", func(t *testing.T) {
		err := v.Validate(&Test{
			A: "--",
			B: "--",
			C: "----",
		})
		if err == nil || err.(errors.Errors)[0].Error() != "Validation failed on field \"B\" with tag \"min\"" {
			t.Error(err)
		}
	})
}
