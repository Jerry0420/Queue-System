package utils

import (
	"errors"
	"fmt"
)

type ServerError struct {
    Code int
    BaseError error
}

func (serverError *ServerError) Error() string {
    return fmt.Sprintf("%v", serverError.BaseError)
}

var (
    // list of all custom errors.
    // check README file for detailed description.
    ServerError40001 *ServerError
)

func init() {
    ServerError40001 = &ServerError{Code: 40001, BaseError: errors.New("some reason...")}
}