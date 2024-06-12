package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
)

func main() {
	handler := chi.NewRouter()

	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)

	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pages.Index().Render(r.Context(), w)
	})

	handler.Get("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		pages.Login().Render(r.Context(), w)
	})

	handler.Get("/auth/signup", func(w http.ResponseWriter, r *http.Request) {
		pages.Signup(&pages.SignupProps{}).Render(r.Context(), w)
	})
	
	handler.Post("/auth/signup", func(w http.ResponseWriter, r *http.Request) {
		fname := r.FormValue("first_name")
		lname := r.FormValue("last_name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		props := pages.SignupProps{}

		if (fname == "") {
			props.Errors.FirstName = "This field is required"
		}

		if (email == "") {
			props.Errors.Email = "This field is required"
		}

		if (password == "") {
			props.Errors.Password = "This field is required"
		} else if (len(password) < 8 || len(password) > 64) {
			props.Errors.Password = "The password must be between 8-64 characters"
		}

		props.Values.FirstName = fname
		props.Values.LastName = lname
		props.Values.Email = email
		props.Values.Password = password

		if r.Header.Get("HX-Request") == "" {
			pages.Signup(&props).Render(r.Context(), w)
		} else {
			pages.SignupForm(&props).Render(r.Context(), w)
		}
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
