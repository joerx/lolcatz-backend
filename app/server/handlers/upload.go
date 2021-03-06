package handlers

import (
	"log"
	"net/http"

	"github.com/joerx/lolcatz-backend/http/errors"
	"github.com/joerx/lolcatz-backend/http/mp"
	"github.com/joerx/lolcatz-backend/s3"
	"github.com/joerx/lolcatz-backend/upload"
)

// UploadHandler handles user uploads
type UploadHandler struct {
	s3cfg   s3.Config
	uploads upload.Service
}

// NewUpload create a new UploadHandler
func NewUpload(s3cfg s3.Config, uploads upload.Service) *UploadHandler {
	return &UploadHandler{s3cfg: s3cfg, uploads: uploads}
}

// CreateUpload accepts user submitted uploads, stores the files in S3 and registers the uploads
// in the application database
func (h *UploadHandler) CreateUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorHandler(w, errors.MethodNotAllowed(r.Method))
		return
	}

	u, err := h.createUpload(w, r)
	if err != nil {
		log.Printf("Error processing upload request - %v", err)
		errorHandler(w, err)
		return
	}

	writeResponseJSON(w, http.StatusOK, u)
}

func (h *UploadHandler) createUpload(w http.ResponseWriter, r *http.Request) (*upload.Upload, error) {
	uf, err := mp.GetUploadedFile(r, "file")
	if err != nil {
		return nil, err
	}

	in := s3.UploadRequest{
		Filename:     uf.TmpFile.Name(),
		OriginalName: uf.Filename,
		ContentType:  uf.ContentType,
	}

	s3Key, err := s3.Upload(in, h.s3cfg)
	if err != nil {
		return nil, err
	}

	u := &upload.Upload{
		Username: "johndoe",
		Filename: uf.Filename,
		S3Url:    s3Key,
	}

	u, err = h.uploads.CreateUpload(r.Context(), u)
	if err != nil {
		return nil, err
	}

	if err := h.presignURL(u); err != nil {
		return nil, err
	}

	return u, nil
}

// FindUploads finds user submitted uploads. Only uploads for the current user will be returned
func (h *UploadHandler) FindUploads(w http.ResponseWriter, r *http.Request) {
	username := "johndoe" // todo: get user from context
	filter := &upload.Filter{Username: &username}

	result, err := h.uploads.FindUploads(r.Context(), filter)
	if err != nil {
		errorHandler(w, err)
		return
	}

	// replace S3 urls with pre-signed ones
	for _, u := range result {
		if err = h.presignURL(u); err != nil {
			errorHandler(w, err)
			return
		}
	}

	writeResponseJSON(w, http.StatusOK, result)
}

func (h *UploadHandler) presignURL(u *upload.Upload) error {
	urlStr, err := s3.Presign(u.S3Url, h.s3cfg)
	if err != nil {
		return err
	}
	u.S3Url = urlStr
	return nil
}
