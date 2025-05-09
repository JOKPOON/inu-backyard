package entity

type ProgramOutcomeRepository interface {
	GetAll(programId string) ([]ProgramOutcome, error)
	GetById(id string) (*ProgramOutcome, error)
	GetByCode(code string) (*ProgramOutcome, error)
	Create(programOutcome *ProgramOutcome) error
	CreateMany(programOutcome []ProgramOutcome) error
	Update(id string, programOutcome *ProgramOutcome) error
	Delete(id string) error
	FilterExisted(ids []string) ([]string, error)
}

type ProgramOutcomeUseCase interface {
	GetAll(programId string) ([]ProgramOutcome, error)
	GetById(id string) (*ProgramOutcome, error)
	GetByCode(code string) (*ProgramOutcome, error)
	Create(programOutcome []CreateProgramOutcome) error
	Update(id string, programOutcome *ProgramOutcome) error
	Delete(id string) error
	FilterNonExisted(ids []string) ([]string, error)
}

type ProgramOutcome struct {
	Id                              string  `json:"id" gorm:"primaryKey;type:char(255)"`
	Code                            string  `json:"code"`
	Description                     string  `json:"description"`
	ExpectedCoursePassingPercentage float64 `json:"expected_course_passing_percentage"`
	Category                        string  `json:"category"`
	ProgramId                       string  `json:"program_id" gorm:"type:char(255)"`

	Programme *Programme `json:"programme,omitempty" gorm:"foreignKey:ProgramId"`
}

type CreateProgramOutcome struct {
	Code                            string  `json:"code" validate:"required"`
	Description                     string  `json:"description" validate:"required"`
	ExpectedCoursePassingPercentage float64 `json:"expected_course_passing_percentage" validate:"required"`
	Category                        string  `json:"category" validate:"required"`
	ProgramId                       string  `json:"program_id" validate:"required"`
}

type CreateProgramOutcomePayload struct {
	ProgramOutcomes []CreateProgramOutcome `json:"program_outcomes" validate:"required,dive"`
}

type UpdateProgramOutcomePayload struct {
	Id                              string  `json:"id" validate:"required"`
	Code                            string  `json:"code" `
	Description                     string  `json:"description" `
	ExpectedCoursePassingPercentage float64 `json:"expected_course_passing_percentage" `
	Category                        string  `json:"category" `
	ProgramId                       string  `json:"program_id" `
}
