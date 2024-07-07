package models

type Bike struct {
	Id             string `db:"id"`
	Model          string `db:"model"`
	Brand          string `db:"brand"`
	Category       string `db:"category"`
	Year           int `db:"year"`
	EngineCapacity int `db:"engine_capacity"`
}
