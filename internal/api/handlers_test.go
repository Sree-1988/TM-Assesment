package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/techops-interviews/service-registry/internal/models"
	"github.com/techops-interviews/service-registry/internal/store"
)

func setupTestHandler() (*Handler, *http.ServeMux) {
	s := store.New()
	h := NewHandler(s)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return h, mux
}

func TestHealth(t *testing.T) {
	_, mux := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", resp["status"])
	}
}

func TestRegisterService(t *testing.T) {
	_, mux := setupTestHandler()

	body := `{"name": "my-service", "endpoint": "http://localhost:8080", "description": "Test service"}`
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var service models.Service
	json.NewDecoder(rec.Body).Decode(&service)

	if service.Name != "my-service" {
		t.Errorf("expected name 'my-service', got %q", service.Name)
	}
}

func TestRegisterService_MissingName(t *testing.T) {
	_, mux := setupTestHandler()

	body := `{"endpoint": "http://localhost:8080"}`
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestRegisterService_Duplicate(t *testing.T) {
	_, mux := setupTestHandler()

	body := `{"name": "my-service", "endpoint": "http://localhost:8080"}`

	// First registration
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Second registration (duplicate)
	req = httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
}

func TestListServices(t *testing.T) {
	_, mux := setupTestHandler()

	// Register a service first
	body := `{"name": "my-service", "endpoint": "http://localhost:8080"}`
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// List services
	req = httptest.NewRequest(http.MethodGet, "/services", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string][]models.Service
	json.NewDecoder(rec.Body).Decode(&resp)

	if len(resp["services"]) != 1 {
		t.Errorf("expected 1 service, got %d", len(resp["services"]))
	}
}

func TestGetService(t *testing.T) {
	_, mux := setupTestHandler()

	// Register a service first
	body := `{"name": "my-service", "endpoint": "http://localhost:8080"}`
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Get the service
	req = httptest.NewRequest(http.MethodGet, "/services/my-service", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var service models.Service
	json.NewDecoder(rec.Body).Decode(&service)

	if service.Name != "my-service" {
		t.Errorf("expected name 'my-service', got %q", service.Name)
	}
}

func TestGetService_NotFound(t *testing.T) {
	_, mux := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/services/nonexistent", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestDeleteService(t *testing.T) {
	_, mux := setupTestHandler()

	// Register a service first
	body := `{"name": "my-service", "endpoint": "http://localhost:8080"}`
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Delete the service
	req = httptest.NewRequest(http.MethodDelete, "/services/my-service", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	// Verify it's gone
	req = httptest.NewRequest(http.MethodGet, "/services/my-service", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d after delete, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdateService(t *testing.T) {
	_, mux := setupTestHandler()

	// Register a service first
	body := `{"name": "my-service", "endpoint": "http://localhost:8080"}`
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Update the service
	updateBody := `{"endpoint": "http://localhost:9090", "description": "Updated"}`
	req = httptest.NewRequest(http.MethodPut, "/services/my-service", bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var service models.Service
	json.NewDecoder(rec.Body).Decode(&service)

	if service.Endpoint != "http://localhost:9090" {
		t.Errorf("expected endpoint 'http://localhost:9090', got %q", service.Endpoint)
	}

	if service.Description != "Updated" {
		t.Errorf("expected description 'Updated', got %q", service.Description)
	}
}

func TestRegisterService_WithTags(t *testing.T) {
	_, mux := setupTestHandler()

	body := `{"name": "tagged-service", "endpoint": "http://localhost:8080", "tags": ["env:production", "team:platform"]}`
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var service models.Service
	json.NewDecoder(rec.Body).Decode(&service)

	if service.Tags["env"] != "production" {
		t.Errorf("expected tag env=production, got %q", service.Tags["env"])
	}
}

func TestListServices_WithTagFilter(t *testing.T) {
	_, mux := setupTestHandler()

	// Register services with tags
	body := `{"name": "prod-service", "endpoint": "http://localhost:8080", "tags": ["env:production"]}`
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	body = `{"name": "dev-service", "endpoint": "http://localhost:8081", "tags": ["env:development"]}`
	req = httptest.NewRequest(http.MethodPost, "/services", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// Filter by tag
	req = httptest.NewRequest(http.MethodGet, "/services?tag=env:production", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string][]*models.Service
	json.NewDecoder(rec.Body).Decode(&resp)

	if len(resp["services"]) != 2 {
		t.Errorf("expected 2 services, got %d", len(resp["services"]))
	}
}

// func TestListServices_MalformedTag(t *testing.T) {
// 	_, mux := setupTestHandler()
//
// 	// TODO: test malformed tag handling
// 	req := httptest.NewRequest(http.MethodGet, "/services?tag=invalid", nil)
// 	rec := httptest.NewRecorder()
// 	mux.ServeHTTP(rec, req)
// }
