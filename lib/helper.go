package lib

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashAndPassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func GenerateUUID() string {
	return uuid.NewString()
}
