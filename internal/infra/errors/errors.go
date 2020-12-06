/*
	Package error provide a custom error implementation so errors can be passed
    across packages easier
*/
package errors

import "fmt"

// Errors a list of errors
type Errors []error

// Error print all errors separated by line breaks
func (e Errors) Error() (message string) {
	for _, err := range e {
		message += fmt.Sprintln(err.Error())
	}
	return
}
