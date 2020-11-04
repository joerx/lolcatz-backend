package middleware

import (
	"log"
	"net/http"
)

type corsConfig struct {
	AllowOrigin string
}

// CorsWithOrigin returns a middleware writing CORS headers for the given origin, all methods and most common headers
func CorsWithOrigin(origin string) Func {
	cf := corsConfig{AllowOrigin: origin}
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			setupCors(w, cf)
			// Stop options requests here
			if r.Method == "OPTIONS" {
				return
			}
			f(w, r)
		}
	}
}

func setupCors(w http.ResponseWriter, cf corsConfig) {
	log.Printf("Setting up cors config")
	w.Header().Set("Access-Control-Allow-Origin", cf.AllowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
