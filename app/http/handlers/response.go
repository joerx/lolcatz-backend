package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joerx/lolcatz-backend/http/errors"
)

// MsgResponse is a simple generic message response
type MsgResponse struct {
	Message string `json:"message"`
}

// ErrorResponse is a simple response with a single error message
type ErrorResponse struct {
	Error string `json:"error"`
}

func writeResponseMsg(w http.ResponseWriter, statusCode int, message string) {
	writeResponseJSON(w, statusCode, MsgResponse{Message: message})
}

func writeResponseJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	res, _ := json.Marshal(data)
	fmt.Fprint(w, string(res))
}

// func internalServerError(w http.ResponseWriter, err error) {
// 	log.Printf("Error %v", err)
// 	errorHandler(w, http.StatusInternalServerError, err.Error())
// }

// func badRequest(w http.ResponseWriter, msg string) {
// 	errorHandler(w, http.StatusBadRequest, msg)
// }

// func methodNotAllowed(w http.ResponseWriter, msg string) {
// 	errorHandler(w, http.StatusMethodNotAllowed, msg)
// }

func errorHandler(w http.ResponseWriter, err error) {
	switch t := err.(type) {
	case errors.Error:
		writeResponseMsg(w, t.StatusCode, t.Message)
	default:
		writeResponseMsg(w, http.StatusInternalServerError, err.Error())
	}
}
