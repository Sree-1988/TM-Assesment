package store

import (
	"sync"
	"time"

	"github.com/techops-interviews/service-registry/internal/models"
)

// Store provides thread-safe storage for registered services.
type Store struct {
	mu       sync.RWMutex
	services map[string]*models.Service
}

// New creates a new empty store.
func New() *Store {
	return &Store{
		services: make(map[string]*models.Service),
	}
}

// Register adds a new service to the store.
func (s *Store) Register(req models.RegisterServiceRequest) (*models.Service, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[req.Name]; exists {
		return nil, models.ErrServiceExists
	}

	now := time.Now()
	service := &models.Service{
		Name:           req.Name,
		Endpoint:       req.Endpoint,
		HealthCheckURL: req.HealthCheckURL,
		Description:    req.Description,
		Status:         models.StatusUnknown,
		RegisteredAt:   now,
		UpdatedAt:      now,
	}

	s.services[req.Name] = service
	return service, nil
}

// Get retrieves a service by name.
func (s *Store) Get(name string) (*models.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	service, exists := s.services[name]
	if !exists {
		return nil, models.ErrServiceNotFound
	}
	return service, nil
}

// List returns all registered services.
func (s *Store) List() []*models.Service {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.Service, 0, len(s.services))
	for _, svc := range s.services {
		result = append(result, svc)
	}
	return result
}

// Update modifies an existing service.
func (s *Store) Update(name string, req models.UpdateServiceRequest) (*models.Service, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	service, exists := s.services[name]
	if !exists {
		return nil, models.ErrServiceNotFound
	}

	if req.Endpoint != "" {
		service.Endpoint = req.Endpoint
	}
	if req.HealthCheckURL != "" {
		service.HealthCheckURL = req.HealthCheckURL
	}
	if req.Description != "" {
		service.Description = req.Description
	}
	service.UpdatedAt = time.Now()

	return service, nil
}

// Delete removes a service from the store.
func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.services[name]; !exists {
		return models.ErrServiceNotFound
	}

	delete(s.services, name)
	return nil
}

// UpdateStatus updates the health status of a service.
func (s *Store) UpdateStatus(name string, status models.ServiceStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	service, exists := s.services[name]
	if !exists {
		return models.ErrServiceNotFound
	}

	service.Status = status
	service.UpdatedAt = time.Now()
	return nil
}
