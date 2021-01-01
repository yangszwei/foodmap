/* Package validator provide validation service for objects (object package) */
package validator

import (
	"foodmap/internal/infra/delivery"
	"foodmap/internal/infra/errors"
	"regexp"

	govalidator "github.com/go-playground/validator/v10"
)

// New create a Validator
func New() *Validator {
	v := new(Validator)
	v.v = govalidator.New()
	_ = v.v.RegisterValidation("hh-mm", func(fl govalidator.FieldLevel) bool {
		return regexp.MustCompile("^(0?[1-9]|1[012]):([0-5][0-9])[ap]m$").
			MatchString(fl.Field().String())
	})
	return v
}

// Validator wrapper of govalidator.Validate so the rest of the app don't need
// to know anything about it
type Validator struct {
	v *govalidator.Validate
}

// Validate validate struct and return errors.Errors, i should be the pointer
// of the instance
func (v *Validator) Validate(i interface{}) error {
	if e := v.v.Struct(i); e != nil {
		var errs errors.Errors
		for _, err := range e.(govalidator.ValidationErrors) {
			errs = append(errs, errors.NewValidationError(err.Field(), err.Tag()))
		}
		return errs
	}
	return nil
}

// IsFieldsValid returns a list of unrecognized fields
func IsFieldsValid(validTags []string, fields string) (result []string) {
	tags := delivery.Split(fields)
	if len(tags) == 0 {
		return nil
	}
	for _, tag := range tags {
		var exist bool
		if !exist {
			for _, t := range validTags {
				if tag == t {
					exist = true
					break
				}
			}
		}
		if !exist {
			result = append(result, tag)
		}
	}
	return
}
