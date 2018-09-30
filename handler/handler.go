package handler

import (
	"net/http"
)

type handler struct{}

// New retorna um http.Handler
func New() http.Handler {
	h := handler{}

	mux := http.NewServeMux()

	mux.HandleFunc("/", h.delay)

	return mux
}
