package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type statusRecoder struct {
	http.ResponseWriter
	code int
}

func (r *statusRecoder) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

var httpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "go_api_http_requests_total",
		Help: "Total HTTP requests by path",
	},
	[]string{"path", "method", "code"},
)

func init() {
	prometheus.MustRegister(httpRequests)
}

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

func withMetrics(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &statusRecoder{ResponseWriter: w, code: 200}
		next(rw, r)
		httpRequests.WithLabelValues(
			path,
			r.Method,
			strconv.Itoa(rw.code),
		).Inc()
	}
}

func main() {
	http.HandleFunc("/healthz", withMetrics("/healthz", healthzHandler))
	http.HandleFunc("/readyz", withMetrics("/readyz",readyzHandler))
	http.HandleFunc("/env", withMetrics("/env", envHandler))
	http.HandleFunc("/file", withMetrics("/file", fileHandler))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
