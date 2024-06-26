package utils

import (
	"crypto/rand"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GenerateRandomBytes(len int) []byte {
	b := make([]byte, len)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatalln(err)
	}

	return b
}

func HandleServerError(w http.ResponseWriter, err error) {
	log.Fatalln(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func HandleHTTPCode(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func FormatInteger(i int) string {
	str := strconv.Itoa(i)
	a := []string{}

	for i := (len(str) % 3) - 1; len(str) > 0; i = 2 {
		a = append(a, str[:i+1])
		str = str[i+1:]
	}

	return strings.Join(a, ",")
}
