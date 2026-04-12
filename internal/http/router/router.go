package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/http/routes"
)

var allowedMethodOrder = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

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

	mux.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range allowedMethodsForPath(mux, requestPath(r)) {
			w.Header().Add("Allow", method)
		}
		response.Error(w, errors.ErrMethodNotAllowed)
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

func requestPath(req *http.Request) string {
	if req.URL.RawPath != "" {
		return req.URL.RawPath
	}
	if req.URL.Path != "" {
		return req.URL.Path
	}
	return "/"
}

func allowedMethodsForPath(router chi.Router, path string) []string {
	allowedMethods := make([]string, 0, len(allowedMethodOrder))
	for _, method := range allowedMethodOrder {
		if router.Match(chi.NewRouteContext(), method, path) {
			allowedMethods = append(allowedMethods, method)
		}
	}

	return allowedMethods
}
