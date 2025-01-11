package entity

import "gorm.io/datatypes"

type CourseRepository interface {
	GetAll() ([]Course, error)
	GetById(id string) (*Course, error)
	GetByUserId(userId string) ([]Course, error)
	Create(course *Course) error
	Update(id string, course *Course) error
	Delete(id string) error
	CreateLinkWithLecturer(courseId string, lecturerId []string) error
	DeleteLinkWithLecturer(courseId string, lecturerId []string) error
}

type CourseUseCase interface {
	GetAll() ([]Course, error)
	GetById(id string) (*Course, error)
	GetByUserId(userId string) ([]Course, error)
	Create(user User, payload CreateCoursePayload) error
	Update(user User, id string, payload UpdateCoursePayload) error
	Delete(user User, id string) error
}

type CriteriaGrade struct {
	A  float64 `json:"criteria_grade_a" gorm:"column:criteria_grade_a" validate:"required"`
	BP float64 `json:"criteria_grade_bp" gorm:"column:criteria_grade_bp" validate:"required"`
	B  float64 `json:"criteria_grade_b" gorm:"column:criteria_grade_b" validate:"required"`
	CP float64 `json:"criteria_grade_cp" gorm:"column:criteria_grade_cp" validate:"required"`
	C  float64 `json:"criteria_grade_c" gorm:"column:criteria_grade_c" validate:"required"`
	DP float64 `json:"criteria_grade_dp" gorm:"column:criteria_grade_dp" validate:"required"`
	D  float64 `json:"criteria_grade_d" gorm:"column:criteria_grade_d" validate:"required"`
	F  float64 `json:"criteria_grade_f" gorm:"column:criteria_grade_f" validate:"required"`
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
	Id                           string         `json:"id" gorm:"primaryKey;type:char(255)"`
	Name                         string         `json:"name"`
	Code                         string         `json:"code"`
	Description                  string         `json:"description"`
	Credit                       int            `json:"credit"`
	ExpectedPassingCloPercentage float64        `json:"expected_passing_clo_percentage"`
	IsPortfolioCompleted         bool           `json:"is_portfolio_completed" gorm:"default:false"`
	PortfolioData                datatypes.JSON `json:"portfolio_data" gorm:"type:json"`

	CriteriaGrade

	ProgrammeId string `json:"programme_id"`
	SemesterId  string `json:"semester_id"`

	Lecturers []*User   `gorm:"many2many:course_lecturer" json:"lecturers"`
	Semester  Semester  `json:"semester"`
	Programme Programme `json:"programme"`
}

type CreateCoursePayload struct {
	SemesterId                   string   `json:"semester_id" validate:"required"`
	LecturerIds                  []string `json:"lecturer_ids" validate:"required"`
	Name                         string   `json:"name" validate:"required"`
	Code                         string   `json:"code" validate:"required"`
	Description                  string   `json:"description" validate:"required"`
	ExpectedPassingCloPercentage float64  `json:"expected_passing_clo_percentage" validate:"required"`
	ProgrammeId                  string   `json:"programme_id" validate:"required"`
	CriteriaGrade
}

type UpdateCoursePayload struct {
	SemesterId                   string   `json:"semester_id" validate:"required"`
	LecturerIds                  []string `json:"lecturer_ids" validate:"required"`
	Name                         string   `json:"name" validate:"required"`
	Code                         string   `json:"code" validate:"required"`
	Description                  string   `json:"description" validate:"required"`
	ExpectedPassingCloPercentage float64  `json:"expected_passing_clo_percentage" validate:"required"`
	ProgrammeId                  string   `json:"programme_id" validate:"required"`
	CriteriaGrade
}
