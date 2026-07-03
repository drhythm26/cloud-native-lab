package main

import (
    "net/http"
    "encoding/json"
    "log"
)

type body struct {
    Server string `json:"server"`
    Status string `json:"status"`
}


func rootHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(body{
        Server: "go-api",
        Status: "0.1.0",
    })
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(body{
        Server: "healthz",
        Status: "ok",
    })
}

func main() {
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/healthz", healthzHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
