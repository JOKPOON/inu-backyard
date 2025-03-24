package entity

type StudentOutcomeRepository interface {
	GetAll(programId string) ([]StudentOutcome, error)
	GetById(id string) (*StudentOutcome, error)
	Create(studentOutcome *StudentOutcome) error
	CreateMany(studentOutcome []*StudentOutcome) error
	Update(id string, studentOutcome *StudentOutcome) error
	Delete(id string) error
	FilterExisted(ids []string) ([]string, error)

	GetAllSubSO() ([]SubStudentOutcome, error)
	GetSubSOById(id string) (*SubStudentOutcome, error)
	UpdateSubSO(id string, payload *SubStudentOutcome) error
	DeleteSubSO(id string) error
	CreateSubSO(subStudentOutcome *SubStudentOutcome) error
	CreateManySubSO(subStudentOutcome []*SubStudentOutcome) error
	FilterExistedSubSO(ids []string) ([]string, error)
}

type StudentOutcomeUseCase interface {
	GetAll(programId string) ([]StudentOutcome, error)
	GetById(id string) (*StudentOutcome, error)
	Create(payload []CreateStudentOutcome) error
	Update(id string, payload *UpdateStudentOutcomePayload) error
	Delete(id string) error
	FilterNonExisted(ids []string) ([]string, error)

	GetAllSubSO() ([]SubStudentOutcome, error)
	GetSubSOById(id string) (*SubStudentOutcome, error)
	UpdateSubSO(id string, payload *UpdateSubStudentOutcomePayload) error
	DeleteSubSO(id string) error
	CreateSubSO(payload []CreateSubStudentOutcome) error
	FilterNonExistedSubSO(ids []string) ([]string, error)
}

type StudentOutcome struct {
	Id              string `gorm:"primaryKey;type:char(255)" json:"id"`
	Code            string `validate:"required" json:"code"`
	DescriptionThai string `validate:"required" json:"description_thai"`
	DescriptionEng  string `json:"description_eng"`
	ProgramId       string `json:"program_id" gorm:"type:char(255)"`

	SubStudentOutcomes []SubStudentOutcome `json:"sub_student_outcomes"`
}

type SubStudentOutcome struct {
	Id               string `gorm:"primaryKey;type:char(255)" json:"id"`
	Code             string `validate:"required" json:"code"`
	DescriptionThai  string `validate:"required" json:"description_thai"`
	DescriptionEng   string `json:"description_eng"`
	StudentOutcomeId string `json:"student_outcome_id" gorm:"type:char(255)"`
}

type CreateSubStudentOutcome struct {
	Code             string `validate:"required" json:"code"`
	DescriptionThai  string `validate:"required" json:"description_thai"`
	DescriptionEng   string `json:"description_eng"`
	StudentOutcomeId string `json:"student_outcome_id"`
}

type CreateSubStudentOutcomePayload struct {
	SubStudentOutcomes []CreateSubStudentOutcome `json:"sub_student_outcomes" validate:"required,dive"`
}

type CreateStudentOutcome struct {
	Code               string                    `validate:"required" json:"code"`
	DescriptionThai    string                    `validate:"required" json:"description_thai"`
	DescriptionEng     string                    `json:"description_eng"`
	SubStudentOutcomes []CreateSubStudentOutcome `json:"sub_student_outcomes" validate:"required,dive"`
	ProgramId          string                    `json:"program_id" validate:"required"`
}

type CreateStudentOutcomePayload struct {
	StudentOutcomes []CreateStudentOutcome `json:"student_outcomes" validate:"required,dive"`
}

type UpdateStudentOutcomePayload struct {
	Id              string `json:"id" validate:"required"`
	Code            string `json:"code"`
	DescriptionThai string `json:"description_thai"`
	DescriptionEng  string `json:"description_eng"`
	ProgramId       string `json:"program_id"`
}

type UpdateSubStudentOutcomePayload struct {
	Code             string `json:"code"`
	DescriptionThai  string `json:"description_thai"`
	DescriptionEng   string `json:"description_eng"`
	StudentOutcomeId string `json:"student_outcome_id"`
}
