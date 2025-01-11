package entity

type Programme struct {
	Id   string `json:"id" gorm:"primaryKey;type:char(255)"`
	Name string `json:"name" gorm:"unique;not null"`
}

type ProgrammeRepository interface {
	GetAll() ([]Programme, error)
	GetByName(name string) (*Programme, error)
	GetById(id string) (*Programme, error)
	Create(programme *Programme) error
	Update(name string, programme *Programme) error
	Delete(name string) error
	FilterExisted(names []string) ([]string, error)
}

type ProgrammeUseCase interface {
	GetAll() ([]Programme, error)
	GetByName(name string) (*Programme, error)
	GetById(id string) (*Programme, error)
	Create(name string) error
	Update(name string, programme *UpdateProgrammePayload) error
	Delete(name string) error
	FilterNonExisted(names []string) ([]string, error)
}

type CreateProgrammePayload struct {
	Name string `json:"name" validate:"required"`
}

type UpdateProgrammePayload struct {
	Name string `json:"name" validate:"required"`
}
