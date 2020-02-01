package middleware

import (
	"log"
	"net/http"
)

// Func is a middleware function that allows pre-/post-processing of requests
// It returns a wrapped handler function allowing to intercept the HTTP call made
// to the original handler.
type Func func(http.HandlerFunc) http.HandlerFunc

// Logging is a simple request logging middleware
func Logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		f(w, r)
	}
}

// Cors is a middleware to write generic cors headers
func Cors(origin string) Func {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			// stop options requests here
			if r.Method == "OPTIONS" {
				return
			}
			f(w, r)
		}
	}

}

// Chain allows to easily chain a sclice of middlewares around a handler
func Chain(ms []Func, f http.HandlerFunc) http.HandlerFunc {
	handler := f
	for _, m := range ms {
		handler = m(handler)
	}
	return handler
}
