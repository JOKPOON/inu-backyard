package entity

import "github.com/team-inu/inu-backyard/entity"

type CreateUserPayload struct {
	FirstName string          `json:"firstName" validate:"required"`
	LastName  string          `json:"lastName" validate:"required"`
	Role      entity.UserRole `json:"role" validate:"required"`
	Email     string          `json:"email" validate:"required,email"`
	Password  string          `json:"password" validate:"required"`
}

type UpdateUserPayload struct {
	FirstName string          `json:"firstName" validate:"required"`
	LastName  string          `json:"lastName" validate:"required"`
	Email     string          `json:"email" validate:"required"`
	Role      entity.UserRole `json:"role" validate:"required"`
}

type ChangePasswordPayload struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type CreateBulkUserPayload struct {
	Users []CreateUserPayload `json:"users" validate:"dive"`
}
