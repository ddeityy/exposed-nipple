package handler

import (
	"net/http"
	"nipple/internal/rcon"
	"text/template"

	"github.com/charmbracelet/log"
)

type rootHandler struct {
	rconClient *rcon.Client
	lg         *log.Logger
}

func NewRootHandler(rconClient *rcon.Client, lg *log.Logger) *rootHandler {
	return &rootHandler{
		rconClient,
		lg,
	}
}

func (h *rootHandler) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.rconClient.GetServerStatus()
	if err != nil {
		h.lg.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("internal/http/templates/index.html"))
	err = tmpl.Execute(w, status)
	if err != nil {
		h.lg.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
