package utils

import (
	"encoding/json"
	"net/http"
	"bytes"
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

func JsonResponse(w http.ResponseWriter, response interface{}, serverError *ServerError) {
	w.Header().Set("Content-Type", "application/json")

	var code int
	var statusCode int

	if serverError != nil {
		code = serverError.Code
	} else {
		code = 20001
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