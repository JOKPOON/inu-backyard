package mocks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

// Mocking PayloadValidator
type MockValidator struct {
	mock.Mock
}

// ValidateAuth implements validator.PayloadValidator.
func (m *MockValidator) ValidateAuth(ctx *fiber.Ctx) (string, error) {
	panic("unimplemented")
}

func (m *MockValidator) Validate(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	args := m.Called(payload, ctx)
	return args.Bool(0), args.Error(1)
}
