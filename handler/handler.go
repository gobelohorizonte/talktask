package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type handler struct {
	r *rand.Rand
}

// New retorna um http.Handler
func New() http.Handler {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	h := handler{r}

	mux := http.NewServeMux()

	mux.HandleFunc("/", h.delay)

	return mux
}

func (h *handler) delay(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second)
	fmt.Fprint(w, "ok")
}
