package handler

import (
	"fmt"
	"net/http"

	"github.com/waltton/talktask/acd"
)

func (h *handler) delay(w http.ResponseWriter, r *http.Request) {
	h.jobs <- acd.Job{}
	fmt.Fprint(w, "ok")
}
