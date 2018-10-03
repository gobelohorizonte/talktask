package handler

import (
	"net/http"

	"github.com/urfave/negroni"
)

type handler struct{}

// New retorna um http.Handler
func New() http.Handler {
	h := handler{}

	mux := http.NewServeMux()

	mux.HandleFunc("/", h.delay)

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(mux)

	return n
}
