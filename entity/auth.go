package entity

import "github.com/gofiber/fiber/v2"

type AuthUseCase interface {
	Authenticate(header string) (*User, error)
	SignIn(payload SignInPayload, ipAddress string, userAgent string) (*fiber.Cookie, error)
	SignOut(header string) (*fiber.Cookie, error)
	ChangePassword(userId string, oldPassword string, newPassword string) error
	ForgotPassword(email string) error
	ResetPassword(email string, token string, newPassword string) error
	GetSession(email string) (string, error)
}

type SignInPayload struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}
