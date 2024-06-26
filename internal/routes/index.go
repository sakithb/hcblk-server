package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
)

type IndexHandler struct {
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Index().Render(r.Context(), w)
	})

	return r
}
