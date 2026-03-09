package models

import "errors"

var (
	ErrNameRequired     = errors.New("service name is required")
	ErrEndpointRequired = errors.New("service endpoint is required")
	ErrServiceNotFound  = errors.New("service not found")
	ErrServiceExists    = errors.New("service already exists")
	ErrInvalidTag       = errors.New("tag must be 1-50 characters, alphanumeric with hyphens only")
	ErrEmptyTag         = errors.New("tags cannot contain empty strings")
)
