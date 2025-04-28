package entity

type CourseLearningOutcomeRepository interface {
	GetAll() ([]CourseLearningOutcome, error)
	GetById(id string) (*CourseLearningOutcome, error)
	GetByCourseId(courseId string) ([]CourseLearningOutcome, error)
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
	GetProgramLearningOutcomesBySubProgramLearningOutcomeId(sploIds []string) ([]ProgramLearningOutcome, error)
	GetStudentOutcomesBySubStudentOutcomeId(subStudentOutcomeIds []string) ([]StudentOutcome, error)
}

type CourseLearningOutcomeUseCase interface {
	GetAll() ([]GetCloResponse, error)
	GetById(id string) (*CourseLearningOutcome, error)
	GetByCourseId(courseId string) ([]GetCloResponse, error)
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
	DescriptionTH                       string  `json:"description_th"`
	DescriptionEN                       string  `json:"description_en"`
	Type                                string  `json:"type"`
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
	DescriptionTH                       string   `json:"description_th"`
	DescriptionEN                       string   `json:"description_en"`
	Type                                string   `json:"type" validate:"required"`
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
	DescriptionTH                       string
	DescriptionEN                       string
	Type                                string
	ExpectedPassingAssignmentPercentage float64
	ExpectedPassingStudentPercentage    float64
	Status                              string
	ProgramOutcomeId                    string
}

type CreateCourseLearningOutcomePayload struct {
	Code                                string   `json:"code" validate:"required"`
	DescriptionTH                       string   `json:"description_th"`
	DescriptionEN                       string   `json:"description_en"`
	Type                                string   `json:"type" validate:"required"`
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
	DescriptionTH                       string  `json:"description_th"`
	DescriptionEN                       string  `json:"description_en"`
	Type                                string  `json:"type"`
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
	StudentID     int           `json:"student_id"`
	StudentNameTH string        `json:"student_name_th"`
	StudentNameEN string        `json:"student_name_en"`
	CLOs          []PassOutcome `json:"clos"`
	POs           []PassOutcome `json:"pos"`
	PLOs          []PassOutcome `json:"plos"`
	SOs           []PassOutcome `json:"sos"`
}

type CLOResult struct {
	StudentId                        int     `json:"student_id"`
	StudentFirstNameTH               string  `json:"student_first_name_th"`
	StudentLastNameTH                string  `json:"student_last_name_th"`
	StudentFirstNameEN               string  `json:"student_first_name_en"`
	StudentLastNameEN                string  `json:"student_last_name_en"`
	CLOId                            string  `json:"clo_id"` // Change to string
	CLOCode                          string  `json:"clo_code"`
	Score                            float64 `json:"score"`
	MaxScore                         int     `json:"max_score"`
	ExpectedPassingAssignmentPercent float64 `gorm:"column:expected_passing_assignment_percentage" json:"expected_passing_assignment_percentage"`
	ExpectedScorePercent             float64 `gorm:"column:expected_score_percentage" json:"expected_score_percentage"`
}

type GetCloResponse struct {
	Id                                  string  `json:"id" gorm:"primaryKey;type:char(255)"`
	Code                                string  `json:"code"`
	DescriptionTH                       string  `json:"description_th"`
	DescriptionEN                       string  `json:"description_en"`
	Type                                string  `json:"type"`
	ExpectedPassingAssignmentPercentage float64 `json:"expected_passing_assignment_percentage"`
	ExpectedPassingStudentPercentage    float64 `json:"expected_passing_student_percentage"`
	Status                              string  `json:"status"`

	ProgramOutcomes         []*ProgramOutcome        ` json:"program_outcomes"`
	ProgramLearningOutcomes []ProgramLearningOutcome `json:"program_learning_outcomes"`
	SubStudentOutcomes      []StudentOutcome         ` json:"student_outcomes"`
}
