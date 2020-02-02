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
	writeReponse(w, statusCode, MsgResponse{Message: message})
}

func writeReponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	res, _ := json.Marshal(data)
	fmt.Fprint(w, string(res))
}

func errorHandler(w http.ResponseWriter, err error) {
	log.Printf("Error %v", err)
	writeResponseMsg(w, http.StatusInternalServerError, "Error processing uploaded file")
	return
}
