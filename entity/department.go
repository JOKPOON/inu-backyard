package entity

type DepartmentRepository interface {
	GetAll() ([]Department, error)
	GetById(id string) (*Department, error)
	Create(department *Department) error
	Update(department *Department) error
	Delete(name string) error
	FilterExisted(names []string) ([]string, error)
}

type DepartmentUseCase interface {
	GetAll() ([]Department, error)
	GetById(id string) (*Department, error)
	Create(department *Department) error
	Update(department *Department) error
	Delete(id string) error
	FilterNonExisted(names []string) ([]string, error)
}

type Department struct {
	Id        string `json:"id" gorm:"primaryKey;type:char(255)"`
	NameTH    string `json:"name_th" gorm:"type:char(255)"`
	NameEN    string `json:"name_en" gorm:"type:char(255)"`
	FacultyId string `json:"faculty_id" gorm:"type:char(255)"`

	Faculty    Faculty     `gorm:"foreignKey:FacultyId" json:"-"`
	Programmes []Programme `gorm:"foreignKey:DepartmentId" json:"programmes"`
}

type CreateDepartmentRequestPayload struct {
	NameTH    string `json:"name_th" gorm:"type:char(255)"`
	NameEN    string `json:"name_en" gorm:"type:char(255)"`
	FacultyId string `json:"faculty_id" gorm:"type:char(255)"`
}

type UpdateDepartmentRequestPayload struct {
	NameTH    string `json:"name_th" gorm:"type:char(255)"`
	NameEN    string `json:"name_en" gorm:"type:char(255)"`
	FacultyId string `json:"faculty_id" gorm:"type:char(255)"`
}
