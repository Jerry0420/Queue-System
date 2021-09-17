package main

import (
    "fmt"
    "net/http"
    "log"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
    host     = "db"
    port     = 5432
    user     = "jerry"
    password = "jerry0315"
    dbname   = "queue_system"
  )

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello world~~")
}

func main() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        panic(err)
    }

    http.HandleFunc("/", hello)
    err = http.ListenAndServe("0.0.0.0:8000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}