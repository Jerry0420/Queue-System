package main

import (
    "fmt"
    "net/http"
    "log"
    "github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`./config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
    env := viper.GetString(`env`)
    fmt.Fprintf(w, env)
}

func main() {
    http.HandleFunc("/", hello)
    err := http.ListenAndServe("0.0.0.0:8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}