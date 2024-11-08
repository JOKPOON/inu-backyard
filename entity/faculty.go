package entity

type FacultyRepository interface {
	GetAll() ([]Faculty, error)
	GetByName(id string) (*Faculty, error)
	Create(faculty *Faculty) error
	Update(faculty *Faculty, newName string) error
	Delete(name string) error
}

type FacultyUseCase interface {
	GetAll() ([]Faculty, error)
	GetByName(name string) (*Faculty, error)
	Create(faculty *Faculty) error
	Update(faculty *Faculty, newName string) error
	Delete(name string) error
}

type Faculty struct {
	Name string `json:"name" gorm:"primaryKey;type:char(255)"`
}

type CreateFacultyRequestPayload struct {
	Name string `json:"name" validate:"required"`
}

type UpdateFacultyRequestPayload struct {
	NewName string `json:"new_name" validate:"required"`
}

type DeleteFacultyRequestPayload struct {
	Name string `json:"name" validate:"required"`
}
