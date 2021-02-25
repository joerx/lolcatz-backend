package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/joerx/lolcatz-backend/db"
)

type healthHandler struct {
	db db.DB
}

type healthHandlerReponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Hostname string `json:"hostname"`
}

// Health creates a health check handler
func Health(db db.DB) http.HandlerFunc {
	h := &healthHandler{db}
	return h.handle
}

func (h healthHandler) getHealthStatus(ctx context.Context) (r *healthHandlerReponse, err error) {
	if err := h.db.Ping(ctx); err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &healthHandlerReponse{
		Status:   "ok",
		Database: "ok",
		Hostname: hostname,
	}, nil
}

func (h healthHandler) handle(w http.ResponseWriter, r *http.Request) {
	resp, err := h.getHealthStatus(r.Context())
	if err != nil {
		errorHandler(w, err)
		return
	}
	writeResponse(w, 200, resp)
}
