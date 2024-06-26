package main

import (
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sakithb/hcblk-server/internal/db"
	"github.com/sakithb/hcblk-server/internal/middleware"
	"github.com/sakithb/hcblk-server/internal/routes"
	"github.com/sakithb/hcblk-server/internal/server"
	"github.com/sakithb/hcblk-server/internal/services"
)

func main() {
	db := db.New()

	sm := scs.New()
	sm.Lifetime = time.Hour * 24 * 30
	sm.IdleTimeout = time.Hour * 24 * 7
	sm.Store = sqlite3store.New(db.DB)

	h := chi.NewRouter()

	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)
	h.Use(middleware.StripSlashes)
	h.Use(sm.LoadAndSave)
	h.Use(mw.Authentication(sm))

	authService := &services.AuthService{DB: db}
	userService := &services.UserService{DB: db}

	h.Mount("/", routes.NewIndexHandler().Router())
	h.Mount("/assets", routes.NewAssetsHandler().Router())
	h.Mount("/auth", routes.NewAuthHandler(authService, userService, sm).Router())
	h.Mount("/user", routes.NewUserHandler(userService).Router())

	server.StartServer(h)
}
