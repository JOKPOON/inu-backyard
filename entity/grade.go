package entity

type Grade struct {
	Id         string  `json:"id" gorm:"primaryKey;type:char(255)"`
	StudentId  string  `json:"studentId"`
	SemesterId string  `json:"semesterId"`
	Grade      float64 `json:"grade"`

	Semester *Semester `json:"semester,omitempty"`
	Student  *Student  `json:"student,omitempty"`
}

type StudentGrade struct {
	StudentId string
	Grade     float64
}

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
