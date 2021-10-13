package domain

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
    ServerError40401 *ServerError
)

func init() {
    ServerError40401 = &ServerError{Code: 40401, BaseError: errors.New("unsupported url route")}
}