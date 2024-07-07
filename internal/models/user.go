package models

import "time"

type User struct {
	Id        string    `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	JoinedAt  time.Time `db:"joined_at"`
}
