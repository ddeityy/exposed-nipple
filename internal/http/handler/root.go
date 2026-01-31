package handler

import (
	"net/http"
	"nipple/internal/manager"
	"text/template"

	"github.com/charmbracelet/log"
)

type rootHandler struct {
	connManager manager.ConnectManager
	lg          *log.Logger
}

func NewRootHandler(connManager manager.ConnectManager, lg *log.Logger) *rootHandler {
	return &rootHandler{
		connManager,
		lg,
	}
}

func (h *rootHandler) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.connManager.GetServerStatus()
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
