package entity

type StudentRepository interface {
	GetById(id string) (*Student, error)
	GetAll() ([]Student, error)
	GetByParams(params *Student, limit int, offset int) ([]Student, error)
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
	GetByParams(params *Student, limit int, offset int) ([]Student, error)
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
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          string  `json:"email"`
	ProgrammeName  string  `json:"programme_name"`
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

type CreateStudentPayload struct {
	Id        string   `json:"id" validate:"required"`
	FirstName string   `json:"first_name" validate:"required"`
	LastName  string   `json:"last_name" validate:"required"`
	GPAX      *float64 `json:"gpax" validate:"required"`
	MathGPA   *float64 `json:"math_gpa" validate:"required"`
	EngGPA    *float64 `json:"eng_gpa" validate:"required"`
	SciGPA    *float64 `json:"sci_gpa" validate:"required"`
	School    string   `json:"school" validate:"required"`
	City      string   `json:"city" validate:"required"`
	Email     string   `json:"email" validate:"required"`
	Year      string   `json:"year" validate:"required"`
	Admission string   `json:"admission" validate:"required"`
	Remark    string   `json:"remark"`

	ProgrammeName  string `json:"programme_name" validate:"required"`
	DepartmentName string `json:"department_name" validate:"required"`
}

type GetStudentsByParamsPayload struct {
	Year           string `json:"year"`
	ProgrammeName  string `json:"programme_name"`
	DepartmentName string `json:"department_name"`
}

type CreateBulkStudentsPayload struct {
	Students []CreateStudentPayload `json:"students" validate:"dive"`
}

type UpdateStudentPayload struct {
	Id        string   `json:"id" validate:"required"`
	FirstName string   `json:"first_name" validate:"required"`
	LastName  string   `json:"last_name" validate:"required"`
	GPAX      *float64 `json:"gpax" validate:"required"`
	MathGPA   *float64 `json:"math_gpa" validate:"required"`
	EngGPA    *float64 `json:"eng_gpa" validate:"required"`
	SciGPA    *float64 `json:"sci_gpa" validate:"required"`
	School    string   `json:"school" validate:"required"`
	City      string   `json:"city" validate:"required"`
	Email     string   `json:"email" validate:"required"`
	Year      string   `json:"year" validate:"required"`
	Admission string   `json:"admission" validate:"required"`
	Remark    *string  `json:"remark" validate:"required"`

	ProgrammeName  string `json:"programme_name" validate:"required"`
	DepartmentName string `json:"department_name" validate:"required"`
}
