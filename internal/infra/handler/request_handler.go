package handler

import (
	"github.com/go-chi/chi/v5"
)

// RequestHandler structure.
type RequestHandler struct {
	Router chi.Router
}

func New() *RequestHandler {
	router := chi.NewRouter()

	return &RequestHandler{
		Router: router,
	}
}
