package routes

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AssetsHandler struct {
}

func NewAssetsHandler() *AssetsHandler {
	return &AssetsHandler{}
}

func (h *AssetsHandler) Router() chi.Router {
	r := chi.NewRouter()

	if _, dev := os.LookupEnv("DEV_MODE"); dev {
		r.Use(middleware.NoCache)
	}

	r.Use(middleware.Compress(5, "image/*", "text/*"))

	r.Get("/users/{userID}", h.Users)
	r.Get("/*", h.Generic)

	return r
}

func (h *AssetsHandler) Users(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	if f, err := os.Open("./assets/dist/users/" + userID); os.IsNotExist(err) {
		f, err = os.Open("./assets/dist/user-circle.webp")
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		_, err := f.WriteTo(w)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	} else {
		_, err := f.WriteTo(w)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}

func (h *AssetsHandler) Generic(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/assets/",
		http.FileServer(http.Dir("./assets/dist/"))).ServeHTTP(w, r)
}
