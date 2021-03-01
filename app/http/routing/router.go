package routing

import (
	"net/http"

	"github.com/joerx/lolcatz-backend/http/middleware"
)

// Router is essentially a small wrapper around http.ServeMux to set up http handlers and filters
type Router struct {
	filters []middleware.Func
	mux     *http.ServeMux
}

// NewRouter creates a new request router
func NewRouter() *Router {
	m := http.NewServeMux()
	return &Router{mux: m}
}

// Handle registers a new request handler
func (d *Router) Handle(path string, f http.HandlerFunc) {
	d.mux.HandleFunc(path, middleware.Chain(d.filters, f))
}

// Filter adds a new middleware to the filter chain
func (d *Router) Filter(f middleware.Func) {
	d.filters = append(d.filters, f)
}

// ServeHTTP implements the http.Handler interface for Router
func (d *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.mux.ServeHTTP(w, r)
}
