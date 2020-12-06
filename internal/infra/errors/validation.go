package errors

import "fmt"

// NewValidationError create a ValidationError
func NewValidationError(field, tag string) (err ValidationError) {
	err.Field = field
	err.Tag = tag
	return
}

// ValidationErr validation error
type ValidationError struct {
	Field string
	Tag   string
}

// Error print the error as string
func (v ValidationError) Error() string {
	f := "Validation failed on field \"%s\" with tag \"%s\""
	return fmt.Sprintf(f, v.Field, v.Tag)
}
