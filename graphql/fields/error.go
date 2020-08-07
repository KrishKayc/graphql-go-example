package fields

import "fmt"

//InvalidArgumentError resturned for an invalid arg type
type InvalidArgumentError struct {
	arg          string
	expectedType string
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("Invalid Type for argument '%s'. Expect type '%s'", e.arg, e.expectedType)
}

//RequiredArgumentError is returned when required arg is not present
type RequiredArgumentError struct {
	arg string
}

func (e *RequiredArgumentError) Error() string {
	return fmt.Sprintf("Argument '%s' is required and is missing.", e.arg)
}
