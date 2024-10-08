package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/sakithb/hcblk-server/internal/services"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
	"github.com/sakithb/hcblk-server/internal/utils"
)

type AuthHandler struct {
	Sessions    *scs.SessionManager
	AuthService *services.AuthService
	UserService *services.UserService
	SMTPAuth    smtp.Auth
}

func NewAuthHandler(as *services.AuthService, us *services.UserService, sm *scs.SessionManager) *AuthHandler {
	return &AuthHandler{
		Sessions:    sm,
		AuthService: as,
		UserService: us,
		SMTPAuth:    smtp.PlainAuth("", "postmaster@hcblk.live", "dda83504bd78534fe107e08f59867b63-2b91eb47-4e042e39", "smtp.eu.mailgun.org"),
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

	r.Get("/forgot_password", func(w http.ResponseWriter, r *http.Request) {
		pages.ForgotPassword(&pages.ForgotPasswordProps{}).Render(r.Context(), w)
	})

	r.Get("/reset", func(w http.ResponseWriter, r *http.Request) {
		pages.Reset(&pages.ResetProps{}).Render(r.Context(), w)
	})

	r.Get("/verify", h.GetVerify)
	r.Get("/logout", h.GetLogout)
	r.Post("/login", h.PostLogin)
	r.Post("/signup", h.PostSignup)
	r.Post("/forgot_password", h.PostForgotPassword)
	r.Post("/reset", h.PostReset)

	return r
}

func (h *AuthHandler) GetVerify(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("t")
	if t == "" {
		http.Error(w, "Invalid token", 400)
		return
	}

	u, err := h.AuthService.VerifySignupToken(t)
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

	err = h.UserService.CreateUser(u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	err = h.AuthService.DeleteSignupToken(t)
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	http.Redirect(w, r, "/auth/login", http.StatusFound)
}

func (h *AuthHandler) GetLogout(w http.ResponseWriter, r *http.Request) {
	h.Sessions.Clear(r.Context())
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	rememberMe := r.FormValue("remember_me")

	props := pages.LoginProps{
		Email:      email,
		Password:   password,
		RememberMe: rememberMe != "",
	}

	defer pages.LoginForm(&props).Render(r.Context(), w)

	valid, err := h.AuthService.VerifyPassword(password, email)
	if err != nil {
		log.Println(err)
		props.ServerError = true
		return
	}

	if !valid {
		props.Invalid = true
		return
	}

	u, err := h.UserService.GetUserByEmail(email)
	if err != nil {
		log.Println(err)
		props.ServerError = true
	} else {
		h.Sessions.Put(r.Context(), "user", u)
		h.Sessions.RememberMe(r.Context(), props.RememberMe)

		w.Header().Add("HX-Redirect", "/")
	}
}

func (h *AuthHandler) PostSignup(w http.ResponseWriter, r *http.Request) {
	fname := r.FormValue("first_name")
	lname := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	props := pages.SignupProps{}
	props.Values.FirstName = fname
	props.Values.LastName = lname
	props.Values.Email = email
	props.Values.Password = password

	defer pages.SignupForm(&props).Render(r.Context(), w)

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

	if props.Errors.Email != "" ||
		props.Errors.FirstName != "" ||
		props.Errors.LastName != "" ||
		props.Errors.Password != "" {
		return
	}

	hash := h.AuthService.GenerateHash(password)
	ou := &services.OnboardingUser{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  hash,
	}

	t, err := h.AuthService.GenerateSignupToken(ou)
	if err != nil {
		log.Println(err)
		props.ServerError = true
		return
	}

	url := "https://" + r.Host + "/auth/verify?t=" + url.QueryEscape(t)
	err = smtp.SendMail(
		"smtp.eu.mailgun.org:587",
		h.SMTPAuth,
		"sakith@hcblk.live",
		[]string{email},
		[]byte(fmt.Sprintf(
			"To: %s\r\n"+
				"Subject: Activate your hcblk.live account\r\n"+
				"\r\n"+
				"Hello %s,\r\n"+
				"Visit this link to verify: %s\r\n"+
				"Thank you for using hcblk.live\r\n",
			email,
			fname,
			url,
		)),
	)

	if err != nil {
		log.Println(err)
		props.ServerError = true
	} else {
		props.Emailed = true
	}
}

func (h *AuthHandler) PostForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	props := pages.ForgotPasswordProps{
		Email: email,
	}

	defer pages.ForgotPasswordForm(&props).Render(r.Context(), w)

	if email == "" {
		props.EmailError = true
		return
	}

	exists, err := h.UserService.UserExists(email)
	if err != nil {
		log.Println(err)
		props.ServerError = true
		return
	} else if !exists {
		props.Success = true
		return
	}

	t, err := h.AuthService.GenerateResetToken(email)
	if err != nil {
		log.Println(err)
		props.ServerError = true
		return
	}

	url := "https://" + r.Host + "/auth/reset?t=" + url.QueryEscape(t)
	err = smtp.SendMail(
		"smtp.eu.mailgun.org:587",
		h.SMTPAuth,
		"sakith@hcblk.live",
		[]string{email},
		[]byte(fmt.Sprintf(
			"To: %s\r\n"+
				"Subject: Reset your hcblk.live account password\r\n"+
				"\r\n"+
				"Hello user,\r\n"+
				"Visit this link to reset your password: %s\r\n"+
				"Thank you for using hcblk.live\r\n",
			email,
			url,
		)),
	)

	if err != nil {
		log.Println(err)
		props.ServerError = true
	} else {
		props.Success = true
	}
}

func (h *AuthHandler) PostReset(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")

	props := pages.ResetProps{
		Password: password,
	}

	defer pages.ResetForm(&props).Render(r.Context(), w)

	if password == "" {
		props.PasswordError = "This field is required"
		return
	} else if len(password) < 8 || len(password) > 64 {
		props.PasswordError = "The password must be between 8-64 characters"
		return
	}

	referer, err := url.Parse(r.Referer())
	if err != nil {
		log.Println(err)
		props.ServerError = true
		return
	}

	if referer.Host != r.Host || referer.Path != r.URL.Path {
		props.ServerError = true
		return
	}

	t := referer.Query().Get("t")
	if t == "" {
		props.ServerError = true
		return
	}

	email, err := h.AuthService.VerifyResetToken(t)
	if err != nil {
		log.Println(err)
		props.ServerError = true
		return
	}

	u, err := h.UserService.GetUserByEmail(email)
	if err != nil {
		log.Println(err)
		props.ServerError = true
		return
	}

	err = h.AuthService.ChangePassword(u.Id, password)
	if err != nil {
		log.Println(err)
		props.ServerError = true
		return
	}

	err = h.AuthService.DeleteResetToken(t);
	if err != nil {
		log.Println(err)
		props.ServerError = true
		return
	}

	w.Header().Add("HX-Redirect", "/")
}
