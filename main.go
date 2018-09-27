package main

import (
	"log"
	"net/http"

	"github.com/waltton/talktask/handler"
)

func main() {
	h := handler.New()

	log.Panic(http.ListenAndServe("0.0.0.0:8080", h))
}
