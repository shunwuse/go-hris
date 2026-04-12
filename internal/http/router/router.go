package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/http/routes"
)

// Router composes the HTTP middleware and route tree.
type Router struct {
	mux chi.Router
}

func New(
	middlewares middlewares.CommonMiddlewares,
	routes routes.Routes,
) *Router {
	mux := chi.NewRouter()

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, errors.ErrNotFound)
	})

	// Setup common middlewares.
	middlewares.Setup(mux)

	// Setup routes.
	routes.Setup(mux)

	return &Router{
		mux: mux,
	}
}

// ServeHTTP delegates requests to the composed chi router.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
