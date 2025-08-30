package handler

import (
	"app"
	"fmt"
	"net/http"
)

type Handler struct {
	App *app.App
}

func NewHandler(a *app.App) Handler {
	return Handler{
		App: a,
	}
}

func (handler *Handler) Healthz(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(rw, "server is ok")
}
