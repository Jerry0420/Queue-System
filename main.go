package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
    "io"
	"encoding/json"
	"github.com/jerry0420/queue-system/config"
	"github.com/jerry0420/queue-system/logging"
	"github.com/jerry0420/queue-system/utils"
)

func middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logger := logging.NewLogger([]string{"requestID", "duration", "code"}, false)        
        ctx := context.WithValue(r.Context(), "requestID", "aaaaaaaaaa")
        r = r.WithContext(ctx)
        
        responseWrapper := &utils.ResponseWrapper{
            ResponseWriter: w,
            Buffer: &bytes.Buffer{},
        }

        next(responseWrapper, r)
        
        var wrapperResponse *utils.ResponseFormat
        json.Unmarshal(responseWrapper.Buffer.Bytes(), &wrapperResponse)
        ctx = context.WithValue(r.Context(), "code", wrapperResponse.Code)
        ctx = context.WithValue(ctx, "duration", 3)
        
        io.Copy(w, responseWrapper.Buffer)
        r = r.WithContext(ctx)
        
        logger.INFOf(r.Context(), "hello world %d", 1234)
    }
}

func hello(w http.ResponseWriter, r *http.Request) {

    serverConfig := config.NewConfig()
    
    utils.JsonResponse(w, map[string]interface{}{
        "page": serverConfig.CONTEXT_TIMEOUT(),
        "list": []string{"hello", "world"},
    }, utils.ServerError40001)
}

func main() {
    http.HandleFunc("/", middleware(hello))
    err := http.ListenAndServe("0.0.0.0:8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}