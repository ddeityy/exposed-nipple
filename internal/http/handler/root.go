package handler

import (
	"embed"
	"net/http"
	"nipple/internal/manager"
	"text/template"

	"github.com/charmbracelet/log"

	_ "embed"
)

//go:embed templates
var templates embed.FS

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

	tmpl := template.Must(template.ParseFS(templates, "templates/*.html"))
	err = tmpl.ExecuteTemplate(w, "index.html", status)
	if err != nil {
		h.lg.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
