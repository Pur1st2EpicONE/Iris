// Package v1 provides version 1 API handlers for the Iris application.
package v1

import (
	"Iris/internal/service"
)

// Handler represents the v1 API handler and holds references
// to the underlying service layer.
type Handler struct {
	service service.Service
}

// NewHandler creates a new v1 API Handler using the provided service.
func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}
