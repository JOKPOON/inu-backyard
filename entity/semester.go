package entity

type Semester struct {
	Id               string `json:"id" gorm:"primaryKey;type:char(255)"`
	Year             int    `json:"year"`
	SemesterSequence string `json:"semester_sequence"`
}

type CreateSemesterPayload struct {
	Year             int    `json:"year" validate:"required"`
	SemesterSequence string `json:"semester_sequence" validate:"required"`
}

type UpdateSemesterPayload struct {
	Year             int    `json:"year" validate:"required"`
	SemesterSequence string `json:"semester_sequence" validate:"required"`
}

type SemesterRepository interface {
	GetAll() ([]Semester, error)
	Get(year int, semesterSequence string) (*Semester, error)
	GetById(id string) (*Semester, error)
	Create(semester *Semester) error
	Update(semester *Semester) error
	Delete(id string) error
}

type SemesterUseCase interface {
	GetAll() ([]Semester, error)
	GetById(id string) (*Semester, error)
	Get(year int, semesterSequence string) (*Semester, error)
	Create(year int, semesterSequence string) error
	Update(semester *Semester) error
	Delete(id string) error
}
