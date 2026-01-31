package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"nipple/internal/http/router"
	"nipple/internal/provider"
)

type Server struct {
	srv           *http.Server
	cancelBaseCtx context.CancelFunc
}

func New(prov provider.Provider) *Server {
	root := router.New(prov)

	cfg := prov.Config().HTTP

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      root.Handler(),
		ReadTimeout:  cfg.ReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{
		srv: srv,
	}
}

func (h *Server) Info() string {
	return h.srv.Addr
}

func (h *Server) Run(ctx context.Context) error {
	baseCtx, cancel := context.WithCancel(ctx)
	h.cancelBaseCtx = cancel

	h.srv.BaseContext = func(_ net.Listener) context.Context {
		return baseCtx
	}

	return h.srv.ListenAndServe()
}

func (h *Server) Stop(ctx context.Context) error {
	h.srv.SetKeepAlivesEnabled(false)
	h.cancelBaseCtx()

	err := h.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("h.srv.Shutdown: %w", err)
	}

	return nil
}
