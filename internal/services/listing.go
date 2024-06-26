package services

import (
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sakithb/hcblk-server/internal/models"
)

type ListingService struct {
	DB *sqlx.DB
}

func (s *ListingService) GetListingById(id string) (*models.Listing, error) {
	l := &models.Listing{}
	err := s.DB.Get(&l, "SELECT id, seller_id, model_slug, description, price, mileage, used, location_slug, phone_nos, listed_at  FROM listings WHERE id = ?", id)

	return l, err
}

func (s *ListingService) GetListingsBySeller(sellerId string) ([]models.Listing, error) {
	l := []models.Listing{}
	err := s.DB.Select(&l, "SELECT id, seller_id, model_slug, description, price, mileage, used, location_slug, phone_nos, listed_at  FROM listings WHERE seller_id = ?", sellerId)

	return l, err
}

func (s *ListingService) CreateListing(sellerId string, modelSlug string, description string, price uint32, mileage uint32, used bool, locationSlug string, phoneNos []string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec("INSERT INTO listings(id, seller_id, model_slug, description, price, mileage, used, location_slug, phone_nos) VALUES(?,?,?,?,?,?,?,?,?)", id, sellerId, modelSlug, description, price, mileage, used, locationSlug, strings.Join(phoneNos, ","))
	if err != nil {
		return err
	}

	return nil
}
