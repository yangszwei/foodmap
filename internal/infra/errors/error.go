package errors

import (
	"fmt"
	"strings"
)

// Error object
type Error struct {
	Name    string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Name, e.Message)
}

// New create a error
func New(error string, message ...string) error {
	var e Error
	e.Name = error
	if len(message) > 0 {
		e.Message = strings.Join(message, ",")
	}
	return e
}
