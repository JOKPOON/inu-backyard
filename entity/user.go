package entity

import (
	"github.com/gofiber/fiber/v2"
)

type UserRole string

const (
	UserRoleLecturer         UserRole = "LECTURER"
	UserRoleModerator        UserRole = "MODERATOR"
	UserRoleHeadOfCurriculum UserRole = "HEAD_OF_CURRICULUM"
	UserRoleTABEEManager     UserRole = "TABEE_MANAGER"
)

func (u User) IsRoles(expectedRoles []UserRole) bool {
	for _, expectedRole := range expectedRoles {
		if u.Role == expectedRole {
			return true
		}
	}

	return false
}

type UserRepository interface {
	GetAll() ([]User, error)
	GetById(id string) (*User, error)
	GetByParams(params *User, limit int, offset int) ([]User, error)
	GetByEmail(email string) (*User, error)
	GetBySessionId(sessionId string) (*User, error)
	Create(user *User) error
	CreateMany(users []User) error
	Update(id string, user *User) error
	Delete(id string) error
}

type UserUseCase interface {
	GetAll() ([]User, error)
	GetById(id string) (*User, error)
	GetByParams(params *User, limit int, offset int) ([]User, error)
	GetByEmail(email string) (*User, error)
	Create(CreateUserPayload) error
	CreateMany(users []User) error
	Update(id string, user *User) error
	Delete(id string) error
	GetBySessionId(sessionId string) (*User, error)
	CheckUserRole(ctx *fiber.Ctx, userId string, role UserRole) error
	GetUserFromCtx(ctx *fiber.Ctx) (*User, error)
}

type User struct {
	Id        string   `json:"id" gorm:"primaryKey;type:char(255)"`
	Email     string   `json:"email" gorm:"unique"`
	Password  string   `json:"password"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Role      UserRole `json:"role" gorm:"default:'LECTURER'"`
}

type CreateUserPayload struct {
	FirstName string   `json:"first_name" validate:"required"`
	LastName  string   `json:"last_name" validate:"required"`
	Role      UserRole `json:"role" validate:"required"`
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required"`
}

type UpdateUserPayload struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email" `
	Role      UserRole `json:"role" `
}

type ChangePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type CreateBulkUserPayload struct {
	Users []CreateUserPayload `json:"users" validate:"dive"`
}
