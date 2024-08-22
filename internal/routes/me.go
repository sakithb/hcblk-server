package routes

import (
	"database/sql"
	"image"
	"image/png"
	_ "image/jpeg"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sakithb/hcblk-server/internal/models"
	"github.com/sakithb/hcblk-server/internal/services"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
	"github.com/sakithb/hcblk-server/internal/utils"
)

type MeHandler struct {
	UserService    *services.UserService
	UIService      *services.UIService
	ListingService *services.ListingService
	AuthService    *services.AuthService
}

func NewMeHandler(us *services.UserService, is *services.UIService, ls *services.ListingService, as *services.AuthService) *MeHandler {
	return &MeHandler{
		UserService:    us,
		UIService:      is,
		ListingService: ls,
		AuthService:    as,
	}
}

func (h *MeHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/settings", h.GetSettings)
	r.Get("/listings", h.GetListings)
	r.Get("/listings/create", h.GetCreateListing)

	r.Get("/listings/create/models", h.GetCreateListingModels)
	r.Get("/listings/create/years", h.GetCreateListingYears)
	r.Get("/listings/create/districts", h.GetCreateListingDistricts)
	r.Get("/listings/create/cities", h.GetCreateListingCities)

	r.Post("/listings/create", h.PostCreateListing)
	r.Post("/settings/password", h.PostSettingsPassword)
	r.Post("/settings/photo", h.PostSettingsPhoto)

	return r
}

func (h *MeHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	pages.UserSettings().Render(r.Context(), w)
}

func (h *MeHandler) GetCreateListing(w http.ResponseWriter, r *http.Request) {
	brands, err := h.UIService.ListBikeBrands()
	if err != nil {
		utils.HandleServerError(w, err)
	}

	provinces, err := h.UIService.ListProvinces()
	if err != nil {
		utils.HandleServerError(w, err)
	}

	pages.CreateListing(brands, provinces).Render(r.Context(), w)
}

