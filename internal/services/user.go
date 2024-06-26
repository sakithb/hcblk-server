package services

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sakithb/hcblk-server/internal/models"
)

type UserService struct {
	DB *sqlx.DB
}

func (s *UserService) GetUserById(id string) (*models.User, error) {
	u := models.User{}
	err := s.DB.Get(&u, "SELECT id, first_name, last_name, email, joined_at FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	u := models.User{}
	err := s.DB.Get(&u, "SELECT id, first_name, last_name, email, joined_at FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UserService) CreateUser(fname string, lname string, email string, hash string) error {
	id, err := uuid.NewRandom()

	if err == nil {
		_, err = s.DB.Exec(
			"INSERT INTO users(id, first_name, last_name, email, password) VALUES(?,?,?,?,?)",
			id.String(),
			fname,
			lname,
			email,
			hash,
		)
	}

	return err
}

func (s *UserService) UserExists(email string) (bool, error) {
	var count int
	err := s.DB.Get(&count, "SELECT COUNT(id) FROM users WHERE email = ?", email)

	return count > 1, err
}
