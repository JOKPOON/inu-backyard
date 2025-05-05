package entity

type ProgramLearningOutcomeRepository interface {
	GetAll(programId string) ([]ProgramLearningOutcome, error)
	GetById(id string) (*ProgramLearningOutcome, error)
	Create(programLearningOutcome *ProgramLearningOutcome) error
	CreateMany(programLearningOutcome []ProgramLearningOutcome) error
	Update(id string, programLearningOutcome *ProgramLearningOutcome) error
	Delete(id string) error
	FilterExisted(ids []string) ([]string, error)

	GetSubPLO(subPloId string) (*SubProgramLearningOutcome, error)
	GetSubPloByPloId(ploId string) ([]SubProgramLearningOutcome, error)
	GetSubPloByCode(code string, programme string, year int) (*SubProgramLearningOutcome, error)
	GetAllSubPlo() ([]SubProgramLearningOutcome, error)
	CreateSubPLO(programLearningOutcome []SubProgramLearningOutcome) error
	UpdateSubPLO(id string, programLearningOutcome *SubProgramLearningOutcome) error
	DeleteSubPLO(id string) error
	FilterExistedSubPLO(subPloIds []string) ([]string, error)
}

type ProgramLearningOutcomeUseCase interface {
	GetAll(programId string) ([]ProgramLearningOutcome, error)
	GetById(id string) (*ProgramLearningOutcome, error)
	Create(dto []CreateProgramLearningOutcome) error
	Update(id string, programLearningOutcome *ProgramLearningOutcome) error
	Delete(id string) error
	FilterNonExisted(ids []string) ([]string, error)

	GetSubPLO(subPloId string) (*SubProgramLearningOutcome, error)
	GetSubPloByPloId(ploId string) ([]SubProgramLearningOutcome, error)
	GetSubPloByCode(code string, programme string, year int) (*SubProgramLearningOutcome, error)
	GetAllSubPlo() ([]SubProgramLearningOutcome, error)
	CreateSubPLO(dto []CreateSubProgramLearningOutcome) error
	UpdateSubPLO(id string, programLearningOutcome *SubProgramLearningOutcome) error
	DeleteSubPLO(id string) error
	FilterNonExistedSubPLO(subPloIds []string) ([]string, error)
}

type ProgramLearningOutcome struct {
	Id                              string  `json:"id" gorm:"primaryKey;type:char(255)"`
	Code                            string  `json:"code"`
	DescriptionThai                 string  `json:"description_thai"`
	DescriptionEng                  string  `json:"description_eng"`
	ExpectedCoursePassingPercentage float64 `json:"expected_course_passing_percentage"`
	ProgramId                       string  `json:"program_id" gorm:"type:char(255)"`

	SubProgramLearningOutcomes []SubProgramLearningOutcome `json:"sub_program_learning_outcomes"`
	Programme                  *Programme                  `json:"programme,omitempty" gorm:"foreignKey:ProgramId"`
}

type CreateProgramLearningOutcome struct {
	Code                            string                            `validate:"required" json:"code"`
	DescriptionThai                 string                            `validate:"required" json:"description_thai"`
	DescriptionEng                  string                            `json:"description_eng"`
	ExpectedCoursePassingPercentage float64                           `validate:"required" json:"expected_course_passing_percentage"`
	ProgramId                       string                            `validate:"required" json:"program_id"`
	SubProgramLearningOutcomes      []CreateSubProgramLearningOutcome `json:"sub_program_learning_outcomes"`
}

type CreateProgramLearningOutcomePayload struct {
	ProgramLearningOutcomes []CreateProgramLearningOutcome `json:"program_learning_outcomes" validate:"required,dive"`
}

type UpdateProgramLearningOutcomePayload struct {
	Id                              string  `json:"id" validate:"required" `
	Code                            string  `json:"code"`
	DescriptionThai                 string  `json:"description_thai"`
	DescriptionEng                  *string `json:"description_eng"`
	ExpectedCoursePassingPercentage float64 `json:"expected_course_passing_percentage"`
	ProgramId                       string  `json:"program_id"`
}

type SubProgramLearningOutcome struct {
	Id                       string `json:"id" gorm:"primaryKey;type:char(255)"`
	Code                     string `json:"code"`
	DescriptionThai          string `json:"description_thai"`
	DescriptionEng           string `json:"description_eng"`
	ProgramLearningOutcomeId string `json:"program_learning_outcome_id"`
}

type CreateSubProgramLearningOutcomeDto struct {
	Code                     string `json:"code"`
	DescriptionThai          string `json:"description_thai"`
	DescriptionEng           string `json:"description_eng"`
	ProgramLearningOutcomeId string `json:"program_learning_outcome_id"`
}

type CreateSubProgramLearningOutcome struct {
	Code                     string `validate:"required" json:"code"`
	DescriptionThai          string `validate:"required" json:"description_thai"`
	DescriptionEng           string `json:"description_eng"`
	ProgramLearningOutcomeId string `validate:"required" json:"program_learning_outcome_id"`
}

type CreateSubProgramLearningOutcomePayload struct {
	SubProgramLearningOutcomes []CreateSubProgramLearningOutcome `json:"sub_program_learning_outcomes" validate:"required,dive"`
}

type UpdateSubProgramLearningOutcome struct {
	Id                       string  `json:"id" validate:"required"`
	Code                     string  `json:"code"`
	DescriptionThai          string  `json:"description_thai"`
	DescriptionEng           *string `json:"description_eng"`
	ProgramLearningOutcomeId string  `json:"program_learning_outcome_id"`
}

type UpdateSubProgramLearningOutcomePayload struct {
	SubProgramLearningOutcomes []UpdateSubProgramLearningOutcome `json:"sub_program_learning_outcomes" validate:"required,dive"`
}
