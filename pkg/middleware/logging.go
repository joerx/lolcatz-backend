package middleware

import (
	"log"
	"net/http"
)

// Logging is a simple request logging middleware
func Logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		f(w, r)
	}
}
