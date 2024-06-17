package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/crypto/argon2"
)

func GenerateRandomBytes(len int) []byte  {
	b := make([]byte, len)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln(err)
	}

	return b
}


func GenerateHashFromPassword(pwd []byte) string {
	const (
		m uint32 = 46 * 1024;
		t uint32 = 1;
		p uint8 = 1;
	)

	salt := GenerateRandomBytes(16)
	hash := argon2.IDKey(pwd, salt, t, m, p, 32)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, m, t, p, b64Salt, b64Hash)
	return encodedHash
}
