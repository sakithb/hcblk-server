package routes

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/sakithb/hcblk-server/internal/models"
	"github.com/sakithb/hcblk-server/internal/services"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
	"github.com/sakithb/hcblk-server/internal/utils"
)

type AuthHandler struct {
	Sessions    *scs.SessionManager
	AuthService *services.AuthService
	UserService *services.UserService
}

func NewAuthHandler(as *services.AuthService, us *services.UserService, sm *scs.SessionManager) *AuthHandler {
	return &AuthHandler{
		Sessions:    sm,
		AuthService: as,
		UserService: us,
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
	r.Get("/logout", h.Logout)
	r.Post("/login", h.Login)
	r.Post("/signup", h.Signup)

	return r
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("t")
	if t == "" {
		http.Error(w, "Invalid token", 400)
		return
	}

	u, err := h.AuthService.VerifyToken(t)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid token", 400)
		} else {
			utils.HandleServerError(w, err)
		}

		return
	}

	exists, err := h.UserService.UserExists(u.Email)
	if exists {
		http.Error(w, "User with email already exists", 400)
		return
	}

	err = h.UserService.CreateUser(u.FirstName, u.LastName, u.Email, u.Hash)
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	err = h.AuthService.DeleteToken(t)
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	rememberMe := r.FormValue("remember_me")

	props := pages.LoginProps{
		Email:      email,
		Password:   password,
		RememberMe: rememberMe != "",
	}

	valid, err := h.AuthService.VerifyPassword(password, email)
	if err != nil {
		log.Println(err)
		props.ServerError = true
		pages.LoginForm(&props).Render(r.Context(), w)
		return
	}

	if !valid {
		props.Invalid = true
		pages.LoginForm(&props).Render(r.Context(), w)
		return
	}

	u, err := h.UserService.GetUserByEmail(email)
	if err != nil {
		log.Println(err)
		props.ServerError = true
		pages.LoginForm(&props).Render(r.Context(), w)
	} else {
		h.Sessions.Put(r.Context(), "user", u)
		h.Sessions.RememberMe(r.Context(), props.RememberMe)

		w.Header().Add("HX-Redirect", "/")
	}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	fname := r.FormValue("first_name")
	lname := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	props := pages.SignupProps{}
	props.Values.FirstName = fname
	props.Values.LastName = lname
	props.Values.Email = email
	props.Values.Password = password

	if fname == "" {
		props.Errors.FirstName = "This field is required"
	}

	if email == "" {
		props.Errors.Email = "This field is required"
	}

	if password == "" {
		props.Errors.Password = "This field is required"
	} else if len(password) < 8 || len(password) > 64 {
		props.Errors.Password = "The password must be between 8-64 characters"
	}

	if props.Errors.Email == "" && props.Errors.FirstName == "" && props.Errors.LastName == "" && props.Errors.Password == "" {
		hash := h.AuthService.GenerateHash(password)
		ou := &models.OnboardingUser{
			FirstName: fname,
			LastName:  lname,
			Email:     email,
			Hash:      hash,
		}

		t, err := h.AuthService.GenerateToken(ou)
		if err != nil {
			log.Println(err)
			props.ServerError = true
		} else {
			log.Println("http://localhost:3000/auth/verify?t=" + url.QueryEscape(t))
			props.Emailed = true
		}
	}

	pages.SignupForm(&props).Render(r.Context(), w)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Sessions.Clear(r.Context())
	http.Redirect(w, r, "/", http.StatusFound)
}
