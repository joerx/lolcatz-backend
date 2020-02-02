package middleware

import (
	"net/http"
)

// Func is a middleware function that allows pre-/post-processing of requests
// It returns a wrapped handler function allowing to intercept the HTTP call made
// to the original handler.
type Func func(http.HandlerFunc) http.HandlerFunc

// Chain allows to easily chain a sclice of middlewares around a handler
func Chain(ms []Func, f http.HandlerFunc) http.HandlerFunc {
	handler := f
	for _, m := range ms {
		handler = m(handler)
	}
	return handler
}
