package services

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sakithb/hcblk-server/internal/models"
	"github.com/sakithb/hcblk-server/internal/utils"
	"golang.org/x/crypto/argon2"
)

type AuthService struct {
	DB *sqlx.DB
}

const (
	TOKEN_LENGTH = 32
)

const (
	MEMORY uint32 = 46 * 1024;
	TIME uint32 = 1;
	THREADS uint8 = 1;
	LENGTH uint32 = 32;
)

func (s *AuthService) GenerateHash(pwd string) string {
	salt := utils.GenerateRandomBytes(16)
	hash := argon2.IDKey([]byte(pwd), salt, TIME, MEMORY, THREADS, LENGTH)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, MEMORY, TIME, THREADS, b64Salt, b64Hash)
	return encodedHash
}

func (s *AuthService) VerifyPassword(pwd string, email string) (bool, error) {
	var storedHash string
	err := s.DB.Get(&storedHash, "SELECT password FROM users WHERE email = ?", email)
	if err != nil {
		return false, err
	}

	components := strings.Split(storedHash, "$")
	b64StoredHash := components[len(components) - 1]
	b64Salt := components[len(components) - 2]

	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(pwd), salt, TIME, MEMORY, THREADS, LENGTH)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return subtle.ConstantTimeCompare([]byte(b64Hash), []byte(b64StoredHash)) > 0, nil
}

func (s *AuthService) GenerateToken(u *models.OnboardingUser) (string, error) {
	bytes := utils.GenerateRandomBytes(TOKEN_LENGTH)
	token := base64.StdEncoding.EncodeToString(bytes)

	_, err := s.DB.Exec(
		"INSERT INTO tokens VALUES(?, ?, ?, ?, ?)",
		token,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Hash,
	)

	return token, err
}

func (s *AuthService) VerifyToken(token string) (*models.OnboardingUser, error) {
	u := &models.OnboardingUser{}
	err := s.DB.Get(&u, "SELECT first_name, last_name, email, hash FROM tokens WHERE token = ?", token)

	return u, err
}

func (s *AuthService) DeleteToken(token string) error {
	_, err := s.DB.Exec("DELETE FROM tokens WHERE token = ?", token)
	return err
}
