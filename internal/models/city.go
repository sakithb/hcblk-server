package models

type City struct {
	Id       string `db:"id"`
	City     string `db:"city"`
	District string `db:"district"`
	Province string `db:"province"`
}
