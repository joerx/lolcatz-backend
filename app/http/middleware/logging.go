package middleware

import (
	"log"
	"net/http"
)

// Logging is a simple request logging middleware
func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		next(w, r)
	}
}
