package entity

type StudentRepository interface {
	GetById(id string) (*Student, error)
	GetAll() ([]Student, error)
	GetByParams(query string, params *Student, limit int, offset int) ([]Student, error)
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
	GetByParams(query string, params *Student, limit int, offset int) ([]Student, error)
	CreateMany(student []CreateStudentPayload) error
	Update(id string, student *UpdateStudentPayload) error
	Delete(id string) error
	FilterExisted(studentIds []string) ([]string, error)
	FilterNonExisted(studentIds []string) ([]string, error)

	GetAllSchools() ([]string, error)
	GetAllAdmissions() ([]string, error)
}

type Student struct {
	Id             string  `gorm:"primaryKey;type:char(255)" json:"id"`
	FirstNameTH    string  `json:"first_name_th"`
	LastNameTH     string  `json:"last_name_th"`
	FirstNameEN    string  `json:"first_name_en"`
	LastNameEN     string  `json:"last_name_en"`
	Email          string  `json:"email"`
	ProgrammeId    string  `json:"programme_id"`
	DepartmentName string  `json:"department_name"`
	GPAX           float64 `json:"gpax"`
	MathGPA        float64 `json:"math_gpa"`
	EngGPA         float64 `json:"eng_gpa"`
	SciGPA         float64 `json:"sci_gpa"`
	School         string  `json:"school"`
	City           string  `json:"city"`
	Year           string  `json:"year"`
	Admission      string  `json:"admission"`
	Remark         string  `json:"remark"`

	Programme  *Programme  `json:"programme,omitempty"`
	Department *Department `json:"department,omitempty"`
}

type CreateBulkStudentsPayload struct {
	Students []CreateStudentPayload `json:"students" validate:"dive"`
}

type CreateStudentPayload struct {
	Id          string   `json:"id" validate:"required"`
	FirstNameTH string   `json:"first_name_th" validate:"required"`
	LastNameTH  string   `json:"last_name_th" validate:"required"`
	FirstNameEN string   `json:"first_name_en" validate:"required"`
	LastNameEN  string   `json:"last_name_en" validate:"required"`
	GPAX        *float64 `json:"gpax" validate:"required"`
	MathGPA     *float64 `json:"math_gpa" validate:"required"`
	EngGPA      *float64 `json:"eng_gpa" validate:"required"`
	SciGPA      *float64 `json:"sci_gpa" validate:"required"`
	School      string   `json:"school" validate:"required"`
	City        string   `json:"city" validate:"required"`
	Email       string   `json:"email" validate:"required"`
	Year        string   `json:"year" validate:"required"`
	Admission   string   `json:"admission" validate:"required"`
	Remark      string   `json:"remark"`

	ProgrammeId    string `json:"programme_id" validate:"required"`
	DepartmentName string `json:"department_name" validate:"required"`
}

type GetStudentsByParamsPayload struct {
	Year           string `json:"year"`
	ProgrammeId    string `json:"programme_id"`
	DepartmentName string `json:"department_name"`
}

type UpdateStudentPayload struct {
	Id          string   `json:"id" validate:"required"`
	FirstNameTH string   `json:"first_name_th" validate:"required"`
	LastNameTH  string   `json:"last_name_th" validate:"required"`
	FirstNameEN string   `json:"first_name_en" validate:"required"`
	LastNameEN  string   `json:"last_name_en" validate:"required"`
	GPAX        *float64 `json:"gpax" validate:"required"`
	MathGPA     *float64 `json:"math_gpa" validate:"required"`
	EngGPA      *float64 `json:"eng_gpa" validate:"required"`
	SciGPA      *float64 `json:"sci_gpa" validate:"required"`
	School      string   `json:"school" validate:"required"`
	City        string   `json:"city" validate:"required"`
	Email       string   `json:"email" validate:"required"`
	Year        string   `json:"year" validate:"required"`
	Admission   string   `json:"admission" validate:"required"`
	Remark      *string  `json:"remark" validate:"required"`

	ProgrammeId    string `json:"programme_id" validate:"required"`
	DepartmentName string `json:"department_name" validate:"required"`
}
