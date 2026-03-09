package models

import "time"

// ServiceStatus represents the current health status of a service.
type ServiceStatus string

const (
	StatusUnknown   ServiceStatus = "unknown"
	StatusHealthy   ServiceStatus = "healthy"
	StatusUnhealthy ServiceStatus = "unhealthy"
)

// Service represents a registered service in the registry.
type Service struct {
	Name           string        `json:"name"`
	Endpoint       string        `json:"endpoint"`
	HealthCheckURL string        `json:"health_check_url,omitempty"`
	Status         ServiceStatus `json:"status"`
	Description    string        `json:"description,omitempty"`
	Tags           []string      `json:"tags,omitempty"`
	RegisteredAt   time.Time     `json:"registered_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// RegisterServiceRequest is the payload for registering a new service.
type RegisterServiceRequest struct {
	Name           string   `json:"name"`
	Endpoint       string   `json:"endpoint"`
	HealthCheckURL string   `json:"health_check_url,omitempty"`
	Description    string   `json:"description,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

// UpdateServiceRequest is the payload for updating an existing service.
type UpdateServiceRequest struct {
	Endpoint       string   `json:"endpoint,omitempty"`
	HealthCheckURL string   `json:"health_check_url,omitempty"`
	Description    string   `json:"description,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

// Validate checks that required fields are present for registration.
func (r RegisterServiceRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Endpoint == "" {
		return ErrEndpointRequired
	}
	// Validate tags if provided
	for _, tag := range r.Tags {
		if err := ValidateTag(tag); err != nil {
			return err
		}
	}
	return nil
}

// ValidateTag checks that a tag meets the requirements: non-empty, alphanumeric with hyphens, max 50 chars
func ValidateTag(tag string) error {
	if tag == "" {
		return ErrEmptyTag
	}
	if len(tag) > 50 {
		return ErrInvalidTag
	}
	// Validate alphanumeric and hyphens only
	for _, r := range tag {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
			return ErrInvalidTag
		}
	}
	return nil
}
