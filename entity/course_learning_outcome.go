package entity

type CourseLearningOutcomeRepository interface {
	GetAll() ([]CourseLearningOutcome, error)
	GetById(id string) (*CourseLearningOutcome, error)
	GetByCourseId(courseId string) ([]CourseLearningOutcomeWithPO, error)
	Create(courseLearningOutcome *CourseLearningOutcome) error
	CreateLinkProgramOutcome(id string, programOutcomeIds []string) error
	CreateLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeIds []string) error
	CreateLinkSubStudentOutcome(id string, subStudentOutcomeIds []string) error
	Update(id string, courseLearningOutcome *CourseLearningOutcome) error
	Delete(id string) error
	DeleteLinkProgramOutcome(id string, programOutcomeId string) error
	DeleteLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeId string) error
	DeleteLinkSubStudentOutcome(id string, subStudentOutcomeId string) error
	FilterExisted(ids []string) ([]string, error)
}

type CourseLearningOutcomeUseCase interface {
	GetAll() ([]CourseLearningOutcome, error)
	GetById(id string) (*CourseLearningOutcome, error)
	GetByCourseId(courseId string) ([]CourseLearningOutcomeWithPO, error)
	Create(dto CreateCourseLearningOutcomePayload) error
	CreateLinkProgramOutcome(id string, programOutcomeIds []string) error
	CreateLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeIds []string) error
	CreateLinkSubStudentOutcome(id string, subStudentOutcomeIds []string) error
	Update(id string, dto UpdateCourseLearningOutcomePayload) error
	Delete(id string) error
	DeleteLinkProgramOutcome(id string, programOutcomeId string) error
	DeleteLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeId string) error
	DeleteLinkSubStudentOutcome(id string, subStudentOutcomeId string) error
	FilterNonExisted(ids []string) ([]string, error)
}

type CourseLearningOutcome struct {
	Id                                  string  `json:"id" gorm:"primaryKey;type:char(255)"`
	Code                                string  `json:"code"`
	Description                         string  `json:"description"`
	ExpectedPassingAssignmentPercentage float64 `json:"expected_passing_assignment_percentage"`
	ExpectedPassingStudentPercentage    float64 `json:"expected_passing_student_percentage"`
	Status                              string  `json:"status"`
	CourseId                            string  `json:"course_id"`

	ProgramOutcomes            []*ProgramOutcome            `gorm:"many2many:clo_po" json:"program_outcomes"`
	SubProgramLearningOutcomes []*SubProgramLearningOutcome `gorm:"many2many:clo_subplo" json:"sub_program_learning_outcomes"`
	SubStudentOutcomes         []*SubStudentOutcome         `gorm:"many2many:clo_subso" json:"sub_student_outcomes"`
	Assignments                []*Assignment                `gorm:"many2many:clo_assignment" json:"assignments"`
	Course                     Course                       `json:"course"`
}

type CourseLearningOutcomeDal struct {
	Code                                string  `json:"code"`
	Description                         string  `json:"description"`
	ExpectedPassingAssignmentPercentage float64 `json:"expected_passing_assignment_percentage"`
	ExpectedPassingStudentPercentage    float64 `json:"expected_passing_student_percentage"`
	Status                              string  `json:"status"`
	CourseId                            string  `json:"course_id"`
}

type CourseLearningOutcomeWithOutcomeDal struct {
	CourseLearningOutcomeDal
	ProgramOutcomes            []*ProgramOutcome            `gorm:"many2many:clo_po" json:"program_outcomes"`
	SubProgramLearningOutcomes []*SubProgramLearningOutcome `gorm:"many2many:clo_subplo" json:"sub_program_learning_outcomes"`
	SubStudentOutcomes         []*SubStudentOutcome         `gorm:"many2many:clo_subso" json:"sub_student_outcomes"`
}

type CourseLearningOutcomeWithAssignmentDal struct {
	CourseLearningOutcomeDal
	AssignmentId string `json:"assignment_id"`
}

