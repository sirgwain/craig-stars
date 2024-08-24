//go:build !wasi && !wasm

package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
	"github.com/sirgwain/craig-stars/cs"
)

func (s *server) techs(w http.ResponseWriter, r *http.Request) {
	rest.RenderJSON(w, cs.StaticTechStore)
}

func (s *server) tech(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		render.Render(w, r, ErrNotFound)
		return
	}

	tech := cs.StaticTechStore.GetTech(name)
	if tech == nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	rest.RenderJSON(w, tech)

}
