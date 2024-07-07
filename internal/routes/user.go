package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sakithb/hcblk-server/internal/models"
	"github.com/sakithb/hcblk-server/internal/services"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
	"github.com/sakithb/hcblk-server/internal/utils"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(us *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

func (h *UserHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/{id}", h.GetUser)

	return r
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.HandleHTTPCode(w, 400)
		return
	}

	var u *models.User
	if id == "me" {
		ut, ok := r.Context().Value("user").(models.User)
		if !ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		u = &ut
	} else {
		ut, err := h.UserService.GetUserById(id)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.HandleHTTPCode(w, 400)
			} else {
				utils.HandleServerError(w, err)
			}

			return
		}

		u = ut
	}


	pages.Profile(u).Render(r.Context(), w)
}
