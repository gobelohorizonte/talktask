package handler

import (
	"net/http"

	"github.com/urfave/negroni"
	"github.com/waltton/talktask/acd"
)

type handler struct {
	jobs chan acd.Job
}

// New retorna um http.Handler
func New(jobs chan acd.Job) http.Handler {
	h := handler{}

	mux := http.NewServeMux()

	mux.HandleFunc("/", h.delay)

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(mux)

	return n
}
