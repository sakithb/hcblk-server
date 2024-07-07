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

	r.Get("/users/{id}", h.GetUserPic)
	r.Get("/*", h.GetAll)

	return r
}

func (h *AssetsHandler) GetUserPic(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if f, err := os.Open("assets/dist/users/" + id + ".png"); os.IsNotExist(err) {
		f, err = os.Open("assets/dist/user-circle.webp")
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

func (h *AssetsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/assets/",
		http.FileServer(http.Dir("assets/dist/"))).ServeHTTP(w, r)
}
