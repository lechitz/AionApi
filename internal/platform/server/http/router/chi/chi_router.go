// Package chi implementa a porta Router sobre chi v5.
package chi

import (
	"net/http"

	chiv5 "github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// rtr implementa a porta Router sobre chi v5.
type rtr struct {
	chi         chiv5.Router
	errorHandle func(http.ResponseWriter, *http.Request, error)
}

// New cria um Router baseado em chi v5, sem middlewares por padrão.
// Middlewares de plataforma (request-id, recovery, cors...) devem ser aplicados
// no composer da plataforma antes de registrar os módulos de domínio.
func New() ports.Router {
	return &rtr{chi: chiv5.NewRouter()}
}

// Use aplica middlewares no pipeline do chi.
// Aceita N middlewares; ignora nils defensivamente.
func (r *rtr) Use(mw ...ports.Middleware) {
	for _, m := range mw {
		if m != nil {
			r.chi.Use(m)
		}
	}
}

// Group cria uma subárvore de rotas com 'prefix'.
// A função fn recebe um Router isolado daquela subárvore.
func (r *rtr) Group(prefix string, fn func(ports.Router)) {
	r.chi.Route(prefix, func(cr chiv5.Router) {
		fn(&rtr{chi: cr})
	})
}

// GroupWith cria uma subárvore aplicando um middleware específico apenas nela.
func (r *rtr) GroupWith(m ports.Middleware, fn func(ports.Router)) {
	r.chi.Group(func(gr chiv5.Router) {
		if m != nil {
			gr.Use(m)
		}
		fn(&rtr{chi: gr})
	})
}

// Mount pendura um http.Handler pronto em um prefixo (ex.: /graphql).
func (r *rtr) Mount(prefix string, h http.Handler) {
	r.chi.Mount(prefix, h)
}

// Handle registra uma rota para um método arbitrário.
func (r *rtr) Handle(method, path string, h http.Handler) {
	r.chi.Method(method, path, h)
}

// Atalhos para verbos comuns.
func (r *rtr) GET(path string, h http.Handler)    { r.Handle(http.MethodGet, path, h) }
func (r *rtr) POST(path string, h http.Handler)   { r.Handle(http.MethodPost, path, h) }
func (r *rtr) PUT(path string, h http.Handler)    { r.Handle(http.MethodPut, path, h) }
func (r *rtr) DELETE(path string, h http.Handler) { r.Handle(http.MethodDelete, path, h) }

// SetNotFound define o handler 404 customizado.
func (r *rtr) SetNotFound(h http.Handler) {
	r.chi.NotFound(h.ServeHTTP)
}

// SetMethodNotAllowed define o handler 405 customizado.
func (r *rtr) SetMethodNotAllowed(h http.Handler) {
	r.chi.MethodNotAllowed(h.ServeHTTP)
}

// SetError armazena um handler de erro centralizado (opcional).
// OBS: o chi não invoca isso automaticamente; a plataforma pode utilizá-lo
// dentro de um middleware de erro/recovery para padronizar respostas.
func (r *rtr) SetError(h func(http.ResponseWriter, *http.Request, error)) {
	r.errorHandle = h
}

// ServeHTTP cumpre a interface http.Handler e delega para o chi.
func (r *rtr) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}
