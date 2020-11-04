package handlers

import (
	"net/http"

	"github.com/joerx/lolcatz-backend/pkg/db"
	"github.com/joerx/lolcatz-backend/pkg/s3"
)

type listHandler struct {
	db    *db.Client
	s3cfg s3.Config
}

// ListUploads returns a handler that allows listing of uploaded images
func ListUploads(s3cfg s3.Config, db *db.Client) http.HandlerFunc {
	h := &listHandler{s3cfg: s3cfg, db: db}
	return h.handle
}

func (h *listHandler) handle(w http.ResponseWriter, r *http.Request) {
	uploads, err := h.db.ListUploads("johndoe")
	if err != nil {
		errorHandler(w, err)
		return
	}

	result := make([]db.Upload, len(uploads))

	for i, u := range uploads {
		urlStr, err := s3.Presign(u.S3Url, h.s3cfg)
		if err != nil {
			errorHandler(w, err)
			return
		}
		u.S3Url = urlStr
		result[i] = u
	}

	writeResponse(w, http.StatusOK, result)
}
