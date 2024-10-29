package mocks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/team-inu/inu-backyard/entity"
)

// Mocking AuthUseCase
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) ChangePassword(userId, oldPassword, newPassword string) error {
	args := m.Called(userId, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockAuthUseCase) Authenticate(sessionId string) (*entity.User, error) {
	args := m.Called(sessionId)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthUseCase) SignIn(email string, password string, ipAddress string, userAgent string) (*fiber.Cookie, error) {
	args := m.Called(email, password, ipAddress, userAgent)
	return args.Get(0).(*fiber.Cookie), args.Error(1)
}

func (m *MockAuthUseCase) SignOut(header string) (*fiber.Cookie, error) {
	args := m.Called(header)
	return args.Get(0).(*fiber.Cookie), args.Error(1)
}
