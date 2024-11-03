package utils

import (
	errs "github.com/team-inu/inu-backyard/entity/error"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.New(errs.ErrUserPassword, "cannot hash password", err)
	}
	return string(bcryptPassword), nil
}

func CheckPassword(hash string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errs.New(errs.ErrUserPassword, "password is incorrect")
	}
	return nil
}
