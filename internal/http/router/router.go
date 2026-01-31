package router

import (
	"net/http"

	"nipple/internal/http/middleware"
	"nipple/internal/provider"

	"github.com/charmbracelet/log"
)

type router struct {
	root *http.ServeMux
	mws  []func(http.Handler) http.Handler
	lg   *log.Logger
}

func New(prov provider.Provider) *router {
	root := http.NewServeMux()

	RegisterRootHandler(prov, root)

	r := router{
		root,
		nil,
		prov.Logger(),
	}

	r.initMiddlewares()

	return &r
}

func (r *router) Use(mws ...func(http.Handler) http.Handler) {
	r.mws = append(r.mws, mws...)
}

func (r *router) Handler() http.Handler {
	res := http.Handler(r.root)
	mwLen := len(r.mws)
	for i := range r.mws {
		res = r.mws[mwLen-i-1](res)
	}
	return res
}

func (r *router) initMiddlewares() {
	mw := middleware.New(r.lg)
	r.Use(mw.RecoverMiddleware)
	r.Use(mw.RequestLogger)
}
