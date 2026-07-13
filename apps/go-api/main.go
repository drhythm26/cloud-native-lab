package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type statusRecoder struct {
	http.ResponseWriter
	code int
}

func (r *statusRecoder) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_api_http_requests_total",
			Help: "Total HTTP requests by path",
		},
		[]string{"path", "method", "code"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_api_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpDuration)
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
		start := time.Now()
		next(rw, r)
		duration := time.Since(start).Seconds()
		code := strconv.Itoa(rw.code)
		httpRequests.WithLabelValues(
			path,
			r.Method,
			code,
		).Inc()
		httpDuration.WithLabelValues(path).Observe(duration)
	}
}

func main() {
	http.HandleFunc("/healthz", withMetrics("/healthz", healthzHandler))
	http.HandleFunc("/readyz", withMetrics("/readyz", readyzHandler))
	http.HandleFunc("/env", withMetrics("/env", envHandler))
	http.HandleFunc("/file", withMetrics("/file", fileHandler))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
