package entity

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserRole string

const (
	UserRoleLecturer         UserRole = "LECTURER"
	UserRoleModerator        UserRole = "MODERATOR"
	UserRoleHeadOfCurriculum UserRole = "HEAD_OF_CURRICULUM"
	UserRoleAUNQAManager     UserRole = "AUN-QA_MANAGER"
	UserRoleTABEEManager     UserRole = "TABEE_MANAGER"
	UserRoleABETManager      UserRole = "ABET_MANAGER"
)

var Roles = []UserRole{
	UserRoleLecturer,
	UserRoleModerator,
	UserRoleHeadOfCurriculum,
	UserRoleAUNQAManager,
	UserRoleTABEEManager,
	UserRoleABETManager,
}

func (u User) IsRoles(expectedRoles []UserRole) bool {
	for _, expectedRole := range expectedRoles {
		role := strings.Split(string(u.Role), ",")
		for _, r := range role {
			if r == string(expectedRole) {
				return true
			}
		}
	}

	return false
}

type UserRepository interface {
	GetAll(query string, offset int, limit int) (*Pagination, error)
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
	GetAll(query string, index string, size string) (*Pagination, error)
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
	Id                 string   `json:"id" gorm:"primaryKey;type:char(255)"`
	Email              string   `json:"email" gorm:"unique"`
	Password           string   `json:"password"`
	FirstNameTH        string   `json:"first_name_th"`
	LastNameTH         string   `json:"last_name_th"`
	FirstNameEN        string   `json:"first_name_en"`
	LastNameEN         string   `json:"last_name_en"`
	AcademicPositionTH string   `json:"academic_position_th"`
	AcademicPositionEN string   `json:"academic_position_en"`
	Role               UserRole `json:"role" gorm:"default:'LECTURER'"`
	DegreeTH           string   `json:"degree_th"`
	DegreeEN           string   `json:"degree_en"`
}

type CreateUserPayload struct {
	FirstNameTH        string   `json:"first_name_th" validate:"required"`
	LastNameTH         string   `json:"last_name_th" validate:"required"`
	FirstNameEN        string   `json:"first_name_en" validate:"required"`
	LastNameEN         string   `json:"last_name_en" validate:"required"`
	Role               UserRole `json:"role" validate:"required"`
	AcademicPositionTH string   `json:"academic_position_th"`
	AcademicPositionEN string   `json:"academic_position_en"`
	Email              string   `json:"email" validate:"required,email"`
	Password           string   `json:"password"`
	DegreeTH           string   `json:"degree_th" validate:"required"`
	DegreeEN           string   `json:"degree_en" validate:"required"`
}

type UpdateUserPayload struct {
	FirstNameTH        string   `json:"first_name_th"`
	LastNameTH         string   `json:"last_name_th"`
	FirstNameEN        string   `json:"first_name_en"`
	LastNameEN         string   `json:"last_name_en"`
	Email              string   `json:"email"`
	AcademicPositionTH string   `json:"academic_position_th"`
	AcademicPositionEN string   `json:"academic_position_en"`
	Role               UserRole `json:"role" `
	DegreeTH           string   `json:"degree_th"`
	DegreeEN           string   `json:"degree_en"`
}

type ChangePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type ForgotPasswordPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordPayload struct {
	Email       string `json:"email" validate:"required,email"`
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type CreateBulkUserPayload struct {
	Users []CreateUserPayload `json:"users" validate:"dive"`
}

type Pagination struct {
	Size      int         `json:"size"`
	Page      int         `json:"page"`
	TotalPage int         `json:"total_page"`
	Total     int64       `json:"total"`
	Data      interface{} `json:"data"`
}
