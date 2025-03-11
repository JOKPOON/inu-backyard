package entity

import (
	"time"

	errs "github.com/team-inu/inu-backyard/entity/error"
)

type CourseStreamType string

const (
	UpCourseStreamType   CourseStreamType = "UPSTREAM"
	DownCourseStreamType CourseStreamType = "DOWNSTREAM"
)

type CourseStreamRepository interface {
	Get(id string) (*CourseStream, error)
	GetByQuery(query CourseStream) ([]CourseStream, error)
	Create(courseStream *CourseStream) error
	Update(id string, courseStream *CourseStream) error
	Delete(id string) error
}

type CourseStreamsUseCase interface {
	Get(id string) (*CourseStream, error)
	GetByFromCourseId(courseId string) ([]CourseStream, error)
	GetByTargetCourseId(courseId string) ([]CourseStream, error)
	Create(CreateCourseStreamPayload) error
	Update(id string, comment string) error
	Delete(id string) error
}

type CourseStream struct {
	Id             string           `json:"id"`
	StreamType     CourseStreamType `json:"stream_type"`
	Comment        string           `json:"comment"`
	SenderId       string           `json:"sender_id"`
	FromCourseId   string           `json:"from_course_id"`
	TargetCourseId string           `json:"target_course_id"`
	CreatedAt      time.Time        `json:"created_at"`

	User         User   `json:"user" gorm:"foreignKey:SenderId"`
	FromCourse   Course `json:"from_course" gorm:"foreignKey:FromCourseId"`
	TargetCourse Course `json:"target_course" gorm:"foreignKey:TargetCourseId"`
}

type CreateCourseStreamPayload struct {
	FromCourseId   string           `json:"from_course_id" validate:"required"`
	TargetCourseId string           `json:"target_course_id" validate:"required"`
	StreamType     CourseStreamType `json:"stream_type" validate:"required"`
	Comment        string           `json:"comment" validate:"required"`
	SenderId       string           `json:"sender_id" validate:"required"`
}

type GetCourseStreamPayload struct {
	FromCourseId   string `json:"from_course_id"`
	TargetCourseId string `json:"target_course_id"`
}

func (p GetCourseStreamPayload) Validate() *errs.DomainError {
	if p.TargetCourseId != "" && p.FromCourseId != "" {
		return errs.New(errs.ErrPayloadValidator, "targetCourseId OR fromCourseId only one")
	}

	if p.TargetCourseId == "" && p.FromCourseId == "" {
		return errs.New(errs.ErrPayloadValidator, "must have query at least targetCourseId OR fromCourseId only one")
	}

	return nil
}
