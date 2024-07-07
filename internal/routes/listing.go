package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sakithb/hcblk-server/internal/services"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
	"github.com/sakithb/hcblk-server/internal/utils"
)

type ListingHandler struct {
	ListingService *services.ListingService
}

func NewListingHandler(ls *services.ListingService) *ListingHandler {
	return &ListingHandler{
		ListingService: ls,
	}
}

func (h *ListingHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/{id}", h.GetListingId)

	return r
}

func (h *ListingHandler) GetListingId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := uuid.Validate(id)
	if err != nil {
		utils.HandleHTTPCode(w, 400)
		return
	}

	l, err := h.ListingService.GetListingById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.HandleHTTPCode(w, 400)
			return
		} else {
			utils.HandleServerError(w, err)
			return
		}
	}

	pages.Listing(l).Render(r.Context(), w)
}
