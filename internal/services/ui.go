package services

import "github.com/jmoiron/sqlx"

type UIService struct {
	DB *sqlx.DB
}

func (s *UIService) ListBikeBrands() ([]string, error) {
	brands := []string{}
	err := s.DB.Select(&brands, "SELECT DISTINCT brand FROM models")

	return brands, err
}

func (s *UIService) ListBikeModelsByBrand(brand string) ([]string, error) {
	models := []string{}
	err := s.DB.Select(&models, "SELECT model FROM models WHERE brand = ?", brand)

	return models, err
}

func (s *UIService) ListBikeYearsByBrandAndModel(brand string, model string) ([]string, error) {
	years := []string{}
	err := s.DB.Select(&years, "SELECT year FROM models WHERE brand = ? AND model = ?", brand)

	return years, err
}

func (s *UIService) ListProvinces() ([]string, error) {
	provinces := []string{}
	err := s.DB.Select(&provinces, "SELECT DISTINCT province FROM cities")

	return provinces, err
}

func (s *UIService) ListDistrictsByProvince(province string) ([]string, error) {
	districts := []string{}
	err := s.DB.Select(&districts, "SELECT district FROM cities WHERE province = ?", province)

	return districts, err
}

func (s *UIService) ListCitiesByDistrictAndProvince(province string, district string) ([]string, error) {
	cities := []string{}
	err := s.DB.Select(&cities, "SELECT city FROM cities WHERE province = ? AND district = ?", district)

	return cities, err
}

