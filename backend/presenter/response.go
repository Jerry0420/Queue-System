package presenter

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jerry0420/queue-system/backend/domain"
)

type ResponseWrapper struct {
	http.ResponseWriter
	Buffer *bytes.Buffer
}

func (responseWrapper *ResponseWrapper) Write(p []byte) (int, error) {
	return responseWrapper.Buffer.Write(p)
}

type ResponseFormat struct {
	Code interface{} `json:"code"`
	Data interface{} `json:"data"`
}

func JsonResponseOK(w http.ResponseWriter, response interface{}) {
	JsonResponse(w, response, nil)
}

func JsonResponse(w http.ResponseWriter, response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")

	var code int
	var statusCode int

	if err == nil {
		code = 20001
	} else {
		serverError, ok := err.(*domain.ServerError)
		if !ok {
			// Because err must be ServerError type.
			// Using panic to prevent any hidden error when develping...
			panic("err must be ServerError")
		}
		code = serverError.Code
	}

	statusCode = code / 100 //http sattus code is the first two digits of code.
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(
		&ResponseFormat{
			Code: code,
			Data: response,
		},
	)
}