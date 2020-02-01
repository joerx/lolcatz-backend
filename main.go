package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joerx/lolcatz-backend/pkg/middleware"
	"github.com/joerx/lolcatz-backend/pkg/s3"
)

const bucket = "sandbox-lolcatz-be-storage-468871832330"
const region = "ap-southeast-1"

// MsgResponse is a simple message response
type MsgResponse struct {
	Message string `json:"message"`
}

func main() {
	log.Println("Starting server")

	logging := middleware.Logging
	cors := middleware.Cors("http://localhost:3000")
	warez := []middleware.Func{logging, cors}

	addHandler("/", warez, handleRoot)
	addHandler("/upload", warez, handleUpload)

	http.ListenAndServe(":8000", nil)
}

func addHandler(path string, warez []middleware.Func, f http.HandlerFunc) {
	http.HandleFunc(path, middleware.Chain(warez, f))
}

func writeResponseMsg(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	res, _ := json.Marshal(MsgResponse{Message: message})
	fmt.Fprint(w, string(res))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	writeResponseMsg(w, 200, "ok")
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s /upload", r.Method)

	if r.Method != "POST" {
		writeResponseMsg(w, 405, "Method not supported")
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeResponseMsg(w, 400, "Can't parse this as multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
	}

	log.Printf("Received file %s with header %v", header.Filename, header.Header)

	tmpFile, err := downloadTempFile(file, header)
	if err != nil {
		writeResponseMsg(w, 500, "Error processing uploaded file")
		return
	}

	in := s3.UploadRequest{
		Filename:     tmpFile.Name(),
		OriginalName: header.Filename,
		ContentType:  header.Header.Get("Content-Type"),
	}

	cfg := s3.Config{
		Bucket: bucket,
		Region: region,
	}

	if err := s3.Upload(in, cfg); err != nil {
		writeResponseMsg(w, 500, "Error processing uploaded file")
		return
	}

	writeResponseMsg(w, 200, "All good")
}

func downloadTempFile(f multipart.File, h *multipart.FileHeader) (*os.File, error) {
	log.Printf("Received file %s, size is %d", h.Filename, h.Size)

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
