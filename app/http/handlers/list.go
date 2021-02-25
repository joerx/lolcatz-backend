package handlers

import (
	"net/http"

	"github.com/joerx/lolcatz-backend"
	"github.com/joerx/lolcatz-backend/s3"
)

type listHandler struct {
	uploads lolcatz.UploadService
	s3cfg   s3.Config
}

// ListUploads returns a handler that allows listing of uploaded images
func ListUploads(s3cfg s3.Config, uploads lolcatz.UploadService) http.HandlerFunc {
	h := &listHandler{s3cfg: s3cfg, uploads: uploads}
	return h.handle
}

func (h *listHandler) handle(w http.ResponseWriter, r *http.Request) {
	username := "johndoe"
	filter := &lolcatz.UploadFilter{Username: &username}

	result, err := h.uploads.FindUploads(r.Context(), filter)
	if err != nil {
		errorHandler(w, err)
		return
	}

	// replace S3 urls with pre-signed ones
	for _, u := range result {
		urlStr, err := s3.Presign(u.S3Url, h.s3cfg)
		if err != nil {
			errorHandler(w, err)
			return
		}
		u.S3Url = urlStr
	}

	writeResponse(w, http.StatusOK, result)
}