func (h *MeHandler) GetCreateListingModels(w http.ResponseWriter, r *http.Request) {
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

func (h *MeHandler) GetCreateListingYears(w http.ResponseWriter, r *http.Request) {
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

func (h *MeHandler) GetCreateListingDistricts(w http.ResponseWriter, r *http.Request) {
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

func (h *MeHandler) GetCreateListingCities(w http.ResponseWriter, r *http.Request) {
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

func (h *MeHandler) GetListings(w http.ResponseWriter, r *http.Request) {
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

	pages.UserListings(ls).Render(r.Context(), w)
}

func (h *MeHandler) PostCreateListing(w http.ResponseWriter, r *http.Request) {
	brand := r.FormValue("brand")
	model := r.FormValue("model")
	year := r.FormValue("year")
	province := r.FormValue("province")
	district := r.FormValue("district")
	city := r.FormValue("city")
	price := r.FormValue("price")
	mileage := r.FormValue("mileage")
	description := r.FormValue("description")
	condition := r.FormValue("condition")
	phoneNos := strings.Split(r.FormValue("phone_nos"), ",")

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

	props := pages.CreateListingFormProps{}
	props.Values.Brand = brand
	props.Values.Model = model
	props.Values.Year = year
	props.Values.Province = province
	props.Values.District = district
	props.Values.City = city
	props.Values.Price = price
	props.Values.Mileage = mileage
	props.Values.Description = description
	props.Values.Condition = condition
	props.Values.PhoneNos = phoneNos

	yearInt, err := strconv.Atoi(props.Values.Year)
	if err != nil || yearInt < 0 {
		props.Errors.Model = "Please enter a valid year"
	}

	mileageInt, err := strconv.Atoi(props.Values.Mileage)
	if err != nil || mileageInt < 0 {
		props.Errors.Price = "Please enter a valid price"
	}

	priceInt, err := strconv.Atoi(props.Values.Price)
	if err != nil || priceInt < 0 {
		props.Errors.Price = "Please enter a valid price"
	}

	for _, phoneNo := range phoneNos {
		if phoneNo == "" || len(phoneNo) != 10 || phoneNo[0] != '0' {
			props.Errors.PhoneNos = "One or more phone numbers are invalid"
		}
	}

	if condition != "New" && condition != "Used" {
		props.Errors.Condition = "Please select a valid condition"
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
			mt := mime.TypeByExtension(path.Ext(img.Filename))
			if mt != "image/png" && mt != "image/jpeg" {
				props.Errors.Images = "Images should be in png or jpeg format"
				break
			} else if img.Size > 1024*1024*3 {
				props.Errors.Images = "Please make sure each image is less than 3MB"
				break
			}
		}
	}

	if props.Errors.Model != "" || props.Errors.Price != "" || props.Errors.Location != "" || props.Errors.Images != "" {
		pages.CreateListingForm(brands, provinces, &props).Render(r.Context(), w)
	} else {
		u, ok := r.Context().Value("user").(models.User)
		if !ok {
			w.Header().Add("HX-Redirect", "/auth/login")
			return
		}

		bikeId, err := h.UIService.GetBikeIdByBrandModelYear(brand, model, yearInt)
		if err != nil {
			utils.HandleServerError(w, err)
			return
		}

		cityId, err := h.UIService.GetCityIdByCityDistrictProvince(city, district, province)
		if err != nil {
			utils.HandleServerError(w, err)
			return
		}

		readers := []io.Reader{}
		for _, img := range images {
			f, err := img.Open()
			if err != nil {
				utils.HandleServerError(w, err)
				return
			}

			readers = append(readers, f)
		}

		err = h.ListingService.CreateListing(u.Id, bikeId, description, priceInt, mileageInt, condition == "Used", cityId, phoneNos, readers)
		if err != nil {
			utils.HandleServerError(w, err)
			return
		}

		w.Header().Add("HX-Redirect", "/me/listings")
	}
}

func (h *MeHandler) PostSettingsPassword(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	oldPwd := r.FormValue("old_password")
	newPwd := r.FormValue("new_password")
	confirmPwd := r.FormValue("confirm_password")

	props := &pages.UserSettingsChangePasswordFormProps{}
	props.OldPassword = oldPwd
	props.NewPassword = newPwd
	props.ConfirmPassword = confirmPwd

	if oldPwd != "" && newPwd != "" && confirmPwd != "" {
		if newPwd != confirmPwd {
			props.Error = "Passwords do not match"
		} else if len(newPwd) < 8 || len(newPwd) > 64 {
			props.Error = "The password must be between 8-64 characters"
		} else {
			ok, err := h.AuthService.VerifyPassword(oldPwd, u.Email)
			if err != nil {
				if err == sql.ErrNoRows {
					props.Error = "Old password is incorrect"
				} else {
					utils.HandleServerError(w, err)
					return
				}
			} else if !ok {
				props.Error = "Old password is incorrect"
			} else {
				err = h.AuthService.ChangePassword(u.Id, newPwd)
				if err != nil {
					utils.HandleServerError(w, err)
					return
				}

				props.Success = true
			}
		}

	} else {
		props.Error = "All fields are required"
	}

	pages.UserSettingsChangePasswordForm(props).Render(r.Context(), w)
}

func (h *MeHandler) PostSettingsPhoto(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	errr := ""
	success := false

	file, header, err := r.FormFile("profile_pic")
	if err != nil {
		if err == http.ErrMissingFile {
			errr = "Please select a photo"
		} else {
			utils.HandleServerError(w, err)
			return
		}
	}

	mt := mime.TypeByExtension(path.Ext(header.Filename))

	if mt != "image/png" && mt != "image/jpeg" {
		errr = "Images should be in png or jpeg format"
	} else if header.Size > 1024*1024*3 {
		errr = "The photo is too large"
	} else {
		img, _, err := image.Decode(file)
		if err != nil {
			utils.HandleServerError(w, err)
			return
		}

		o, err := os.Create("assets/dist/users/" + u.Id + ".png")
		if err != nil {
			utils.HandleServerError(w, err)
			return
		}

		cimg := utils.CropImageToSquare(img)
		err = png.Encode(o, cimg)
		if err != nil {
			utils.HandleServerError(w, err)
			return
		}

		o.Close()
		success = true
	}

	pages.UserSettingsProfilePhotoForm(errr, success).Render(r.Context(), w)
}
