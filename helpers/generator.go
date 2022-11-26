package helpers

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var APPLICATION_NAME = "ima-svc-management"
var LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour

func GeneratePasswordHash(password []byte) string {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	return string(hashedPassword)
}

func PasswordCompare(password []byte, hashedPassword []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return false, err
	}
	return true, nil
}

// TODO change logic jwt
func GenerateToken(param string) (string, error) {
	key := []byte(param)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(LOGIN_EXPIRATION_DURATION)),
		Issuer:    "ima-svc-management",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return signedString, nil
}
