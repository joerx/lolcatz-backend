package middleware

import (
	"log"
	"net/http"
)

type corsFilter struct {
	AllowOrigin string
}

// HeaderAccessControlAllowOrigin represents the Access-Control-Allow-Origin http header
const HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"

// HeaderAccessControlAllowMethods represents the Access-Control-Allow-Methods http header
const HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"

// HeaderAccessControlAllowHeaders represents the Access-Control-Allow-Headers http header
const HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"

// CorsWithOrigin returns a middleware writing CORS headers for the given origin, all methods and most common headers
func CorsWithOrigin(origin string) Func {
	cf := &corsFilter{AllowOrigin: origin}
	return cf.doFilter
}

func (cf *corsFilter) doFilter(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cf.setupCors(w, r)
		// Stop handling OPTIONS requests here
		if r.Method == http.MethodOptions {
			return
		}
		f(w, r)
	}
}

func (cf *corsFilter) setupCors(w http.ResponseWriter, r *http.Request) {
	log.Printf("Setting up cors config")
	w.Header().Set(HeaderAccessControlAllowOrigin, cf.AllowOrigin)
	w.Header().Set(HeaderAccessControlAllowMethods, "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set(HeaderAccessControlAllowHeaders, "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
