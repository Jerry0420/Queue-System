package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/jerry0420/queue-system/config"
	"github.com/jerry0420/queue-system/logging"
)

func hello(w http.ResponseWriter, r *http.Request) {
    logger := logging.NewLogger([]string{"requestID", "duration"}, false)
    ctx := context.WithValue(r.Context(), "requestID", "aaaaaaaaaa")
    ctx = context.WithValue(ctx, "duration", 3)
    logger.INFOf(ctx, "hello world %d", 1234)

    serverConfig := config.NewConfig()
    fmt.Fprintf(w, serverConfig.ENV())
}

func main() {
    http.HandleFunc("/", hello)
    err := http.ListenAndServe("0.0.0.0:8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}