package entity

type EnrollmentStatus string

const (
	EnrollmentStatusEnroll   EnrollmentStatus = "ENROLL"
	EnrollmentStatusWithdraw EnrollmentStatus = "WITHDRAW"
)

type EnrollmentRepository interface {
	GetAll() ([]Enrollment, error)
	GetById(id string) (*Enrollment, error)
	GetByCourseId(courseId string, query string) ([]Enrollment, error)
	GetByStudentId(studentId string) ([]Enrollment, error)
	Create(enrollment *Enrollment) error
	CreateMany(enrollments []Enrollment) error
	Update(id string, enrollment *Enrollment) error
	Delete(id string) error
	FilterExisted(ids []string) ([]string, error)
	FilterJoinedStudent(studentIds []string, courseId string, withStatus *EnrollmentStatus) ([]string, error)
}

type EnrollmentUseCase interface {
	GetAll() ([]Enrollment, error)
	GetById(id string) (*Enrollment, error)
	GetByCourseId(courseId string, query string) ([]Enrollment, error)
	GetByStudentId(studentId string) ([]Enrollment, error)
	CreateMany(CreateEnrollmentsPayload) error
	Update(id string, status EnrollmentStatus) error
	Delete(id string) error
	FilterJoinedStudent(studentIds []string, courseId string, withStatus *EnrollmentStatus) ([]string, error)
}

type Enrollment struct {
	Id          string           `json:"id" gorm:"primaryKey;type:char(255)"`
	CourseId    string           `json:"course_id"`
	StudentId   string           `json:"student_id"`
	Status      EnrollmentStatus `json:"status" gorm:"type:enum('ENROLL','WITHDRAW')"`
	Email       string           `json:"email" gorm:"->;-:migration"`
	FirstNameTH string           `json:"first_name_th" gorm:"->;-:migration"`
	LastNameTH  string           `json:"last_name_th" gorm:"->;-:migration"`
	FirstNameEN string           `json:"first_name_en" gorm:"->;-:migration"`
	LastNameEN  string           `json:"last_name_en" gorm:"->;-:migration"`

	Course  Course  `json:"-"`
	Student Student `json:"-"`
}

type CreateEnrollmentsPayload struct {
	CourseId   string           `json:"course_id" validate:"required"`
	StudentIds []string         `json:"student_ids" validate:"required"`
	Status     EnrollmentStatus `json:"status" validate:"required"`
}

type UpdateEnrollmentPayload struct {
	Status EnrollmentStatus `json:"status" validate:"required"`
}
