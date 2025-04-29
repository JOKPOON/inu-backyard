package entity

type FacultyRepository interface {
	GetAll() ([]Faculty, error)
	GetById(id string) (*Faculty, error)
	Create(faculty *Faculty) error
	Update(faculty *Faculty) error
	Delete(name string) error
}

type FacultyUseCase interface {
	GetAll() ([]Faculty, error)
	GetById(id string) (*Faculty, error)
	Create(faculty *Faculty) error
	Update(faculty *Faculty) error
	Delete(name string) error
}

type Faculty struct {
	Id          string       `json:"id" gorm:"primaryKey;type:char(255)"`
	NameTH      string       `json:"name_th" gorm:"type:char(255)"`
	NameEN      string       `json:"name_en" gorm:"type:char(255)"`
	Departments []Department `gorm:"foreignKey:FacultyId" json:"departments"`
}

type CreateFacultyRequestPayload struct {
	NameTH string `json:"name_th" validate:"required"`
	NameEN string `json:"name_en" validate:"required"`
}

type UpdateFacultyRequestPayload struct {
	NameTH string `json:"name_th" validate:"required"`
	NameEN string `json:"name_en" validate:"required"`
}
