package models

import "time"

type Seller User

type Listing struct {
	Id          string    `db:"id"`
	Description string    `db:"description"`
	Price       int       `db:"price"`
	Mileage     int       `db:"mileage"`
	Used        bool      `db:"used"`
	PhoneNosRaw string    `db:"phone_nos"`
	ListedAt    time.Time `db:"listed_at"`

	Images      []string
	PhoneNos    []string

	*Bike       `db:"bike"`
	*City       `db:"city"`
	*Seller     `db:"seller"`
}