type CourseLearningOutcomeWithPO struct {
	CourseLearningOutcome
	ProgramOutcomeCode            string  `json:"program_outcome_code"`
	ProgramOutcomeName            string  `json:"program_outcome_name"`
	ProgramLearningOutcomeCode    string  `json:"program_learning_outcome_code"`
	ExpectedPassingCloPercentage  float64 `json:"expected_passing_clo_percentage"`
	ProgramLearningOutcomeName    string  `json:"program_learning_outcome_name"`
	SubProgramLearningOutcomeCode string  `json:"sub_program_learning_outcome_code"`
	SubProgramLearningOutcomeName string  `json:"sub_program_learning_outcome_name"`
}

type CreateCourseLearningOutcomeDto struct {
	Code                                string   `json:"code" validate:"required"`
	Description                         string   `json:"description" validate:"required"`
	ExpectedPassingAssignmentPercentage float64  `json:"expected_passing_assignment_percentage" validate:"required"`
	ExpectedPassingStudentPercentage    float64  `json:"expected_passing_student_percentage" validate:"required"`
	Status                              string   `json:"status"`
	SubProgramLearningOutcomeIds        []string `json:"sub_program_learning_outcome_ids"`
	SubStudentOutcomeIds                []string `json:"sub_student_outcome_ids"`
	ProgramOutcomeIds                   []string `json:"program_outcome_ids"`
	CourseId                            string
}

type UpdateCourseLeaningOutcomeDto struct {
	Code                                string
	Description                         string
	ExpectedPassingAssignmentPercentage float64
	ExpectedPassingStudentPercentage    float64
	Status                              string
	ProgramOutcomeId                    string
}

type CreateCourseLearningOutcomePayload struct {
	Code                                string   `json:"code" validate:"required"`
	Description                         string   `json:"description" validate:"required"`
	ExpectedPassingAssignmentPercentage float64  `json:"expected_passing_assignment_percentage" validate:"required"`
	ExpectedPassingStudentPercentage    float64  `json:"expected_passing_student_percentage" validate:"required"`
	Status                              string   `json:"status"`
	CourseId                            string   `json:"course_id" validate:"required"`
	ProgramOutcomeIds                   []string `json:"program_outcome_ids"`
	SubProgramLearningOutcomeIds        []string `json:"sub_program_learning_outcome_ids"`
	SubStudentOutcomeIds                []string `json:"sub_student_outcome_ids"`
}

type UpdateCourseLearningOutcomePayload struct {
	Code                                string  `json:"code"`
	Description                         string  `json:"description"`
	ExpectedPassingAssignmentPercentage float64 `json:"expected_passing_assignment_percentage"`
	ExpectedPassingStudentPercentage    float64 `json:"expected_passing_student_percentage"`
	Status                              string  `json:"status"`
}

type CreateLinkProgramOutcomePayload struct {
	ProgramOutcomeIds []string `json:"program_outcome_ids" validate:"required"`
}

type CreateLinkSubProgramLearningOutcomePayload struct {
	SubProgramLearningOutcomeIds []string `json:"sub_program_learning_outcome_ids" validate:"required"`
}

type CreateLinkSubStudentOutcomePayload struct {
	SubStudentOutcomeIds []string `json:"sub_student_outcome_ids" validate:"required"`
}

type StudentPassCLOResp struct {
	Clos   []string     `json:"clos"`
	Result []StudentCLO `json:"result"`
}

type StudentCLO struct {
	StudentID int      `json:"student_id"`
	PassCLO   []string `json:"pass_clo"`
}

type CLOResult struct {
	StudentID                     int     `json:"student_id"`
	CLOID                         string  `json:"clo_id"` // Change to string
	CLOCode                       string  `json:"clo_code"`
	PassedAssignments             int     `json:"passed_assignments"`
	TotalAssignments              int     `json:"total_assignments"`
	ExpectedPassingAssignmentPerc float64 `json:"expected_passing_assignment_percentage"`
}
