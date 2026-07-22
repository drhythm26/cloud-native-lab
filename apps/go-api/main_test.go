package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		path       string
		wantCode   int
		wantStatus string
	}{
		{"healthz", healthzHandler, "/healthz", http.StatusOK, "ok"},
		{"readyz", readyzHandler, "/readyz", http.StatusOK, "ok"},
		{"not found", notFoundHandler, "/nope", http.StatusNotFound, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			tt.handler(rr, req)
			if rr.Code != tt.wantCode {
				t.Fatalf("%s: status = %d, want %d", tt.name, rr.Code, tt.wantCode)
			}
			if tt.wantStatus == "" {
				return
			}
			var body map[string]string
			if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
				t.Fatalf("%s: decode body: %v", tt.name, err)
			}
			if body["status"] != tt.wantStatus {
				t.Fatalf("%s: status failed = %q, want %q", tt.name, body["status"], tt.wantStatus)
			}
		})
	}
}

func TestFileHandler(t *testing.T) {
	name := "file"
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir,"a.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("go-api"), 0644)
	handler := fileHandler(dir)
	req := httptest.NewRequest(http.MethodGet, "/file", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("%s: status = %d, want = 200", name, rr.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("%s: decode body: %v", name, err)
	}
	if body["b.txt"] != "go-api" {
		t.Fatalf("%s: failed = %q, want = %q", name, body["b.txt"], "go-api")
	}
}

func TestEnv(t *testing.T) {
	testName := "env"
	data, err := os.ReadFile("./env")
	if err != nil {
		t.Fatalf("%s: failed read env", testName)
	}
	names := strings.Split(string(data), "\n")
	values := []string{"test", "debug", "secret123"}
	for i, name := range names {
		os.Setenv(name, values[i])
	}
	handler := envHandler("./env")
	req := httptest.NewRequest(http.MethodGet, "/env", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("%s: status = %d, want = 200", testName, rr.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("%s: decode body: %v", testName, err)
	}
	for i, name := range names {
		if body[name] != values[i] {
			t.Fatalf("%s: env = %q, failed = %q, want = %q", testName, name, body[name], values[i])
		}
	}
}