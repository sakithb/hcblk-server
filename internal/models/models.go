package models

import "time"

type User struct {
	Id        string `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

type Bike struct {
	Id             string `db:"	id"`
	Model          string `db:"model"`
	Brand          string `db:"brand"`
	Category       string `db:"category"`
	Year           int    `db:"year"`
	EngineCapacity int    `db:"engine_capacity"`
}

type Listing struct {
	Id       string    `db:"	id"`
	Seller   *User     `db:"seller"`
	Model    *Bike     `db:"model"`
	Price    int       `db:"price"`
	District string    `db:"district"`
	City     string    `db:"city"`
	ListedAt time.Time `db:"listed_at"`
}
