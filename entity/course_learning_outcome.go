package entity

type CourseLearningOutcomeRepository interface {
	GetAll() ([]CourseLearningOutcome, error)
	GetById(id string) (*CourseLearningOutcome, error)
	GetByCourseId(courseId string) ([]CourseLearningOutcomeWithPO, error)
	Create(courseLearningOutcome *CourseLearningOutcome) error
	CreateLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeId []string) error
	Update(id string, courseLearningOutcome *CourseLearningOutcome) error
	Delete(id string) error
	DeleteLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeId string) error
	FilterExisted(ids []string) ([]string, error)
}

type CourseLearningOutcomeUseCase interface {
	GetAll() ([]CourseLearningOutcome, error)
	GetById(id string) (*CourseLearningOutcome, error)
	GetByCourseId(courseId string) ([]CourseLearningOutcomeWithPO, error)
	Create(dto CreateCourseLearningOutcomePayload) error
	CreateLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeId []string) error
	Update(id string, dto UpdateCourseLearningOutcomePayload) error
	Delete(id string) error
	DeleteLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeId string) error
	FilterNonExisted(ids []string) ([]string, error)
}

type CourseLearningOutcome struct {
	Id                                  string  `json:"id" gorm:"primaryKey;type:char(255)"`
	Code                                string  `json:"code"`
	Description                         string  `json:"description"`
	ExpectedPassingAssignmentPercentage float64 `json:"expected_passing_assignment_percentage"`
	ExpectedPassingStudentPercentage    float64 `json:"expected_passing_student_percentage"`
	Status                              string  `json:"status"`
	ProgramOutcomeId                    string  `json:"program_outcome_id"`
	CourseId                            string  `json:"course_id"`

	SubProgramLearningOutcomes []*SubProgramLearningOutcome `gorm:"many2many:clo_subplo" json:"sub_program_learning_outcomes"`
	Assignments                []*Assignment                `gorm:"many2many:clo_assignment" json:"-"`
	ProgramOutcome             ProgramOutcome               `json:"-"`
	Course                     Course                       `json:"-"`
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
	Code                                string
	Description                         string
	ExpectedPassingAssignmentPercentage float64
	ExpectedPassingStudentPercentage    float64
	Status                              string
	SubProgramLearningOutcomeIds        []string
	ProgramOutcomeId                    string
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
	Status                              string   `json:"status" validate:"required"`
	ProgramOutcomeId                    string   `json:"program_outcome_id" validate:"required"`
	CourseId                            string   `json:"course_id" validate:"required"`
	SubProgramLearningOutcomeIds        []string `json:"sub_program_learning_outcome_ids" validate:"required"`
}

type UpdateCourseLearningOutcomePayload struct {
	Code                                string  `json:"code"`
	Description                         string  `json:"description"`
	ExpectedPassingAssignmentPercentage float64 `json:"expected_passing_assignment_percentage" validate:"required"`
	ExpectedPassingStudentPercentage    float64 `json:"expected_passing_student_percentage" validate:"required"`
	Status                              string  `json:"status" validate:"required"`
	ProgramOutcomeId                    string  `json:"program_outcome_id" validate:"required"`
}

type CreateLinkSubProgramLearningOutcomePayload struct {
	SubProgramLearningOutcomeId []string `json:"sub_program_learning_outcome_id" validate:"required"`
}
