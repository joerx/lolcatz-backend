package routing

import (
	"net/http"

	"github.com/joerx/lolcatz-backend/pkg/middleware"
)

// Router is simple helper to set up http handlers and filters
type Router struct {
	filters []middleware.Func
}

// NewRouter creates a new request router
func NewRouter() *Router {
	return &Router{}
}

// Handle registers a new request handler
func (d *Router) Handle(path string, f http.HandlerFunc) {
	http.HandleFunc(path, middleware.Chain(d.filters, f))
}

// Filter adds a new middleware to the filter chain
func (d *Router) Filter(f middleware.Func) {
	d.filters = append(d.filters, f)
}
