package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/techops-interviews/service-registry/internal/models"
	"github.com/techops-interviews/service-registry/internal/store"
)

// Handler provides HTTP handlers for the service registry API.
type Handler struct {
	store *store.Store
}

// NewHandler creates a new API handler with the given store.
func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

// RegisterRoutes sets up the HTTP routes for the API.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /services", h.ListServices)
	mux.HandleFunc("POST /services", h.RegisterService)
	mux.HandleFunc("GET /services/{name}", h.GetService)
	mux.HandleFunc("PUT /services/{name}", h.UpdateService)
	mux.HandleFunc("DELETE /services/{name}", h.DeleteService)
	mux.HandleFunc("GET /services/{name}/health", h.CheckServiceHealth)
}

// Health returns the health status of the registry itself.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ListServices returns all registered services.
func (h *Handler) ListServices(w http.ResponseWriter, r *http.Request) {
	services := h.store.List()
	writeJSON(w, http.StatusOK, map[string]any{"services": services})
}

// RegisterService handles POST /services to register a new service.
func (h *Handler) RegisterService(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	service, err := h.store.Register(req)
	if err != nil {
		if errors.Is(err, models.ErrServiceExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to register service")
		return
	}

	writeJSON(w, http.StatusCreated, service)
}

// GetService handles GET /services/{name} to retrieve a service.
func (h *Handler) GetService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "service name is required")
		return
	}

	service, err := h.store.Get(name)
	if err != nil {
		if errors.Is(err, models.ErrServiceNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get service")
		return
	}

	writeJSON(w, http.StatusOK, service)
}

// UpdateService handles PUT /services/{name} to update a service.
func (h *Handler) UpdateService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "service name is required")
		return
	}

	var req models.UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	service, err := h.store.Update(name, req)
	if err != nil {
		if errors.Is(err, models.ErrServiceNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update service")
		return
	}

	writeJSON(w, http.StatusOK, service)
}

// DeleteService handles DELETE /services/{name} to deregister a service.
func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "service name is required")
		return
	}

	err := h.store.Delete(name)
	if err != nil {
		if errors.Is(err, models.ErrServiceNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete service")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CheckServiceHealth checks the health of a specific service by calling its health endpoint.
func (h *Handler) CheckServiceHealth(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "service name is required")
		return
	}

	service, err := h.store.Get(name)
	if err != nil {
		if errors.Is(err, models.ErrServiceNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get service")
		return
	}

	if service.HealthCheckURL == "" {
		writeJSON(w, http.StatusOK, map[string]any{
			"name":   service.Name,
			"status": models.StatusUnknown,
			"reason": "no health check URL configured",
		})
		return
	}

	status := checkHealth(service.HealthCheckURL)
	h.store.UpdateStatus(name, status)

	writeJSON(w, http.StatusOK, map[string]any{
		"name":   service.Name,
		"status": status,
	})
}

// checkHealth performs an HTTP GET to the health check URL and returns the status.
func checkHealth(url string) models.ServiceStatus {
	resp, err := http.Get(url)
	if err != nil {
		return models.StatusUnhealthy
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return models.StatusHealthy
	}
	return models.StatusUnhealthy
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// extractServiceName gets the service name from a URL path like /services/{name}
// This is a helper for older Go versions without PathValue support.
func extractServiceName(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
