package entity

type ScoreRepository interface {
	GetAll() ([]Score, error)
	GetById(id string) (*Score, error)
	GetByAssignmentId(assignmentId string) ([]Score, error)
	GetByUserId(userId string) ([]Score, error)
	GetByStudentId(studentId string) ([]Score, error)
	Create(score *Score) error
	CreateMany(score []Score) error
	Update(id string, score *Score) error
	Delete(id string) error
	FilterSubmittedScoreStudents(assignmentId string, studentIds []string) ([]string, error)
}

type ScoreUseCase interface {
	GetAll() ([]Score, error)
	GetById(id string) (*Score, error)
	GetByAssignmentId(assignmentId string) (*AssignmentScore, error)
	GetByUserId(userId string) ([]Score, error)
	GetByStudentId(studentId string) ([]Score, error)
	CreateMany(userId string, assignmentId string, studentScores []StudentScore) error
	Update(user User, scoreId string, score float64) error
	Delete(user User, id string) error
	FilterSubmittedScoreStudents(assignmentId string, studentIds []string) ([]string, error)
}

type Score struct {
	Id           string  `json:"id" gorm:"primaryKey;type:char(255)"`
	Score        float64 `json:"score"`
	StudentId    string  `json:"student_id"`
	UserId       string  `json:"user_id"`
	AssignmentId string  `json:"assignment_id"`

	Email       string `json:"email" gorm:"->;-:migration"`
	FirstNameTH string `json:"first_name_th" gorm:"->;-:migration"`
	LastNameTH  string `json:"last_name_th" gorm:"->;-:migration"`
	FirstNameEN string `json:"first_name_en" gorm:"->;-:migration"`
	LastNameEN  string `json:"last_name_en" gorm:"->;-:migration"`

	Student    Student    `json:"-"`
	User       User       `json:"-"`
	Assignment Assignment `json:"-"`
}

type AssessmentClos struct {
	CLOId          string `json:"clo_id"`
	CLOCode        string `json:"clo_code"`
	CLODescription string `json:"clo_description"`
}

type AssignmentScore struct {
	CLOs                   []AssessmentClos `json:"clos"`
	IsIncludedInClo        bool             `json:"is_included_in_clo"`
	MaxScore               int              `json:"max_score"`
	ExpectedPassingStudent float64          `json:"expected_passing_student"`
	ExpectedPassingScore   float64          `json:"expected_passing_score"`
	Description            string           `json:"description"`
	Scores                 []Score          `json:"scores"`
	SubmittedAmount        int              `json:"submitted_amount"`
	EnrolledAmount         int              `json:"enrolled_amount"`
}

type CreateScoreRequestPayload struct {
	StudentId    string  `json:"student_id" validate:"required"`
	Score        float64 `json:"score" validate:"required"`
	UserId       string  `json:"user_id" validate:"required"`
	AssignmentId string  `json:"assignment_id" validate:"required"`
}

type StudentScore struct {
	StudentId string   `json:"student_id" validate:"required"`
	Score     *float64 `json:"score" validate:"required"`
}

type BulkCreateScoreRequestPayload struct {
	StudentScores []StudentScore `json:"student_scores" validate:"dive"`
	AssignmentId  string         `json:"assignment_id" validate:"required"`
}

type UpdateScoreRequestPayload struct {
	Score float64 `json:"score" validate:"required"`
}
