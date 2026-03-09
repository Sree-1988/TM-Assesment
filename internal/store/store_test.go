package store

import (
	"testing"

	"github.com/techops-interviews/service-registry/internal/models"
)

func TestStore_Register(t *testing.T) {
	s := New()

	req := models.RegisterServiceRequest{
		Name:        "test-service",
		Endpoint:    "http://localhost:8080",
		Description: "A test service",
	}

	service, err := s.Register(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if service.Name != req.Name {
		t.Errorf("expected name %q, got %q", req.Name, service.Name)
	}

	if service.Endpoint != req.Endpoint {
		t.Errorf("expected endpoint %q, got %q", req.Endpoint, service.Endpoint)
	}

	if service.Status != models.StatusUnknown {
		t.Errorf("expected status %q, got %q", models.StatusUnknown, service.Status)
	}
}

func TestStore_Register_Duplicate(t *testing.T) {
	s := New()

	req := models.RegisterServiceRequest{
		Name:     "test-service",
		Endpoint: "http://localhost:8080",
	}

	_, err := s.Register(req)
	if err != nil {
		t.Fatalf("unexpected error on first register: %v", err)
	}

	_, err = s.Register(req)
	if err != models.ErrServiceExists {
		t.Errorf("expected ErrServiceExists, got %v", err)
	}
}

func TestStore_Get(t *testing.T) {
	s := New()

	req := models.RegisterServiceRequest{
		Name:     "test-service",
		Endpoint: "http://localhost:8080",
	}
	s.Register(req)

	service, err := s.Get("test-service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if service.Name != "test-service" {
		t.Errorf("expected name %q, got %q", "test-service", service.Name)
	}
}

func TestStore_Get_NotFound(t *testing.T) {
	s := New()

	_, err := s.Get("nonexistent")
	if err != models.ErrServiceNotFound {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}
}

func TestStore_List(t *testing.T) {
	s := New()

	services := s.List()
	if len(services) != 0 {
		t.Errorf("expected empty list, got %d services", len(services))
	}

	s.Register(models.RegisterServiceRequest{Name: "svc1", Endpoint: "http://svc1"})
	s.Register(models.RegisterServiceRequest{Name: "svc2", Endpoint: "http://svc2"})

	services = s.List()
	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}
}

func TestStore_Update(t *testing.T) {
	s := New()

	s.Register(models.RegisterServiceRequest{
		Name:     "test-service",
		Endpoint: "http://localhost:8080",
	})

	updated, err := s.Update("test-service", models.UpdateServiceRequest{
		Endpoint:    "http://localhost:9090",
		Description: "Updated description",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.Endpoint != "http://localhost:9090" {
		t.Errorf("expected endpoint %q, got %q", "http://localhost:9090", updated.Endpoint)
	}

	if updated.Description != "Updated description" {
		t.Errorf("expected description %q, got %q", "Updated description", updated.Description)
	}
}

func TestStore_Update_NotFound(t *testing.T) {
	s := New()

	_, err := s.Update("nonexistent", models.UpdateServiceRequest{Endpoint: "http://new"})
	if err != models.ErrServiceNotFound {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}
}

func TestStore_Delete(t *testing.T) {
	s := New()

	s.Register(models.RegisterServiceRequest{Name: "test-service", Endpoint: "http://localhost:8080"})

	err := s.Delete("test-service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.Get("test-service")
	if err != models.ErrServiceNotFound {
		t.Errorf("expected service to be deleted, got %v", err)
	}
}

func TestStore_Delete_NotFound(t *testing.T) {
	s := New()

	err := s.Delete("nonexistent")
	if err != models.ErrServiceNotFound {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}
}

func TestStore_UpdateStatus(t *testing.T) {
	s := New()

	s.Register(models.RegisterServiceRequest{Name: "test-service", Endpoint: "http://localhost:8080"})

	err := s.UpdateStatus("test-service", models.StatusHealthy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	service, _ := s.Get("test-service")
	if service.Status != models.StatusHealthy {
		t.Errorf("expected status %q, got %q", models.StatusHealthy, service.Status)
	}
}

func TestStore_FilterByTag(t *testing.T) {
	s := New()

	// Register services with tags
	s.Register(models.RegisterServiceRequest{
		Name:     "api-service",
		Endpoint: "http://api:8080",
		Tags:     []string{"production", "api"},
	})

	s.Register(models.RegisterServiceRequest{
		Name:     "db-service",
		Endpoint: "http://db:5432",
		Tags:     []string{"production", "database"},
	})

	s.Register(models.RegisterServiceRequest{
		Name:     "cache-service",
		Endpoint: "http://cache:6379",
		Tags:     []string{"production", "cache"},
	})

	s.Register(models.RegisterServiceRequest{
		Name:     "dev-service",
		Endpoint: "http://dev:3000",
		Tags:     []string{"development"},
	})

	// Filter by "production" tag
	results := s.FilterByTag("production")
	if len(results) != 3 {
		t.Errorf("expected 3 services with 'production' tag, got %d", len(results))
	}

	// Filter by "api" tag
	results = s.FilterByTag("api")
	if len(results) != 1 {
		t.Errorf("expected 1 service with 'api' tag, got %d", len(results))
	}
	if results[0].Name != "api-service" {
		t.Errorf("expected 'api-service', got %q", results[0].Name)
	}

	// Filter by "development" tag
	results = s.FilterByTag("development")
	if len(results) != 1 {
		t.Errorf("expected 1 service with 'development' tag, got %d", len(results))
	}

	// Filter by non-existent tag
	results = s.FilterByTag("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected 0 services with non-existent tag, got %d", len(results))
	}
}

func TestValidateTag(t *testing.T) {
	testCases := []struct {
		tag   string
		valid bool
	}{
		// Valid tags
		{"production", true},
		{"dev", true},
		{"api-service", true},
		{"python-3-11", true},
		{"a", true},
		{"a1b2c3", true},
		{"hyphen-test-123", true},

		// Invalid tags
		{"", false},                    // empty
		{"invalid tag", false},         // space
		{"invalid.tag", false},         // dot
		{"invalid_tag", false},         // underscore
		{"invalid@tag", false},         // special char
		{string(make([]byte, 51)), false}, // too long (51 chars)
	}

	for _, tc := range testCases {
		err := models.ValidateTag(tc.tag)
		if tc.valid && err != nil {
			t.Errorf("expected %q to be valid, got error: %v", tc.tag, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("expected %q to be invalid, but validation passed", tc.tag)
		}
	}
}
