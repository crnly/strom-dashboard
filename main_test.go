package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDashboardContainsBothViews(t *testing.T) {
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	for _, text := range []string{`id="dashboard"`, `id="system"`, "Beispieldaten", "diagrams/system.svg"} {
		if !strings.Contains(response.Body.String(), text) {
			t.Errorf("dashboard response does not contain %q", text)
		}
	}
}

func TestPublicAssetsAndPrivatePaths(t *testing.T) {
	handler := newHandler()
	for _, tc := range []struct {
		path string
		want int
	}{
		{"/styles.css", http.StatusOK},
		{"/diagrams/system.svg", http.StatusOK},
		{"/diagrams/dataflow.svg", http.StatusOK},
		{"/diagrams/database.svg", http.StatusOK},
		{"/diagrams/deployment.svg", http.StatusOK},
		{"/diagrams/system.d2", http.StatusOK},
		{"/diagrams/", http.StatusNotFound},
		{"/diagrams/missing.svg", http.StatusNotFound},
		{"/main.go", http.StatusNotFound},
		{"/go.mod", http.StatusNotFound},
		{"/README.md", http.StatusNotFound},
		{"/.env", http.StatusNotFound},
	} {
		t.Run(tc.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if response.Code != tc.want {
				t.Errorf("status = %d, want %d", response.Code, tc.want)
			}
		})
	}
}

func TestHealthReportsWebProcess(t *testing.T) {
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected response: status=%d headers=%v", response.Code, response.Header())
	}
	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil || result.Status != "ok" {
		t.Fatalf("invalid health response: %s (error: %v)", response.Body.String(), err)
	}
}

func TestReadOnlyRoutesRejectWrites(t *testing.T) {
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/", nil))
	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", response.Code)
	}
}
