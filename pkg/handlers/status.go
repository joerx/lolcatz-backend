package handlers

import "net/http"

// Status is the base route handler responding with a simple status message
func Status(w http.ResponseWriter, r *http.Request) {
	writeResponseMsg(w, http.StatusOK, "ok")
}
