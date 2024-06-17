package main

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sakithb/hcblk-server/internal/db"
	"github.com/sakithb/hcblk-server/internal/routes"
	"github.com/sakithb/hcblk-server/internal/server"
)

func main() {
	db := db.New()

	sm := scs.New()
	sm.Lifetime = time.Hour * 24 * 30
	sm.IdleTimeout = time.Hour * 24 * 7

	h := chi.NewRouter()

	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)
	h.Use(sm.LoadAndSave)

	h.Mount("/", routes.NewIndexHandler().Route())
	h.Mount("/assets", routes.NewAssetsHandler().Route())
	h.Mount("/auth", routes.NewAuthHandler(db, sm).Router())

	server.StartServer(h)
}
