package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	out := map[string]string{}
	out["APP_ENV"] = os.Getenv("APP_ENV")
	out["LOG_LEVEL"] = os.Getenv("LOG_LEVEL")
	out["API_TOKEN"] = os.Getenv("API_TOKEN")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir("/etc/go-api")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	files := map[string]string{}
	for _, e := range entries {
		path := "/etc/go-api/" + e.Name()
		info, err := os.Stat(path)
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			files[e.Name()] = "read error: " + err.Error()
			continue
		}
		files[e.Name()] = string(data)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func main() {
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/readyz", readyzHandler)
	http.HandleFunc("/env", envHandler)
	http.HandleFunc("/file", fileHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
