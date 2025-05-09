package entity

type StudentRepository interface {
	GetById(id string) (*Student, error)
	GetAll() ([]Student, error)
	GetByParams(query string, params *Student, limit int, offset int) (*GetStudentsResponse, error)
	Create(student *Student) error
	CreateMany(student []Student) error
	Update(id string, student *Student) error
	Delete(id string) error
	FilterExisted(studentIds []string) ([]string, error)

	GetAllSchools() ([]string, error)
	GetAllAdmissions() ([]string, error)
}

type StudentUseCase interface {
	GetById(id string) (*Student, error)
	GetAll() ([]Student, error)
	GetByParams(query string, params *Student, limit int, offset int) (*GetStudentsResponse, error)
	CreateMany(student []CreateStudentPayload) error
	Update(id string, student *UpdateStudentPayload) error
	Delete(id string) error
	FilterExisted(studentIds []string) ([]string, error)
	FilterNonExisted(studentIds []string) ([]string, error)

	GetAllSchools() ([]string, error)
	GetAllAdmissions() ([]string, error)
}

type Student struct {
	Id          string `gorm:"primaryKey;type:char(255)" json:"id"`
	FirstNameTH string `json:"first_name_th"`
	LastNameTH  string `json:"last_name_th"`
	FirstNameEN string `json:"first_name_en"`
	LastNameEN  string `json:"last_name_en"`
	Email       string `json:"email"`
	Year        string `json:"year"`

	ProgrammeId string    `json:"programme_id"`
	Programme   Programme `gorm:"foreignKey:ProgrammeId" json:"programme"`

	DepartmentId string     `json:"department_id"`
	Department   Department `gorm:"foreignKey:DepartmentId" json:"department"`
}

type CreateBulkStudentsPayload struct {
	Students []CreateStudentPayload `json:"students" validate:"dive"`
}

type CreateStudentPayload struct {
	Id          string `json:"id" validate:"required"`
	FirstNameTH string `json:"first_name_th"`
	LastNameTH  string `json:"last_name_th"`
	FirstNameEN string `json:"first_name_en" validate:"required"`
	LastNameEN  string `json:"last_name_en" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Year        string `json:"year" validate:"required"`

	ProgrammeId  string `json:"programme_id" validate:"required"`
	DepartmentId string `json:"department_id" validate:"required"`
}

type GetStudentsByParamsPayload struct {
	Year         string `json:"year"`
	ProgrammeId  string `json:"programme_id"`
	DepartmentId string `json:"department_id"`
}

type UpdateStudentPayload struct {
	Id          string `json:"id" validate:"required"`
	FirstNameTH string `json:"first_name_th"`
	LastNameTH  string `json:"last_name_th" `
	FirstNameEN string `json:"first_name_en"`
	LastNameEN  string `json:"last_name_en"`
	Email       string `json:"email" validate:"required"`
	Year        string `json:"year" validate:"required"`

	ProgrammeId  string `json:"programme_id" validate:"required"`
	DepartmentId string `json:"department_id" validate:"required"`
}

type GetStudentsResponse struct {
	Students  []Student `json:"students"`
	Total     int       `json:"total"`
	Page      int       `json:"page"`
	TotalPage int       `json:"total_page"`
}
