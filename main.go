package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jerry0420/queue-system/config"
	"github.com/jerry0420/queue-system/logging"
	"github.com/jerry0420/queue-system/utils"
)

func hello(w http.ResponseWriter, r *http.Request) {
    // logger := logging.NewLogger([]string{"requestID", "duration"}, false)
    ctx := context.WithValue(r.Context(), "requestID", "aaaaaaaaaa")
    ctx = context.WithValue(ctx, "duration", 3)
    // logger.INFOf(ctx, "hello world %d", 1234)

    serverConfig := config.NewConfig()
    
    utils.JsonResponse(w, ctx, map[string]interface{}{
        "page": serverConfig.CONTEXT_TIMEOUT(),
        "list": []string{"hello", "world"},
    }, utils.ServerError40001)

    logger := logging.NewLogger([]string{"requestID", "duration", "code"}, false)
    logger.INFOf(ctx, "hello world %d", 1234)
}

func main() {
    http.HandleFunc("/", hello)
    err := http.ListenAndServe("0.0.0.0:8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}