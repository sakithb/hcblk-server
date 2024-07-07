package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/sakithb/hcblk-server/internal/models"
)

type UIService struct {
	DB *sqlx.DB
}



func (s *UIService) GetBikeById(id string) (*models.Bike, error) {
	var b models.Bike
	err := s.DB.Get(&b, "SELECT * FROM bikes WHERE id = ?", id)

	return &b, err
}

func (s *UIService) GetBikeIdByBrandModelYear(brand string, model string, year int) (string, error) {
	var id string
	err := s.DB.Get(&id, "SELECT id FROM bikes WHERE brand = ? AND model = ? AND year = ?", brand, model, year)

	return id, err
}

func (s *UIService) ListBikeCategories() ([]string, error) {
	cats := []string{}
	err := s.DB.Select(&cats, "SELECT DISTINCT category FROM bikes")

	return cats, err
}

func (s *UIService) ListBikeBrands() ([]string, error) {
	brands := []string{}
	err := s.DB.Select(&brands, "SELECT DISTINCT brand FROM bikes")

	return brands, err
}

func (s *UIService) ListBikeModelsByBrand(brand string) ([]string, error) {
	models := []string{}
	err := s.DB.Select(&models, "SELECT model FROM bikes WHERE brand = ?", brand)

	return models, err
}

func (s *UIService) ListBikeYearsByBrandAndModel(brand string, model string) ([]string, error) {
	years := []string{}
	err := s.DB.Select(&years, "SELECT year FROM bikes WHERE brand = ? AND model = ?", brand, model)

	return years, err
}

func (s *UIService) GetCityById(id string) (*models.City, error) {
	var l models.City
	err := s.DB.Get(&l, "SELECT * FROM cities WHERE id = ?", id)

	return &l, err
}

func (s *UIService) GetCityIdByCityDistrictProvince(city string, district string, province string) (string, error) {
	var id string
	err := s.DB.Get(&id, "SELECT id FROM cities WHERE city = ? AND district = ? AND province = ?", city, district, province)

	return id, err
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
	err := s.DB.Select(&cities, "SELECT city FROM cities WHERE province = ? AND district = ?", province, district)

	return cities, err
}
