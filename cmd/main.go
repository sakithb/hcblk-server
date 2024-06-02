package main

import (
	"log"
	"net/http"

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

	server := http.Server{Addr: "0.0.0.0:3000", Handler: handler }

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
