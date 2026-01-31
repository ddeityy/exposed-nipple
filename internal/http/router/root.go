package router

import (
	"net/http"
	handler "nipple/internal/http/handler"
	"nipple/internal/provider"
)

func RegisterRootHandler(provider provider.Provider, root *http.ServeMux) {
	handler := handler.NewRootHandler(
		provider.RconClient(),
		provider.Logger(),
	)
	root.HandleFunc("GET /", handler.GetServerStatus)
}
