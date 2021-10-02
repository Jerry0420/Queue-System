package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/jerry0420/queue-system/config"
)

func hello(w http.ResponseWriter, r *http.Request) {
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