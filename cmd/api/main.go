package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func loadDBConfig() DBConfig {
	return DBConfig{
		Host:     getenv("DB_HOST", "localhost"),
		Port:     getenv("DB_PORT", "5432"),
		Name:     getenv("DB_NAME", "release_tracker"),
		User:     getenv("DB_USER", "release_tracker"),
		Password: getenv("DB_PASSWORD", "release_tracker"),
		SSLMode:  getenv("DB_SSLMODE", "disable"),
	}
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Name,
		c.User,
		c.Password,
		c.SSLMode,
	)
}

func openDB(ctx context.Context, cfg DBConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

type App struct {
	db *sql.DB
}

func newApp(db *sql.DB) *App {
	return &App{
		db: db,
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response failed: %v", err)
	}
}

func (a *App) healthzHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (a *App) readyzHandler(w http.ResponseWriter, r *http.Request) {
	if a.db == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := a.db.PingContext(ctx); err != nil {
		log.Printf("database readiness check failed: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func (a *App) createReleaseHandler(w http.ResponseWriter, r *http.Request) {
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

	if a.db == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "database is not ready",
		})
		return
	}

	now := time.Now().UTC()
	release := Release{
		ID:          fmt.Sprintf("rel_%d", now.UnixNano()),
		ServiceName: req.ServiceName,
		Version:     req.Version,
		Environment: req.Environment,
		Status:      "pending",
		Owner:       req.Owner,
		CreatedAt:   now.Format(time.RFC3339),
		UpdatedAt:   now.Format(time.RFC3339),
	}

	_, err := a.db.ExecContext(
		r.Context(),
		`INSERT INTO releases (
			id,
			service_name,
			version,
			environment,
			status,
			owner,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		release.ID,
		release.ServiceName,
		release.Version,
		release.Environment,
		release.Status,
		release.Owner,
		now,
		now,
	)
	if err != nil {
		log.Printf("insert release failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create release",
		})
		return
	}

	writeJSON(w, http.StatusCreated, release)
}

func (a *App) listReleasesHandler(w http.ResponseWriter, r *http.Request) {
	if a.db == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "database is not ready",
		})
		return
	}

	rows, err := a.db.QueryContext(
		r.Context(),
		`SELECT
			id,
			service_name,
			version,
			environment,
			status,
			owner,
			created_at,
			updated_at
		FROM releases
		ORDER BY created_at ASC`,
	)
	if err != nil {
		log.Printf("list releases failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to list releases",
		})
		return
	}
	defer rows.Close()

	items := []Release{}
	for rows.Next() {
		var release Release
		var createdAt time.Time
		var updatedAt time.Time

		if err := rows.Scan(
			&release.ID,
			&release.ServiceName,
			&release.Version,
			&release.Environment,
			&release.Status,
			&release.Owner,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Printf("scan release failed: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to list releases",
			})
			return
		}

		release.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		release.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)
		items = append(items, release)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iterate releases failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to list releases",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string][]Release{
		"items": items,
	})
}

func (a *App) getReleaseHandler(w http.ResponseWriter, r *http.Request) {
	if a.db == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "database is not ready",
		})
		return
	}

	id := r.PathValue("id")

	var release Release
	var createdAt time.Time
	var updatedAt time.Time

	err := a.db.QueryRowContext(
		r.Context(),
		`SELECT
			id,
			service_name,
			version,
			environment,
			status,
			owner,
			created_at,
			updated_at
		FROM releases
		WHERE id = $1`,
		id,
	).Scan(
		&release.ID,
		&release.ServiceName,
		&release.Version,
		&release.Environment,
		&release.Status,
		&release.Owner,
		&createdAt,
		&updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "release not found",
		})
		return
	}
	if err != nil {
		log.Printf("get release failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to get release",
		})
		return
	}

	release.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	release.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)

	writeJSON(w, http.StatusOK, release)
}

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", a.healthzHandler)
	mux.HandleFunc("GET /readyz", a.readyzHandler)
	mux.HandleFunc("POST /api/v1/releases", a.createReleaseHandler)
	mux.HandleFunc("GET /api/v1/releases", a.listReleasesHandler)
	mux.HandleFunc("GET /api/v1/releases/{id}", a.getReleaseHandler)
	mux.Handle("GET /metrics", promhttp.Handler())
	return mux
}

func main() {
	port := getenv("PORT", "8080")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := openDB(ctx, loadDBConfig())
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}
	defer db.Close()

	app := newApp(db)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           app.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("release tracker api listening on :%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
