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

func JsonResponseOK(w http.ResponseWriter, response interface{}) {
	JsonResponse(w, response, nil)
}

func JsonResponse(w http.ResponseWriter, response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		serverError, ok := err.(*domain.ServerError)
		if !ok {
			// Because err must be ServerError type.
			// Using panic to prevent any hidden error when develping...
			panic("err must be ServerError")
		}
		errorResponse := map[string]interface{}{"error_code": serverError.Code}
		if config.ServerConfig.ENV() != config.EnvStatus.PROD {
			errorResponse["description"] = serverError.Description
		}
		
		w.WriteHeader(serverError.Code / 100) //http status code is the first two digits of code.
		json.NewEncoder(w).Encode(&errorResponse)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&response)
	}
}
