package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/joerx/lolcatz-backend/pkg/db"
)

type healthHandler struct {
	db *db.Client
}

type healthHandlerReponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Hostname string `json:"hostname"`
}

// Health creates a health check handler
func Health(db *db.Client) http.HandlerFunc {
	h := &healthHandler{db}
	return h.handle
}

func (h healthHandler) getHealthStatus() (*healthHandlerReponse, error) {
	if err := h.db.Ping(5 * time.Second); err != nil {
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
	resp, err := h.getHealthStatus()
	if err != nil {
		errorHandler(w, err)
		return
	}
	writeResponse(w, 200, resp)
}
