package routes

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sakithb/hcblk-server/internal/services"
	"github.com/sakithb/hcblk-server/internal/templates/pages"
	"github.com/sakithb/hcblk-server/internal/utils"
)

type IndexHandler struct {
	ListingService *services.ListingService
	UIService      *services.UIService
}

func NewIndexHandler(ls *services.ListingService, us *services.UIService) *IndexHandler {
	return &IndexHandler{
		ListingService: ls,
		UIService:      us,
	}
}

func (h *IndexHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.Get)

	return r
}

func (h *IndexHandler) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	oobModels := q.Get("brand") != ""
	oobYears := q.Get("model") != ""
	oobDistricts := q.Get("province") != ""
	oobCities := q.Get("district") != ""

	if p := r.Header.Get("HX-Current-URL"); p != "" {
		u, err := url.Parse(p)
		if err == nil {
			for k := range u.Query() {
				if !q.Has(k) {
					q.Add(k, u.Query().Get(k))
				}
			}
		}
	}

	query := q.Get("q")

	sortBy := q.Get("sort_by")
	price := q.Get("price")
	mileage := q.Get("mileage")
	used := q.Get("used")
	brand := q.Get("brand")
	model := q.Get("model")
	year := q.Get("year")
	category := q.Get("category")
	engineCapacity := q.Get("engine_capacity")
	city := q.Get("city")
	district := q.Get("district")
	province := q.Get("province")

	opts := &services.SearchOptions{
		FilterBy: make(map[services.FilterByAttr]interface{}),
		Limit:    10,
	}

	props := &pages.IndexProps{}
	props.Values.SortBy = sortBy
	props.Values.Category = category
	props.Values.Brand = brand
	props.Values.Model = model
	props.Values.Year = year
	props.Values.City = city
	props.Values.District = district
	props.Values.Province = province
	props.Values.Used = used

	if sortBy == "price_asc" {
		opts.SortBy = services.SORT_BY_PRICE
		opts.SortIn = services.SORT_IN_ASC
	} else if sortBy == "price_desc" {
		opts.SortBy = services.SORT_BY_PRICE
		opts.SortIn = services.SORT_IN_DESC
	} else if sortBy == "mileage_asc" {
		opts.SortBy = services.SORT_BY_MILEAGE
		opts.SortIn = services.SORT_IN_ASC
	} else if sortBy == "mileage_desc" {
		opts.SortBy = services.SORT_BY_MILEAGE
		opts.SortIn = services.SORT_IN_DESC
	} else if sortBy == "listed_at_asc" {
		opts.SortBy = services.SORT_BY_LISTED_AT
		opts.SortIn = services.SORT_IN_ASC
	} else {
		opts.SortBy = services.SORT_BY_LISTED_AT
		opts.SortIn = services.SORT_IN_DESC
	}

	priceRange := strings.Split(price, ",")
	if len(priceRange) == 2 {
		priceMin, err := strconv.Atoi(priceRange[0])
		if err != nil {
			priceMin = -1
		}

		priceMax, err := strconv.Atoi(priceRange[1])
		if err != nil {
			priceMax = -1
		}

		opts.FilterBy[services.FILTER_BY_PRICE] = services.Range{Min: priceMin, Max: priceMax}
		props.Values.Price.Min = priceRange[0]
		props.Values.Price.Min = priceRange[1]
	}

	mileageRange := strings.Split(mileage, ",")
	if len(mileageRange) == 2 {
		mileageMin, err := strconv.Atoi(mileageRange[0])
		if err != nil {
			mileageMin = -1
		}

		mileageMax, err := strconv.Atoi(mileageRange[1])
		if err != nil {
			mileageMax = -1
		}

		opts.FilterBy[services.FILTER_BY_MILEAGE] = services.Range{Min: mileageMin, Max: mileageMax}
		props.Values.Mileage.Min = mileageRange[0]
		props.Values.Mileage.Min = mileageRange[1]
	}

	engineCapacityRange := strings.Split(engineCapacity, ",")
	if len(engineCapacityRange) == 2 {
		engineCapacityMin, err := strconv.Atoi(engineCapacityRange[0])
		if err != nil {
			engineCapacityMin = -1
		}

		engineCapacityMax, err := strconv.Atoi(engineCapacityRange[1])
		if err != nil {
			engineCapacityMax = -1
		}

		opts.FilterBy[services.FILTER_BY_ENGINE_CAPACITY] = services.Range{Min: engineCapacityMin, Max: engineCapacityMax}
		props.Values.EngineCapacity.Min = engineCapacityRange[0]
		props.Values.EngineCapacity.Min = engineCapacityRange[1]
	}

	if used == "used" {
		opts.FilterBy[services.FILTER_BY_USED] = true
	} else if used == "brand_new" {
		opts.FilterBy[services.FILTER_BY_USED] = false
	}

	if brand != "" {
		opts.FilterBy[services.FILTER_BY_BRAND] = brand
	}

	if model != "" {
		opts.FilterBy[services.FILTER_BY_MODEL] = model
	}

	if year != "" {
		opts.FilterBy[services.FILTER_BY_YEAR] = year
	}

	if category != "" {
		opts.FilterBy[services.FILTER_BY_CATEGORY] = category
	}

	if city != "" {
		opts.FilterBy[services.FILTER_BY_CITY] = city
	}

	if district != "" {
		opts.FilterBy[services.FILTER_BY_DISTRICT] = district
	}

	if province != "" {
		opts.FilterBy[services.FILTER_BY_PROVINCE] = province
	}

	ls, err := h.ListingService.SearchListings(query, opts)
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	cats, err := h.UIService.ListBikeCategories()
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	brands, err := h.UIService.ListBikeBrands()
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	provinces, err := h.UIService.ListProvinces()
	if err != nil {
		utils.HandleServerError(w, err)
		return
	}

	props.Listings = ls
	props.Categories = cats
	props.Brands = brands
	props.Provinces = provinces

	if r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-Boosted") != "true" {
		rprops := pages.IndexResultsProps{
			Listings: ls,
		}

		if brand == "" {
			rprops.Models = []string{}
			rprops.Years = []string{}
		} else if model == "" {
			rprops.Years = []string{}
		}

		if province == "" {
			rprops.Districts = []string{}
			rprops.Cities = []string{}
		} else if district == "" {
			rprops.Cities = []string{}
		}

		if oobModels {
			models, err := h.UIService.ListBikeModelsByBrand(brand)
			if err != nil {
				utils.HandleServerError(w, err)
				return
			}

			rprops.Models = models
			rprops.Years = []string{}
		}

		if oobYears && brand != "" {
			years, err := h.UIService.ListBikeYearsByBrandAndModel(brand, model)
			if err != nil {
				utils.HandleServerError(w, err)
				return
			}

			rprops.Years = years
		}

		if oobDistricts {
			districts, err := h.UIService.ListDistrictsByProvince(province)
			if err != nil {
				utils.HandleServerError(w, err)
				return
			}

			rprops.Districts = districts
			rprops.Cities = []string{}
		}

		if oobCities {
			cities, err := h.UIService.ListCitiesByDistrictAndProvince(province, district)
			if err != nil {
				utils.HandleServerError(w, err)
				return
			}

			rprops.Cities = cities
		}

		w.Header().Add("HX-Push-URL", "/?"+q.Encode())
		pages.IndexResults(&rprops).Render(r.Context(), w)
	} else {
		if brand != "" {
			models, err := h.UIService.ListBikeModelsByBrand(brand)
			if err != nil {
				utils.HandleServerError(w, err)
				return
			}

			props.Models = models
		}

		if brand != "" && model != "" {
			years, err := h.UIService.ListBikeYearsByBrandAndModel(brand, model)
			if err != nil {
				utils.HandleServerError(w, err)
				return
			}

			props.Years = years
		}

		pages.Index(props).Render(context.WithValue(r.Context(), "search", true), w)
	}
}
