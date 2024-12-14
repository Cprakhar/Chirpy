package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(passwd string) (string, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("couldn't hash the password")
	}
	return string(hashed_password), nil
}

func CheckHashPassword(password , hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password is incorrect")
	}
	return nil
}