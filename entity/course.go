package entity

import "gorm.io/datatypes"

type CourseRepository interface {
	GetAll(query string, year string, program string) ([]Course, error)
	GetById(id string) (*Course, error)
	GetByUserId(userId string, query string, year string, program string) ([]Course, error)
	GetStudentsPassingCLOs(courseId string) (*StudentPassCLOResp, error)
	Create(course *Course) error
	Update(id string, course *Course) error
	Delete(id string) error
	CreateLinkWithLecturer(courseId string, lecturerId []string) error
	DeleteLinkWithLecturer(courseId string, lecturerId []string) error
	ReplaceLecturersForCourse(courseId string, lecturerIds []string) error
}

type CourseUseCase interface {
	GetAll(query string, year string, program string) (*GetAllCourseResponse, error)
	GetById(id string) (*Course, error)
	GetByUserId(userId string, query string, year string, program string) (*GetAllCourseResponse, error)
	GetStudentsPassingCLOs(courseId string) (*StudentPassCLOResp, error)
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
}

func (c CriteriaGrade) IsValid() bool {
	return c.A > c.BP &&
		c.BP > c.B &&
		c.B > c.CP &&
		c.CP > c.C &&
		c.C > c.DP &&
		c.DP > c.D &&
		c.D >= 0
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
	AcademicYear                 string         `json:"academic_year"`
	GraduateYear                 string         `json:"graduate_year"`
	Credit                       int            `json:"credit"`
	ExpectedPassingCloPercentage float64        `json:"expected_passing_clo_percentage"`
	IsPortfolioCompleted         bool           `json:"is_portfolio_completed" gorm:"default:false"`
	PortfolioData                datatypes.JSON `json:"portfolio_data" gorm:"type:json"`
	Result                       datatypes.JSON `json:"result" gorm:"type:json"`

	CriteriaGrade

	ProgrammeId string `json:"programme_id"`
	SemesterId  string `json:"semester_id"`

	Lecturers []*User   `gorm:"many2many:course_lecturer" json:"lecturers"`
	Semester  Semester  `json:"semester" gorm:"foreignKey:SemesterId"`
	Programme Programme `json:"programme" gorm:"foreignKey:ProgrammeId"`
}

type CreateCoursePayload struct {
	SemesterId                   string   `json:"semester_id" validate:"required"`
	LecturerIds                  []string `json:"lecturer_ids" validate:"required"`
	Name                         string   `json:"name" validate:"required"`
	Code                         string   `json:"code" validate:"required"`
	AcademicYear                 string   `json:"academic_year" validate:"required"`
	GraduateYear                 string   `json:"graduate_year" validate:"required"`
	Credit                       int      `json:"credit" validate:"required"`
	Description                  string   `json:"description" validate:"required"`
	ExpectedPassingCloPercentage float64  `json:"expected_passing_clo_percentage" validate:"required"`
	ProgrammeId                  string   `json:"programme_id" validate:"required"`
	CriteriaGrade
}

type UpdateCoursePayload struct {
	SemesterId                   string   `json:"semester_id"`
	LecturerIds                  []string `json:"lecturer_ids" `
	Name                         string   `json:"name" `
	Code                         string   `json:"code" `
	AcademicYear                 string   `json:"academic_year"`
	GraduateYear                 string   `json:"graduate_year" `
	Credit                       int      `json:"credit" `
	Description                  string   `json:"description" `
	ExpectedPassingCloPercentage float64  `json:"expected_passing_clo_percentage" `
	ProgrammeId                  string   `json:"programme_id" `
	CriteriaGrade
}

type Lecturer struct {
	Id     string `json:"id"`
	NameTH string `json:"name_th"`
	NameEN string `json:"name_en"`
}

type Program struct {
	Id     string `json:"id"`
	NameTH string `json:"name_th"`
	NameEN string `json:"name_en"`
}

type CourseSimpleData struct {
	Id           string     `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Lecturers    []Lecturer `json:"lecturers"`
	Description  string     `json:"description"`
	Credit       int        `json:"credit"`
	AcademicYear string     `json:"academic_year"`
	GraduateYear string     `json:"graduate_year"`
	Program      Program    `json:"program"`
	Semester     Semester   `json:"semester"`
}

type GetAllCourseResponse struct {
	Courses []CourseSimpleData `json:"courses"`
}

type PassOutcome struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Pass bool   `json:"pass"`
}

type StudentPassStatus struct {
	StudentID int           `json:"student_id"`
	POPass    []PassOutcome `json:"po_pass"`
	PLOPass   []PassOutcome `json:"plo_pass"`
	SOPass    []PassOutcome `json:"so_pass"`
}
