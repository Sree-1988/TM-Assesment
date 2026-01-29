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
