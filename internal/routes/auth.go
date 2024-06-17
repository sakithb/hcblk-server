package routes

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sakithb/hcblk-server/internal/models"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
	"github.com/sakithb/hcblk-server/internal/utils"
)

type AuthHandler struct {
	Tokens   *sync.Map
	Sessions *scs.SessionManager
	DB       *sqlx.DB
}

func NewAuthHandler(db *sqlx.DB, sm *scs.SessionManager) *AuthHandler {
	return &AuthHandler{
		Tokens:   &sync.Map{},
		Sessions: sm,
		DB:       db,
	}
}

func (h *AuthHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		pages.Login(&pages.LoginProps{}).Render(r.Context(), w)
	})

	r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		pages.Signup(&pages.SignupProps{}).Render(r.Context(), w)
	})

	r.Get("/verify", h.Verify)
	r.Post("/login", h.Login)
	r.Post("/signup", h.Signup)

	return r
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("t")
	if t == "" {
		http.Error(w, "Bad request", 400)
		return
	}

	uv, ok := h.Tokens.Load(t)

	if !ok {
		http.Error(w, "Bad request", 400)
		return
	}

	u := uv.(models.UnverifiedUser)

	id, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", 500)
		return
	}

	_, err = h.DB.Exec("INSERT INTO users VALUES(?,?,?,?,?,?)", id.String(), u.FirstName, u.LastName, u.Email, u.Password, time.Now().Unix())
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", 500)
		return
	}

	http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	rememberMe := r.FormValue("remember_me")

	hash := utils.GenerateHashFromPassword([]byte(password))

	u := models.User{}
	err := h.DB.Get(&u, "SELECT id,first_name,last_name,email FROM users WHERE email = ? AND password = ?", email, hash)

	if err == sql.ErrNoRows {
		props := pages.LoginProps{}

		props.Error = true
		props.Email = email
		props.Password = password
		props.RememberMe = rememberMe != ""

		pages.LoginForm(&props).Render(r.Context(), w)
	} else if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	} else {
		h.Sessions.Put(r.Context(), "user", u)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	fname := r.FormValue("first_name")
	lname := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	props := pages.SignupProps{Success: true}

	if fname == "" {
		props.Errors.FirstName = "This field is required"
		props.Success = false
	}

	if email == "" {
		props.Errors.Email = "This field is required"
		props.Success = false
	}

	if password == "" {
		props.Errors.Password = "This field is required"
		props.Success = false
	} else if len(password) < 8 || len(password) > 64 {
		props.Errors.Password = "The password must be between 8-64 characters"
		props.Success = false
	}

	hash := utils.GenerateHashFromPassword([]byte(password))

	props.Values.FirstName = fname
	props.Values.LastName = lname
	props.Values.Email = email
	props.Values.Password = hash

	if props.Success {
		b := utils.GenerateRandomBytes(32)
		t := base64.StdEncoding.EncodeToString(b)

		h.Tokens.Store(t, models.UnverifiedUser{
			FirstName: fname,
			LastName:  lname,
			Email:     email,
			Password:  hash,
		})

		// TODO: Email token
		log.Println("http://localhost:3000/auth/verify?t=" + url.QueryEscape(t))
	}

	pages.SignupForm(&props).Render(r.Context(), w)
}
