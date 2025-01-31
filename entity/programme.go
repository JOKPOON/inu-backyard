package entity

import "gorm.io/datatypes"

type ProgrammeRepository interface {
	GetAll() ([]Programme, error)
	GetByName(nameTH string, nameEN string) ([]Programme, error)
	GetByNameAndYear(nameTH string, nameEN string, year string) (*Programme, error)
	GetById(id string) (*Programme, error)
	GetAllCourseOutcomeLinked(programmeId string) ([]CourseOutcomes, error)
	GetAllCourseLinkedPO(programmeId string) (*ProgrammeLinkedPO, error)
	GetAllCourseLinkedPLO(programmeId string) (*ProgrammeLinkedPLO, error)
	GetAllCourseLinkedSO(programmeId string) (*ProgrammeLinkedSO, error)
	Create(programme *Programme) error
	Update(name string, programme *Programme) error
	Delete(name string) error
	FilterExisted(nameTH []string, nameEN []string) ([]string, error)

	CreateLinkWithPO(programmeId string, poId string) error
	CreateLinkWithPLO(programmeId string, ploId string) error
	CreateLinkWithSO(programmeId string, soId string) error
}

type ProgrammeUseCase interface {
	GetAll() ([]Programme, error)
	GetByName(nameTH string, nameEN string) ([]Programme, error)
	GetByNameAndYear(nameTH string, nameEN string, year string) (*Programme, error)
	GetById(id string) (*Programme, error)
	GetAllCourseOutcomeLinked(programmeId string) ([]CourseOutcomes, error)
	GetAllCourseLinkedPO(programmeId string) (*ProgrammeLinkedPO, error)
	GetAllCourseLinkedPLO(programmeId string) (*ProgrammeLinkedPLO, error)
	GetAllCourseLinkedSO(programmeId string) (*ProgrammeLinkedSO, error)
	Create(payload CreateProgrammePayload) error
	Update(name string, programme *UpdateProgrammePayload) error
	Delete(name string) error
	FilterNonExisted(namesTH []string, namesEN []string) ([]string, error)

	CreateLinkWithPO(programmeId string, poIds []string) error
	CreateLinkWithPLO(programmeId string, ploIds []string) error
	CreateLinkWithSO(programmeId string, soIds []string) error
}

type Category struct {
	Name   string `json:"name" validate:"required"`
	Credit int    `json:"credit" validate:"required"`
	SubCat []struct {
		Name   string `json:"name"`
		Credit int    `json:"credit"`
	} `json:"sub"`
}

type ProgrammeStructure struct {
	Category     []Category `json:"category" validate:"required"`
	TotalsCredit int        `json:"totals_credit" validate:"required"`
}

type Programme struct {
	Id            string `json:"id" gorm:"primaryKey;type:char(255)"`
	NameTH        string `json:"name_th" gorm:"not null"`
	NameEN        string `json:"name_en" gorm:"not null"`
	DegreeTH      string `json:"degree_th" gorm:"not null"`
	DegreeEN      string `json:"degree_en" gorm:"not null"`
	DegreeShortTH string `json:"degree_short_th" gorm:"not null"`
	DegreeShortEN string `json:"degree_short_en" gorm:"not null"`
	Year          string `json:"year" gorm:"not null"`

	Structure datatypes.JSON `json:"structure" gorm:"type:json"`

	ProgramOutcomes         []*ProgramOutcome         `gorm:"many2many:programme_po" json:"program_outcomes"`
	ProgramLearningOutcomes []*ProgramLearningOutcome `gorm:"many2many:programme_plo" json:"program_learning_outcomes"`
	StudentOutcomes         []*StudentOutcome         `gorm:"many2many:programme_so" json:"student_outcomes"`
}

type CreateProgrammePayload struct {
	NameTH        string `json:"name_th" validate:"required"`
	NameEN        string `json:"name_en" validate:"required"`
	DegreeTH      string `json:"degree_th" validate:"required"`
	DegreeEN      string `json:"degree_en" validate:"required"`
	DegreeShortTH string `json:"degree_short_th" validate:"required"`
	DegreeShortEN string `json:"degree_short_en" validate:"required"`
	Year          string `json:"year" validate:"required"`

	Structure ProgrammeStructure `json:"structure" validate:"required"`
}

type UpdateProgrammePayload struct {
	NameTH        string `json:"name_th" validate:"required"`
	NameEN        string `json:"name_en" validate:"required"`
	DegreeTH      string `json:"degree_th" validate:"required"`
	DegreeEN      string `json:"degree_en" validate:"required"`
	DegreeShortTH string `json:"degree_short_th" validate:"required"`
	DegreeShortEN string `json:"degree_short_en" validate:"required"`
	Year          string `json:"year" validate:"required"`

	Structure ProgrammeStructure `json:"structure" validate:"required"`
}

type AllCourseOutcome struct {
	CourseCode string `json:"course_code"`
	CourseName string `json:"course_name"`
	Program    string `json:"program"`
	POs        []struct {
		POCode string `json:"po_code"`
	} `json:"pos"`
	CLOs []struct {
		CLOCode string `json:"clo_code"`
	} `json:"clos"`
	PLOs []struct {
		PLOCode string `json:"plo_code"`
		SPLOs   []struct {
			SPLOCode string `json:"splo_code"`
		} `json:"splos"`
	} `json:"plos"`
	SOs []struct {
		SOCode string `json:"so_code"`
		SSOs   []struct {
			SSOCode string `json:"sso_code"`
		} `json:"ssos"`
	} `json:"sos"`
}

type ProgrammeOutcomes struct {
	POs            map[string][]string // PO -> List of PLOs
	PLO_SPLO       map[string][]string // PLO -> List of SPLOs
	SO_SSO         map[string][]string // SO -> List of SSOs
	CourseOutcomes []CourseOutcomes    // Course -> List of CLOs, POs, PLOs, SOs
}

type CourseOutcomes struct {
	CourseCode string
	CourseName string
	CLOs       []string
	POs        []string
	PLOs       map[string][]string // PLO -> List of SPLOs
	SOs        map[string][]string // SO -> List of SSOs
}

type ProgrammeLinkedPO struct {
	ProgrammeName   string
	ProgrammeYear   string
	AllPOs          []string
	AllCourse       []string
	CourseLinkedPOs []CourseLinkedPO
}

type ProgrammeLinkedPLO struct {
	ProgrammeName    string
	ProgrammeYear    string
	AllPLOs          map[string][]string
	AllCourse        []string
	CourseLinkedPLOs []CourseLinkedPLO
}

type ProgrammeLinkedSO struct {
	ProgrammeName   string
	ProgrammeYear   string
	AllSOs          map[string][]string
	AllCourse       []string
	CourseLinkedSOs []CourseLinkedSO
}

type CourseLinkedPO struct {
	CourseCode string
	CourseName string
	POs        []string
}

type CourseLinkedPLO struct {
	CourseCode string
	CourseName string
	PLOs       map[string][]string
}

type CourseLinkedSO struct {
	CourseCode string
	CourseName string
	SOs        map[string][]string
}

type LinkProgrammePO struct {
	POIds []string `json:"po_ids" validate:"required"`
}

type LinkProgrammePLO struct {
	PLOIds []string `json:"plo_ids" validate:"required"`
}

type LinkProgrammeSO struct {
	SOIds []string `json:"so_ids" validate:"required"`
}
