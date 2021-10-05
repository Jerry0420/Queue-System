package utils

import (
	"context"
	"encoding/json"
	"net/http"
)

func JsonResponseOK(w http.ResponseWriter, ctx context.Context, response interface{}) {
	JsonResponse(w, ctx, response, nil)
}

func JsonResponse(w http.ResponseWriter, ctx context.Context, response interface{}, serverError *ServerError) {
	w.Header().Set("Content-Type", "application/json")

	var code int
	var statusCode int

	if serverError != nil {
		code = serverError.Code
	} else {
		code = 20001
	}
	
	ctx = context.WithValue(ctx, "code", code)


	statusCode = code / 100 //http sattus code is the first two digits of code.
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(
		&map[string]interface{}{
			"data": response,
			"code": code,
		},
	)
}