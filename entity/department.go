package entity

type DepartmentRepository interface {
	GetAll() ([]Department, error)
	GetByName(id string) (*Department, error)
	Create(department *Department) error
	Update(department *Department, newName string) error
	Delete(name string) error
	FilterExisted(names []string) ([]string, error)
}

type DepartmentUseCase interface {
	GetAll() ([]Department, error)
	GetByName(name string) (*Department, error)
	Create(department *Department) error
	Update(department *Department, newName string) error
	Delete(id string) error
	FilterNonExisted(names []string) ([]string, error)
}

type Department struct {
	Name        string `json:"name" gorm:"type:char(255);unique;not null;primaryKey"`
	FacultyName string `json:"faculty_name"`

	Faculty    Faculty     `gorm:"foreignKey:FacultyName" json:"-"`
	Programmes []Programme `gorm:"foreignKey:DepartmentName" json:"programmes"`
}

type CreateDepartmentRequestPayload struct {
	Name        string `json:"name" validate:"required"`
	FacultyName string `json:"faculty_name" validate:"required"`
}

type UpdateDepartmentRequestPayload struct {
	NewName     string `json:"new_name" validate:"required"`
	FacultyName string `json:"faculty_name" validate:"required"`
}

type DeleteDepartmentRequestPayload struct {
	Name string `json:"name" validate:"required"`
}
