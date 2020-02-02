package handlers

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joerx/lolcatz-backend/pkg/db"
	"github.com/joerx/lolcatz-backend/pkg/s3"
)

// Upload returns a handler that handles file uploads to the given region and bucket
func Upload(s3cfg s3.Config, db *db.Client) http.HandlerFunc {
	h := &uploadHandler{s3cfg: s3cfg, db: db}
	return h.handle
}

type uploadHandler struct {
	s3cfg s3.Config
	db    *db.Client
}

func (h *uploadHandler) handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s /upload", r.Method)

	if r.Method != "POST" {
		writeResponseMsg(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeResponseMsg(w, http.StatusNotAcceptable, "Can't parse this as multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
	}

	tmpFile, err := downloadTempFile(file, header)
	if err != nil {
		errorHandler(w, err)
		return
	}

	in := s3.UploadRequest{
		Filename:     tmpFile.Name(),
		OriginalName: header.Filename,
		ContentType:  header.Header.Get("Content-Type"),
	}

	s3Key, err := s3.Upload(in, h.s3cfg)
	if err != nil {
		errorHandler(w, err)
		return
	}

	u := db.Upload{
		Username: "johndoe",
		Filename: header.Filename,
		S3Url:    s3Key,
	}
	if err := h.db.InsertUpload(u); err != nil {
		errorHandler(w, err)
	}

	log.Println("Upload recorded in database")

	writeResponse(w, http.StatusOK, db.Upload{
		ID:       -1,
		Filename: header.Filename,
		S3Url:    "", // empty url before image has been processed
		Username: "johndoe",
	})
}

func downloadTempFile(f multipart.File, h *multipart.FileHeader) (*os.File, error) {
	log.Printf("Received file '%s', size is %d bytes", h.Filename, h.Size)

	defer f.Close()

	tmpFile, err := ioutil.TempFile("", "lolcatz-upload-")
	if err != nil {
		return nil, err
	}

	defer tmpFile.Close()

	numBytes, err := io.Copy(tmpFile, f)
	if err != nil {
		return nil, err
	}

	log.Printf("%d bytes written to %s", numBytes, tmpFile.Name())
	return tmpFile, nil
}
