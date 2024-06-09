package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sakithb/hcblk-server/views"
)

func main() {
	handler := chi.NewRouter()

	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)

	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})

	handler.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		if _, dev := os.LookupEnv("DEV_MODE"); dev {
			w.Header().Set("Cache-Control", "no-store")
		}

		http.StripPrefix("/assets/",
			http.FileServer(http.Dir("./assets/dist/"))).ServeHTTP(w, r)
	})

	server := http.Server{Addr: ":3000", Handler: handler}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
