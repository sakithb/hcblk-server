package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	USERS_SCHEMA = `
		CREATE TABLE IF NOT EXISTS
		users(
			id         TEXT PRIMARY KEY,
			first_name TEXT,
			last_name  TEXT,
			email      TEXT UNIQUE,
			password   TEXT,
			created_at INTEGER
		)
	`
)

func New() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "assets/data.db")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(USERS_SCHEMA)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
