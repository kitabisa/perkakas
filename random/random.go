package random

import (
	"math/rand"

	"github.com/google/uuid"
)

func UUID() (id string, err error) {
	newId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return newId.String(), nil
}

// String random string with n length
func String(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
