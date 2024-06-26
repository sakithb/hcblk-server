package models

import (
	"encoding/gob"
	"time"
)

type User struct {
	Id        string    `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	JoinedAt  time.Time `db:"joined_at"`
}

type OnboardingUser struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Hash      string `db:"hash"`
}

type Bike struct {
	Slug           string `db:"slug"`
	Model          string `db:"model"`
	Brand          string `db:"brand"`
	Category       string `db:"category"`
	Year           uint16 `db:"year"`
	EngineCapacity uint16 `db:"engine_capacity"`
}

type Listing struct {
	Id           string    `db:"id"`
	SellerId     string    `db:"seller_id"`
	ModelSlug    string    `db:"model_slug"`
	Description  string    `db:"description"`
	Price        uint32    `db:"price"`
	Mileage      uint32    `db:"mileage"`
	Used         bool      `db:"used"`
	LocationSlug uint16    `db:"location_slug"`
	PhoneNos     []string  `db:"phone_nos"`
	ListedAt     time.Time `db:"listed_at"`
	Images       []string
}

func init() {
	gob.Register(User{})
}
