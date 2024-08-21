package main

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
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

	dev, _ := os.LookupEnv("DEV_MODE")
	isDev := dev == "1"

	var loglevel slog.Level
	var out io.Writer

	if isDev {
		loglevel = slog.LevelDebug
		out = os.Stdout
	} else {
		err := os.MkdirAll("/var/log/hcblk/", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		fname := strconv.Itoa(int(time.Now().Unix()))

		access, err := os.Create(path.Join("/var/log/hcblk/", fname+"-access.log"))
		if err != nil {
			log.Fatal(err)
		}

		std, err := os.Create(path.Join("/var/log/hcblk/", fname+".log"))
		if err != nil {
			log.Fatal(err)
		}

		loglevel = slog.LevelInfo
		out = access
		log.SetOutput(std)
	}

	logger := httplog.NewLogger("server", httplog.Options{
		LogLevel:        loglevel,
		Concise:         !isDev,
		RequestHeaders:  true,
		ResponseHeaders: true,
		JSON:            !isDev,
		Writer:          out,
	})

	h := chi.NewRouter()

	h.Use(middleware.RealIP)
	h.Use(middleware.Heartbeat("/ping"))
	h.Use(httplog.RequestLogger(logger))
	h.Use(middleware.Recoverer)
	h.Use(middleware.StripSlashes)
	h.Use(sm.LoadAndSave)
	h.Use(mw.Authentication(sm))

	authService := &services.AuthService{DB: db}
	userService := &services.UserService{DB: db}
	listingService := &services.ListingService{DB: db}
	uiService := &services.UIService{DB: db}

	h.Mount("/", routes.NewIndexHandler(listingService, uiService).Router())
	h.Mount("/assets", routes.NewAssetsHandler().Router())
	h.Mount("/auth", routes.NewAuthHandler(authService, userService, sm).Router())
	h.Mount("/me", routes.NewMeHandler(userService, uiService, listingService, authService).Router())
	h.Mount("/user", routes.NewUserHandler(userService).Router())
	h.Mount("/listing", routes.NewListingHandler(listingService).Router())

	log.Println("Starting server...")
	server.StartServer(h)
}
