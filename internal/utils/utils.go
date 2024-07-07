package utils

import (
	"context"
	"crypto/rand"
	"image"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sakithb/hcblk-server/internal/models"
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

func FormatPhoneNo(no string) string {
	return no[:3] + "-" + no[3:6] + "-" + no[6:]
}

func GetConditionString(used bool) string {
	if used {
		return "Used"
	} else {
		return "Brand new"
	}
}

func GetUserFromContext(ctx context.Context) *models.User {
	u, ok := ctx.Value("user").(models.User)
	if !ok {
		return nil
	} else {
		return &u
	}
}

func CropImageToSquare(img image.Image) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	xs := 0
	ys := 0

	var crect image.Rectangle

	if w == h {
		return img
	} else if w < h {
		crect = image.Rect(0, 0, w, w)
		ys = (h / 2) - (w / 2)
	} else {
		crect = image.Rect(0, 0, h, h)
		xs = (w / 2) - (h / 2)
	}

	c := image.NewRGBA(crect)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x >= xs && x <= (x + crect.Dx()) && y >= ys && y <= (y + crect.Dy()) {
				c.Set(x - xs, y - ys, img.At(x, y))
			}
		}
	}

	return c
}
