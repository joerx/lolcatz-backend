package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// MsgResponse is a simple message response
type MsgResponse struct {
	Message string `json:"message"`
}

func writeResponseMsg(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	res, _ := json.Marshal(MsgResponse{Message: message})
	fmt.Fprint(w, string(res))
}

func errorHandler(err error, w http.ResponseWriter) {
	log.Printf("Error %v", err)
	writeResponseMsg(w, http.StatusInternalServerError, "Error processing uploaded file")
	return
}
