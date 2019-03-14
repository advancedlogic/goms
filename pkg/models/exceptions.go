package models

import "errors"

type Exception struct {
	errors []error
}

func (e Exception) Errors() []error {
	return e.errors
}

func (e Exception) Catch(msg string) {
	e.errors = append(e.errors, errors.New(msg))
}
