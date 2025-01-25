package entity

import "gorm.io/datatypes"

type ProgrammeRepository interface {
	GetAll() ([]Programme, error)
	GetByName(name string) ([]Programme, error)
	GetByNameAndYear(name string, year string) (*Programme, error)
	GetById(id string) (*Programme, error)
	Create(programme *Programme) error
	Update(name string, programme *Programme) error
	Delete(name string) error
	FilterExisted(names []string) ([]string, error)
}

type ProgrammeUseCase interface {
	GetAll() ([]Programme, error)
	GetByName(name string) ([]Programme, error)
	GetByNameAndYear(name string, year string) (*Programme, error)
	GetById(id string) (*Programme, error)
	Create(payload CreateProgrammePayload) error
	Update(name string, programme *UpdateProgrammePayload) error
	Delete(name string) error
	FilterNonExisted(names []string) ([]string, error)
}

type Category struct {
	Name   string `json:"name" validate:"required"`
	Credit int    `json:"credit" validate:"required"`
	SubCat []struct {
		Name   string `json:"name"`
		Credit int    `json:"credit"`
	} `json:"sub"`
}

type ProgrammeStructure struct {
	Category     []Category `json:"category" validate:"required"`
	TotalsCredit int        `json:"totals_credit" validate:"required"`
}

type Programme struct {
	Id   string `json:"id" gorm:"primaryKey;type:char(255)"`
	Name string `json:"name" gorm:"unique;not null"`
	Year string `json:"year" gorm:"not null"`

	Structure datatypes.JSON `json:"structure" gorm:"type:json"`
}

type CreateProgrammePayload struct {
	Name string `json:"name" validate:"required"`
	Year string `json:"year" validate:"required"`

	Structure ProgrammeStructure `json:"structure" validate:"required"`
}

type UpdateProgrammePayload struct {
	Name string `json:"name" validate:"required"`
}
