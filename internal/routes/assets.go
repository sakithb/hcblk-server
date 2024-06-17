package routes

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type AssetsHandler struct {
}

func NewAssetsHandler() *AssetsHandler {
	return &AssetsHandler{}
}

func (h *AssetsHandler) Route() chi.Router {
	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, dev := os.LookupEnv("DEV_MODE"); dev {
			w.Header().Set("Cache-Control", "no-store")
		}

		http.StripPrefix("/assets/",
			http.FileServer(http.Dir("./assets/dist/"))).ServeHTTP(w, r)
	})

	return r
}
