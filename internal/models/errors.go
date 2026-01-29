package models

import "errors"

var (
	ErrNameRequired     = errors.New("service name is required")
	ErrEndpointRequired = errors.New("service endpoint is required")
	ErrServiceNotFound  = errors.New("service not found")
	ErrServiceExists    = errors.New("service already exists")
)
