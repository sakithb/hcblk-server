package routes

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sakithb/hcblk-server/internal/models"
	"github.com/sakithb/hcblk-server/internal/services"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
	"github.com/sakithb/hcblk-server/internal/utils"
)

type UserHandler struct {
	UserService    *services.UserService
	UIService      *services.UIService
	ListingService *services.ListingService
}

func NewUserHandler(us *services.UserService, is *services.UIService) *UserHandler {
	return &UserHandler{
		UserService: us,
		UIService:   is,
	}
}

func (h *UserHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/listings/create", h.CreateListingPage)
	r.Get("/listings/create/models", h.CreateListingModelsPartial)
	r.Get("/listings/create/years", h.CreateListingYearsPartial)
	r.Get("/listings/create/districts", h.CreateListingDistrictsPartial)
	r.Get("/listings/create/cities", h.CreateListingCitiesPartial)

	r.Get("/listings", h.ListingsPage)
	r.Get("/listing/{id}", h.ListingPage)
	r.Post("/listings/create", h.CreateListingFormPartial)

	return r
}

func (h *UserHandler) CreateListingPage(w http.ResponseWriter, r *http.Request) {
	brands, err := h.UIService.ListBikeBrands()
	if err != nil {
		utils.HandleServerError(w, err)
	}

	provinces, err := h.UIService.ListProvinces()
	if err != nil {
		utils.HandleServerError(w, err)
	}

	pages.CreateListing(&pages.CreateListingProps{Brands: brands, Provinces: provinces}).Render(r.Context(), w)
}

func (h *UserHandler) CreateListingModelsPartial(w http.ResponseWriter, r *http.Request) {
	brand := r.URL.Query().Get("brand")
	if brand == "" {
		utils.HandleHTTPCode(w, 400)
		return
	}

	models, err := h.UIService.ListBikeModelsByBrand(brand)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.HandleHTTPCode(w, 400)
		} else {
			utils.HandleServerError(w, err)
		}

		return
	}

	pages.CreateListingModels(brand, "", models).Render(r.Context(), w)
}

func (h *UserHandler) CreateListingYearsPartial(w http.ResponseWriter, r *http.Request) {
	brand := r.URL.Query().Get("brand")
	if brand == "" {
		utils.HandleHTTPCode(w, 400)
		return
	}

	model := r.URL.Query().Get("model")
	if model == "" {
		utils.HandleHTTPCode(w, 400)
		return
	}

	years, err := h.UIService.ListBikeYearsByBrandAndModel(brand, model)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.HandleHTTPCode(w, 400)
		} else {
			utils.HandleServerError(w, err)
		}

		return
	}

	pages.CreateListingYears(brand, model, "", years).Render(r.Context(), w)
}

func (h *UserHandler) CreateListingDistrictsPartial(w http.ResponseWriter, r *http.Request) {
	province := r.URL.Query().Get("province")
	if province == "" {
		utils.HandleHTTPCode(w, 400)
		return
	}

	districts, err := h.UIService.ListDistrictsByProvince(province)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.HandleHTTPCode(w, 400)
		} else {
			utils.HandleServerError(w, err)
		}

		return
	}

	pages.CreateListingDistricts(province, "", districts).Render(r.Context(), w)
}

func (h *UserHandler) CreateListingCitiesPartial(w http.ResponseWriter, r *http.Request) {
	province := r.URL.Query().Get("province")
	if province == "" {
		utils.HandleHTTPCode(w, 400)
		return
	}

	district := r.URL.Query().Get("district")
	if district == "" {
		utils.HandleHTTPCode(w, 400)
		return
	}

	cities, err := h.UIService.ListCitiesByDistrictAndProvince(province, district)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.HandleHTTPCode(w, 400)
		} else {
			utils.HandleServerError(w, err)
		}

		return
	}

	pages.CreateListingCities(province, district, "", cities).Render(r.Context(), w)
}

