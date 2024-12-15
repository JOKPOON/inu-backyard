package utils

import (
	"math/rand"
	"strconv"

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

func ValidatePagination(pageIndex string, pageSize string) (int, int, error) {
	page, _ := strconv.Atoi(pageIndex)
	if page <= 0 {
		page = 1
	}

	size, _ := strconv.Atoi(pageSize)
	switch {
	case size > 100:
		size = 100
	case size <= 0:
		size = 2
	}

	offset := (page - 1) * size

	return offset, size, nil
}

func GenerateRandomInt(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
