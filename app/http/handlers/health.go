package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/joerx/lolcatz-backend/db"
)

// HealthHandler is a http handler that returns application health information
type HealthHandler struct {
	db db.DB
}

type healthCheckReponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Hostname string `json:"hostname"`
}

// NewHealth creates a health check handler
func NewHealth(db db.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) checkHealth(ctx context.Context) (r *healthCheckReponse, err error) {
	if err := h.db.Ping(ctx); err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &healthCheckReponse{
		Status:   "ok",
		Database: "ok",
		Hostname: hostname,
	}, nil
}

// Health performs an application health check and returns an appropriate response
// Response status code will be 200 if health check passes, 500 otherwise
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	resp, err := h.checkHealth(r.Context())
	if err != nil {
		errorHandler(w, err)
		return
	}
	writeResponseJSON(w, 200, resp)
}
