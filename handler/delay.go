package handler

import (
	"fmt"
	"net/http"
	"time"
)

func (h *handler) delay(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second)
	fmt.Fprint(w, "ok")
}
