package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password []byte) string {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	return string(hashedPassword)
}

func PasswordCompare(password []byte, hashedPassword []byte) error {
    err := bcrypt.CompareHashAndPassword(hashedPassword, password)

    return err
}
