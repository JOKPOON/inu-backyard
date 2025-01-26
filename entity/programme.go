package entity

import "gorm.io/datatypes"

type ProgrammeRepository interface {
	GetAll() ([]Programme, error)
	GetByName(nameTH string, nameEN string) ([]Programme, error)
	GetByNameAndYear(nameTH string, nameEN string, year string) (*Programme, error)
	GetById(id string) (*Programme, error)
	Create(programme *Programme) error
	Update(name string, programme *Programme) error
	Delete(name string) error
	FilterExisted(nameTH []string, nameEN []string) ([]string, error)
}

type ProgrammeUseCase interface {
	GetAll() ([]Programme, error)
	GetByName(nameTH string, nameEN string) ([]Programme, error)
	GetByNameAndYear(nameTH string, nameEN string, year string) (*Programme, error)
	GetById(id string) (*Programme, error)
	Create(payload CreateProgrammePayload) error
	Update(name string, programme *UpdateProgrammePayload) error
	Delete(name string) error
	FilterNonExisted(namesTH []string, namesEN []string) ([]string, error)
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
	Id            string `json:"id" gorm:"primaryKey;type:char(255)"`
	NameTH        string `json:"name_th" gorm:"not null"`
	NameEN        string `json:"name_en" gorm:"not null"`
	DegreeTH      string `json:"degree_th" gorm:"not null"`
	DegreeEN      string `json:"degree_en" gorm:"not null"`
	DegreeShortTH string `json:"degree_short_th" gorm:"not null"`
	DegreeShortEN string `json:"degree_short_en" gorm:"not null"`
	Year          string `json:"year" gorm:"not null"`

	Structure datatypes.JSON `json:"structure" gorm:"type:json"`
}

type CreateProgrammePayload struct {
	NameTH        string `json:"name_th" validate:"required"`
	NameEN        string `json:"name_en" validate:"required"`
	DegreeTH      string `json:"degree_th" validate:"required"`
	DegreeEN      string `json:"degree_en" validate:"required"`
	DegreeShortTH string `json:"degree_short_th" validate:"required"`
	DegreeShortEN string `json:"degree_short_en" validate:"required"`
	Year          string `json:"year" validate:"required"`

	Structure ProgrammeStructure `json:"structure" validate:"required"`
}

type UpdateProgrammePayload struct {
	NameTH        string `json:"name_th" validate:"required"`
	NameEN        string `json:"name_en" validate:"required"`
	DegreeTH      string `json:"degree_th" validate:"required"`
	DegreeEN      string `json:"degree_en" validate:"required"`
	DegreeShortTH string `json:"degree_short_th" validate:"required"`
	DegreeShortEN string `json:"degree_short_en" validate:"required"`
	Year          string `json:"year" validate:"required"`

	Structure ProgrammeStructure `json:"structure" validate:"required"`
}
