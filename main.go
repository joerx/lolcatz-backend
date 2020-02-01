package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/joerx/lolcatz-backend/pkg/middleware"
)

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

	defer file.Close()

	log.Printf("Received file %s, size is %d", header.Filename, header.Size)

	tmpFile, err := ioutil.TempFile("", "lolcatz-upload-")
	if err != nil {
		writeResponseMsg(w, 500, "Error processing uploaded file")
		return
	}

	defer tmpFile.Close()
	numBytes, err := io.Copy(tmpFile, file)
	if err != nil {
		writeResponseMsg(w, 500, "Error processing uploaded file")
		return
	}

	log.Printf("%d bytes written to %s", numBytes, tmpFile.Name())

	writeResponseMsg(w, 200, "All good")
}
