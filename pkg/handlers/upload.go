package handlers

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joerx/lolcatz-backend/pkg/s3"
)

type uploadConf struct {
	region string
	bucket string
}

// Upload returns a handler that handles file uploads to the given region and bucket
func Upload(region string, bucket string) http.HandlerFunc {
	cf := uploadConf{region, bucket}
	return func(w http.ResponseWriter, r *http.Request) {
		upload(w, r, cf)
	}
}

func upload(w http.ResponseWriter, r *http.Request, u uploadConf) {
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
		errorHandler(err, w)
		return
	}

	in := s3.UploadRequest{
		Filename:     tmpFile.Name(),
		OriginalName: header.Filename,
		ContentType:  header.Header.Get("Content-Type"),
	}

	cfg := s3.Config{
		Bucket: u.bucket,
		Region: u.region,
	}

	if err := s3.Upload(in, cfg); err != nil {
		errorHandler(err, w)
		return
	}

	writeResponseMsg(w, http.StatusOK, "All good")
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
