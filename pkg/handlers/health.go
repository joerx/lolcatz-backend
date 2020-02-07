package handlers

import (
	"net/http"
	"time"

	"github.com/joerx/lolcatz-backend/pkg/db"
)

type healthHandler struct {
	db *db.Client
}

type healthHandlerReponse struct {
	Database bool `json:"database"`
}

// Health creates a health check handler
func Health(db *db.Client) http.HandlerFunc {
	h := &healthHandler{db}
	return h.handle
}

func (h healthHandler) handle(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(5 * time.Second); err != nil {
		errorHandler(w, err)
	}
	writeResponse(w, 200, healthHandlerReponse{Database: true})
}
