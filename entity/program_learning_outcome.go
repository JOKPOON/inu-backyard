package entity

type SubProgramLearningOutcome struct {
	Id                       string `json:"id" gorm:"primaryKey;type:char(255)"`
	Code                     string `json:"code"`
	DescriptionThai          string `json:"descriptionThai"`
	DescriptionEng           string `json:"descriptionEng"`
	ProgramLearningOutcomeId string `json:"programLearningOutcomeId"`
}

type ProgramLearningOutcome struct {
	Id              string `json:"id" gorm:"primaryKey;type:char(255)"`
	Code            string `json:"code"`
	DescriptionThai string `json:"descriptionThai"`
	DescriptionEng  string `json:"descriptionEng"`
	ProgramYear     int    `json:"programYear"`
	ProgrammeId     string `json:"programmeId"`

	SubProgramLearningOutcomes []SubProgramLearningOutcome
	Programme                  Programme
}

type CreateSubProgramLearningOutcomeDto struct {
	Code                     string
	DescriptionThai          string
	DescriptionEng           string
	ProgramLearningOutcomeId string
}

type CrateProgramLearningOutcomeDto struct {
	Code                       string
	DescriptionThai            string
	DescriptionEng             string
	ProgramYear                int
	ProgrammeName              string
	SubProgramLearningOutcomes []CreateSubProgramLearningOutcomeDto
}

type ProgramLearningOutcomeRepository interface {
	GetAll() ([]ProgramLearningOutcome, error)
	GetById(id string) (*ProgramLearningOutcome, error)
	Create(programLearningOutcome *ProgramLearningOutcome) error
	CreateMany(programLearningOutcome []ProgramLearningOutcome) error
	Update(id string, programLearningOutcome *ProgramLearningOutcome) error
	Delete(id string) error

	GetSubPLO(subPloId string) (*SubProgramLearningOutcome, error)
	GetAllSubPlo() ([]SubProgramLearningOutcome, error)
	CreateSubPLO(programLearningOutcome *SubProgramLearningOutcome) error
	UpdateSubPLO(id string, programLearningOutcome *SubProgramLearningOutcome) error
	DeleteSubPLO(id string) error
	FilterExistedSubPLO(subPloIds []string) ([]string, error)
}

type ProgramLearningOutcomeUseCase interface {
	GetAll() ([]ProgramLearningOutcome, error)
	GetById(id string) (*ProgramLearningOutcome, error)
	Create(dto []CrateProgramLearningOutcomeDto) error
	Update(id string, programLearningOutcome *ProgramLearningOutcome) error
	Delete(id string) error

	GetSubPLO(subPloId string) (*SubProgramLearningOutcome, error)
	GetAllSubPlo() ([]SubProgramLearningOutcome, error)
	CreateSubPLO(code string, descriptionThai string, descriptionEng string, programLearningOutcomeId string) error
	UpdateSubPLO(id string, programLearningOutcome *SubProgramLearningOutcome) error
	DeleteSubPLO(id string) error
	FilterNonExistedSubPLO(subPloIds []string) ([]string, error)
}
