package entity

import "gorm.io/datatypes"

type CriteriaGrade struct {
	A  float64 `json:"criteriaGradeA" gorm:"column:criteria_grade_a" validate:"required"`
	BP float64 `json:"criteriaGradeBP" gorm:"column:criteria_grade_bp" validate:"required"`
	B  float64 `json:"criteriaGradeB" gorm:"column:criteria_grade_b" validate:"required"`
	CP float64 `json:"criteriaGradeCP" gorm:"column:criteria_grade_cp" validate:"required"`
	C  float64 `json:"criteriaGradeC" gorm:"column:criteria_grade_c" validate:"required"`
	DP float64 `json:"criteriaGradeDP" gorm:"column:criteria_grade_dp" validate:"required"`
	D  float64 `json:"criteriaGradeD" gorm:"column:criteria_grade_d" validate:"required"`
	F  float64 `json:"criteriaGradeF" gorm:"column:criteria_grade_f" validate:"required"`
}

func (c CriteriaGrade) IsValid() bool {
	return c.A >= c.BP &&
		c.BP >= c.B &&
		c.B >= c.CP &&
		c.CP >= c.C &&
		c.C >= c.DP &&
		c.DP >= c.D &&
		c.D >= c.F &&
		c.F >= 0
}

func (c CriteriaGrade) CalculateCriteriaWeight(maxScore float64) CriteriaGrade {
	percentage := maxScore / 100

	criteriaGrade := CriteriaGrade{
		A:  c.A * percentage,
		BP: c.BP * percentage,
		B:  c.B * percentage,
		CP: c.CP * percentage,
		C:  c.C * percentage,
		DP: c.DP * percentage,
		D:  c.D * percentage,
		F:  c.F * percentage,
	}
	return criteriaGrade
}

func (c CriteriaGrade) GradeToGPA(grade string) float64 {
	switch grade {
	case "A":
		return 4.0
	case "BP":
		return 3.5
	case "B":
		return 3.0
	case "CP":
		return 2.5
	case "C":
		return 2.0
	case "DP":
		return 1.5
	case "D":
		return 1.0
	case "F":
		return 0
	default:
		return 0
	}
}

type Course struct {
	Id                           string  `json:"id" gorm:"primaryKey;type:char(255)"`
	Name                         string  `json:"name"`
	Code                         string  `json:"code"`
	Curriculum                   string  `json:"curriculum"`
	Description                  string  `json:"description"`
	ExpectedPassingCloPercentage float64 `json:"expectedPassingCloPercentage"`
	IsPortfolioCompleted         *bool   `json:"isPortfolioCompleted" gorm:"default:false"`
	PortfolioData                datatypes.JSON
	AcademicYear                 int `json:"academicYear"`
	GraduateYear                 int `json:"graduateYear"`
	ProgramYear                  int `json:"programYear"`

	SemesterId string `json:"semesterId"`
	UserId     string `json:"userId"`
	CriteriaGrade

	Semester Semester `json:"semester"`
	User     User     `json:"user"`
}

type CreateCoursePayload struct {
	SemesterId                   string  `json:"semesterId" validate:"required"`
	UserId                       string  `json:"userId" validate:"required"`
	Name                         string  `json:"name" validate:"required"`
	Code                         string  `json:"code" validate:"required"`
	Curriculum                   string  `json:"curriculum" validate:"required"`
	Description                  string  `json:"description" validate:"required"`
	ExpectedPassingCloPercentage float64 `json:"expectedPassingCloPercentage" validate:"required"`
	AcademicYear                 int     `json:"academicYear" validate:"required"`
	GraduateYear                 int     `json:"graduateYear" validate:"required"`
	ProgramYear                  int     `json:"programYear" validate:"required"`
	CriteriaGrade
}

type UpdateCoursePayload struct {
	Name                         string  `json:"name" validate:"required"`
	Code                         string  `json:"code" validate:"required"`
	Curriculum                   string  `json:"curriculum" validate:"required"`
	Description                  string  `json:"description" validate:"required"`
	ExpectedPassingCloPercentage float64 `json:"expectedPassingCloPercentage" validate:"required"`
	AcademicYear                 int     `json:"academicYear" validate:"required"`
	GraduateYear                 int     `json:"graduateYear" validate:"required"`
	ProgramYear                  int     `json:"programYear" validate:"required"`
	IsPortfolioCompleted         *bool   `json:"isPortfolioCompleted" validate:"required"`
	CriteriaGrade
}

type CourseRepository interface {
	GetAll() ([]Course, error)
	GetById(id string) (*Course, error)
	GetByUserId(userId string) ([]Course, error)
	Create(course *Course) error
	Update(id string, course *Course) error
	Delete(id string) error
}
type CourseUseCase interface {
	GetAll() ([]Course, error)
	GetById(id string) (*Course, error)
	GetByUserId(userId string) ([]Course, error)
	Create(user User, payload CreateCoursePayload) error
	Update(user User, id string, name string, code string, curriculum string, description string, expectedPassingCloPercentage float64, academicYear int, graduateYear int, programYear int, criteriaGrade CriteriaGrade, isPortfolioCompleted bool) error
	Delete(user User, id string) error
}
