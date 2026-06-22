package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	healthzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := `{"status":"ok"}` + "\n"
	if rec.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestReadyzHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	readyzHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := `{"status":"ready"}` + "\n"
	if rec.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestCreateReleaseHandler(t *testing.T) {
	body := strings.NewReader(`{
                "serviceName": "payment-api",
                "version": "v1.0.0",
                "environment": "dev",
                "owner": "lan"
        }`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/releases", body)
	rec := httptest.NewRecorder()

	createReleaseHandler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	responseBody := rec.Body.String()

	requiredParts := []string{
		`"serviceName":"payment-api"`,
		`"version":"v1.0.0"`,
		`"environment":"dev"`,
		`"status":"pending"`,
		`"owner":"lan"`,
		`"createdAt":`,
		`"updatedAt":`,
	}

	for _, part := range requiredParts {
		if !strings.Contains(responseBody, part) {
			t.Fatalf("expected response body to contain %q, got %s", part, responseBody)
		}
	}
}
