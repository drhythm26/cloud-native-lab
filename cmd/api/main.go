package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Release struct {
	ID          string `json:"id"`
	ServiceName string `json:"serviceName"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
	Owner       string `json:"owner"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type CreateReleaseRequest struct {
	ServiceName string `json:"serviceName"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Owner       string `json:"owner"`
}

var (
	releasesMu sync.Mutex
	releases   = map[string]Release{}
)

func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response failed: %v", err)
	}
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func createReleaseHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateReleaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json body",
		})
		return
	}
	if strings.TrimSpace(req.ServiceName) == "" ||
		strings.TrimSpace(req.Version) == "" ||
		strings.TrimSpace(req.Environment) == "" ||
		strings.TrimSpace(req.Owner) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "serviceName, version, environment, and owner are required",
		})
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	id := "rel_" + strings.ReplaceAll(now, ":", "")
	release := Release{
		ID:          id,
		ServiceName: req.ServiceName,
		Version:     req.Version,
		Environment: req.Environment,
		Status:      "pending",
		Owner:       req.Owner,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	releasesMu.Lock()
	releases[id] = release
	releasesMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(release); err != nil {
		log.Printf("write response failed: %v", err)
	}
}

func listReleasesHandler(w http.ResponseWriter, r *http.Request) {
	releasesMu.Lock()
	items := make([]Release, 0, len(releases))
	for _, release := range releases {
		items = append(items, release)
	}
	releasesMu.Unlock()
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt < items[j].CreatedAt
	})
	writeJSON(w, http.StatusOK, map[string][]Release{
		"items": items,
	})
}

func getReleaseHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	releasesMu.Lock()
	release, ok := releases[id]
	releasesMu.Unlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "release not found",
		})
		return
	}
	writeJSON(w, http.StatusOK, release)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /readyz", readyzHandler)
	mux.HandleFunc("POST /api/v1/releases", createReleaseHandler)
	mux.HandleFunc("GET /api/v1/releases", listReleasesHandler)
	mux.HandleFunc("GET /api/v1/releases/{id}", getReleaseHandler)
	mux.Handle("GET /metrics", promhttp.Handler())
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("release tracker api listening on :%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
