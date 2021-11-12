package presenter

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jerry0420/queue-system/backend/config"
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
	Code        int         `json:"code"`
	Data        interface{} `json:"data"`
	Description string      `json:"description,omitempty"`
}

func JsonResponseOK(w http.ResponseWriter, response interface{}) {
	JsonResponse(w, response, nil)
}

func JsonResponse(w http.ResponseWriter, response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")

	var code int
	var statusCode int
	var description string

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
		description = serverError.Description
	}

	statusCode = code / 100 //http sattus code is the first two digits of code.
	w.WriteHeader(statusCode)

	responseFormat := ResponseFormat{Code: code, Data: response}
	if config.ServerConfig.ENV() != config.EnvStatus.PROD {
		responseFormat.Description = description
	}

	json.NewEncoder(w).Encode(&responseFormat)
}