func (h *UserHandler) ListingPage(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) ListingsPage(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	ls, err := h.ListingService.GetListingsBySeller(u.Id)
	if err != nil && err != sql.ErrNoRows {
		utils.HandleServerError(w, err)
		return
	}

	pages.Listings(ls).Render(r.Context(), w)
}

func (h *UserHandler) CreateListingFormPartial(w http.ResponseWriter, r *http.Request) {
	brand := r.FormValue("brand")
	model := r.FormValue("model")
	year := r.FormValue("year")
	province := r.FormValue("province")
	district := r.FormValue("district")
	city := r.FormValue("city")
	phone := r.FormValue("phone")
	price := r.FormValue("price")

	brands, err := h.UIService.ListBikeBrands()
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	makes, err := h.UIService.ListBikeModelsByBrand(brand)
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	years, err := h.UIService.ListBikeYearsByBrandAndModel(brand, model)
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	provinces, err := h.UIService.ListProvinces()
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	districts, err := h.UIService.ListDistrictsByProvince(province)
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	cities, err := h.UIService.ListCitiesByDistrictAndProvince(province, district)
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	props := pages.CreateListingProps{Brands: brands, Provinces: provinces}
	props.Values.Brand = brand
	props.Values.Model = model
	props.Values.Year = year
	props.Values.Price = price
	props.Values.District = district
	props.Values.City = city
	props.Values.Phone = phone

	yearInt, err := strconv.Atoi(props.Values.Year)
	if err != nil || yearInt < 0 {
		props.Errors.Model = "Please enter a valid year"
	}

	priceInt, err := strconv.Atoi(props.Values.Price)
	if err != nil || priceInt < 0 {
		props.Errors.Price = "Please enter a valid price"
	}

	if phone == "" || len(phone) != 10 || phone[0] != '0' {
		props.Errors.Phone = "Please enter a valid phone number"
	}

	if !slices.Contains(brands, brand) {
		props.Errors.Model = "Please select a valid brand"
	} else if !slices.Contains(makes, model) {
		props.Errors.Model = "Please select a valid model"
	} else if !slices.Contains(years, year) {
		props.Errors.Model = "Please select a valid year"
	}

	if !slices.Contains(provinces, province) {
		props.Errors.Location = "Please select a valid province"
	} else if !slices.Contains(districts, district) {
		props.Errors.Location = "Please select a valid district"
	} else if !slices.Contains(cities, city) {
		props.Errors.Location = "Please select a valid city"
	}

	r.ParseMultipartForm(1024 * 1024 * 32)
	images := r.MultipartForm.File["images"]

	if len(images) < 3 {
		props.Errors.Images = "Please upload atleast 03 images"
	} else {
		for _, img := range images {
			if img.Size > 1024*1024*3 {
				props.Errors.Images = "Please make sure each image is less than 3MB"
				break
			}
		}
	}

	if props.Errors.Model != "" || props.Errors.Price != "" || props.Errors.Location != "" || props.Errors.Images != "" {
		pages.CreateListingForm(&props).Render(r.Context(), w)
	} else {
		u, ok := r.Context().Value("user").(models.User)
		if !ok {
			w.Header().Add("HX-Redirect", "/auth/login")
			return
		}

		err = h.ListingService.CreateListing(u.Id, )
		if err != nil {
			log.Fatalln(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		for _, img := range images {
			imgID, err := uuid.NewRandom()
			if err != nil {
				log.Fatalln(err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			s, err := img.Open()
			if err != nil {
				log.Fatalln(err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			err = os.MkdirAll("./assets/dist/listings/"+id.String()+"/", 0755)
			if err != nil {
				log.Fatalln(err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			f, err := os.Create("./assets/dist/listings/" + id.String() + "/" + imgID.String())
			if err != nil {
				log.Fatalln(err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			_, err = f.ReadFrom(s)
			if err != nil {
				log.Fatalln(err)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			f.Sync()
			s.Close()
			f.Close()
		}

		w.Header().Add("HX-Redirect", "/user/listings")
	}
}
