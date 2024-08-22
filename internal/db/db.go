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
			id         VARCHAR(36) NOT NULL,
			first_name VARCHAR(30) NOT NULL,
			last_name  VARCHAR(30),
			email      VARCHAR(255) UNIQUE,
			password   VARCHAR(120) NOT NULL,
			joined_at  DATETIME NOT NULL DEFAULT(unixepoch()),
			PRIMARY KEY(id)
		)
	`

	LISTINGS_SCHEMA = `
		CREATE TABLE IF NOT EXISTS
		listings(
			id    	    VARCHAR(36) NOT NULL,
			seller_id   VARCHAR(36) NOT NULL,
			bike_id  VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			price 	    INT NOT NULL,
			mileage     INT NOT NULL,
			used        BOOLEAN NOT NULL,
			city_id     INT NOT NULL,
			phone_nos   VARCHAR(255) NOT NULL,
			listed_at   DATETIME NOT NULL DEFAULT(unixepoch()),
			PRIMARY KEY(id),
			FOREIGN KEY(seller_id) REFERENCES users(id),
			FOREIGN KEY(bike_id) REFERENCES bikes(id),
			FOREIGN KEY(city_id) REFERENCES cities(id)
		)
	`

	SIGNUP_TOKENS_SCHEMA = `
		CREATE TABLE IF NOT EXISTS
		signup_tokens(
			token      VARCHAR(50) NOT NULL,
			first_name VARCHAR(30) NOT NULL,
			last_name  VARCHAR(30),
			email      VARCHAR(255) NOT NULL,
			password   VARCHAR(120) NOT NULL,
			PRIMARY KEY(token)
		)
	`

	RESET_TOKENS_SCHEMA = `
		CREATE TABLE IF NOT EXISTS
		reset_tokens(
			token      VARCHAR(50) NOT NULL,
			email      VARCHAR(255) NOT NULL,
			PRIMARY KEY(token)
		)
	`

	SESSIONS_SCHEMA = `
		CREATE TABLE IF NOT EXISTS
		sessions (
			token  TEXT PRIMARY KEY,
			data   BLOB NOT NULL,
			expiry REAL NOT NULL
		);

		CREATE INDEX IF NOT EXISTS
		sessions_expiry_idx ON sessions(expiry);
	`
)

func New() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "assets/db/data.db")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(USERS_SCHEMA)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(LISTINGS_SCHEMA)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(SIGNUP_TOKENS_SCHEMA)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(RESET_TOKENS_SCHEMA)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(SESSIONS_SCHEMA)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
