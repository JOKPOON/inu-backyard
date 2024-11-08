package entity

type GradeRepository interface {
	GetAll() ([]Grade, error)
	GetById(id string) (*Grade, error)
	GetByStudentId(studentId string) ([]Grade, error)
	FilterExisted(studentIds []string, year int, semesterSequence string) ([]string, error)
	Create(grade *Grade) error
	CreateMany(grades []Grade) error
	Update(id string, grade *Grade) error
	Delete(id string) error
}

type GradeUseCase interface {
	GetAll() ([]Grade, error)
	GetById(id string) (*Grade, error)
	GetByStudentId(studentId string) ([]Grade, error)
	FilterExisted(studentIds []string, year int, semesterSequence string) ([]string, error)
	Create(studentId string, year string, grade float64) error
	CreateMany(studentGrades []StudentGrade, year int, semesterSequence string) error
	Update(id string, grade *Grade) error
	Delete(id string) error
}

type Grade struct {
	Id         string  `json:"id" gorm:"primaryKey;type:char(255)"`
	StudentId  string  `json:"student_id"`
	SemesterId string  `json:"semester_id"`
	Grade      float64 `json:"grade"`

	Semester *Semester `json:"semester,omitempty"`
	Student  *Student  `json:"student,omitempty"`
}

type CreateGradePayload struct {
	StudentId string  `json:"student_id" validate:"required"`
	Year      string  `json:"year" validate:"required"`
	Grade     float64 `json:"grade" validate:"required"`
}

type StudentGrade struct {
	StudentId string  `json:"student_id"`
	Grade     float64 `json:"grade"`
}

type CreateManyGradesPayload struct {
	StudentGrade     []StudentGrade `json:"student_grade" validate:"dive"`
	Year             int            `json:"year" validate:"required"`
	SemesterSequence string         `json:"semester_sequence" validate:"required"`
}

type UpdateGradePayload struct {
	StudentId string  `json:"student_id"`
	Grade     float64 `json:"grade"`
}
