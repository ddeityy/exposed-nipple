package middleware

import (
	"net/http"

	"github.com/charmbracelet/log"
)

type middleware struct {
	lg *log.Logger
}

func New(lg *log.Logger) *middleware {
	return &middleware{
		lg,
	}
}

type writer struct {
	http.ResponseWriter
	resStatus int
	resSize   int
}

func newWriter(w http.ResponseWriter) *writer {
	return &writer{
		ResponseWriter: w,
	}
}

func (w *writer) WriteHeader(status int) {
	if w.resStatus == 0 {
		w.resStatus = status
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *writer) Write(body []byte) (int, error) {
	if w.resStatus == 0 {
		w.WriteHeader(http.StatusOK)
	}

	var err error
	w.resSize, err = w.ResponseWriter.Write(body)

	return w.resSize, err
}
