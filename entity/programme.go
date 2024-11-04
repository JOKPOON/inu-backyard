package entity

type Programme struct {
	Name string `json:"name" gorm:"primaryKey;type:char(255)"`
}

type ProgrammeRepository interface {
	GetAll() ([]Programme, error)
	Get(name string) (*Programme, error)
	Create(programme *Programme) error
	Update(name string, programme *Programme) error
	Delete(name string) error
	FilterExisted(names []string) ([]string, error)
}

type ProgrammeUseCase interface {
	GetAll() ([]Programme, error)
	Get(name string) (*Programme, error)
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
