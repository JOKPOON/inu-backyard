package mocks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/team-inu/inu-backyard/entity"
)

// Mocking UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) GetAll() ([]entity.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		if args.Error(1) != nil {
			return nil, args.Error(1)
		}
		return nil, nil
	}
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetById(userId string) (*entity.User, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetBySessionId(sessionId string) (*entity.User, error) {
	args := m.Called(sessionId)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetByParams(params *entity.User, limit, offset int) ([]entity.User, error) {
	args := m.Called(params, limit, offset)
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserUseCase) Create(firstName, lastName, email, password string, role entity.UserRole) error {
	args := m.Called(firstName, lastName, email, password, role)
	return args.Error(0)
}

func (m *MockUserUseCase) CreateMany(users []entity.User) error {
	args := m.Called(users)
	return args.Error(0)
}

func (m *MockUserUseCase) Update(userId string, user *entity.User) error {
	args := m.Called(userId, user)
	return args.Error(0)
}

func (m *MockUserUseCase) Delete(userId string) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *MockUserUseCase) CheckUserRole(ctx *fiber.Ctx, userId string, role entity.UserRole) error {
	args := m.Called(ctx, userId, role)
	return args.Error(0)
}

func (m *MockUserUseCase) GetUserFromCtx(ctx *fiber.Ctx) (*entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entity.User), args.Error(1)
}
