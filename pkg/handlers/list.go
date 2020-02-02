package handlers

import (
	"net/http"

	"github.com/joerx/lolcatz-backend/pkg/db"
)

type listHandler struct {
	db *db.Client
}

// ListUploads returns a handler that allows listing of uploaded images
func ListUploads(db *db.Client) http.HandlerFunc {
	h := &listHandler{db: db}
	return h.handle
}

func (h *listHandler) handle(w http.ResponseWriter, r *http.Request) {
	uploads, err := h.db.ListUploads("johndoe")
	if err != nil {
		errorHandler(w, err)
	}
	writeReponse(w, http.StatusOK, uploads)
}
