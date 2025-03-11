package usecase

import (
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
)

type courseStreamUseCase struct {
	courseStreamRepository entity.CourseStreamRepository
	courseUseCase          entity.CourseUseCase
}

func NewCourseStreamUseCase(
	courseStreamRepository entity.CourseStreamRepository,
	courseUseCase entity.CourseUseCase,
) entity.CourseStreamsUseCase {
	return &courseStreamUseCase{
		courseStreamRepository: courseStreamRepository,
		courseUseCase:          courseUseCase,
	}
}

func (u *courseStreamUseCase) Create(payload entity.CreateCourseStreamPayload) error {
	fromCourse, err := u.courseUseCase.GetById(payload.FromCourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot validate from course id %s while creating stream course", payload.TargetCourseId, err)
	} else if fromCourse == nil {
		return errs.New(errs.ErrCreateCourseStream, "from course id: %s not found", payload.FromCourseId)
	}

	targetCourse, err := u.courseUseCase.GetById(payload.TargetCourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot validate target course id %s while creating stream course", payload.TargetCourseId, err)
	} else if targetCourse == nil {
		return errs.New(errs.ErrCreateCourseStream, "target course id: %s not found", payload.TargetCourseId)
	}

	courseStream := entity.CourseStream{
		Id:             ulid.Make().String(),
		FromCourseId:   payload.FromCourseId,
		TargetCourseId: payload.TargetCourseId,
		StreamType:     payload.StreamType,
		Comment:        payload.Comment,
		SenderId:       payload.SenderId,
		CreatedAt:      time.Now(),
	}

	err = u.courseStreamRepository.Create(&courseStream)
	if err != nil {
		return errs.New(errs.ErrCreateCourseStream, "cannot create stream course", err)
	}

	return nil
}

func (u *courseStreamUseCase) Delete(id string) error {
	courseStream, err := u.Get(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot validate course stream id %s to delete", id, err)
	} else if courseStream == nil {
		return errs.New(errs.ErrDeleteCourseStream, "course stream id %s not found while deleting", id)
	}

	err = u.courseStreamRepository.Delete(id)
	if err != nil {
		return errs.New(errs.ErrDeleteCourseStream, "cannot delete stream course", err)
	}

	return nil
}

func (u *courseStreamUseCase) GetByTargetCourseId(courseId string) ([]entity.CourseStream, error) {
	course, err := u.courseUseCase.GetById(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot validate course id %s while getting stream course", courseId, err)
	} else if course == nil {
		return nil, errs.New(errs.ErrQueryCourseStream, "course id %s not found while getting stream course", courseId)
	}

	courseStream, err := u.courseStreamRepository.GetByQuery(entity.CourseStream{TargetCourseId: courseId})
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get stream course of course id: %s", courseStream, err)
	}

	return courseStream, nil
}

func (u *courseStreamUseCase) GetByFromCourseId(courseId string) ([]entity.CourseStream, error) {
	course, err := u.courseUseCase.GetById(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot validate course id %s while getting stream course", courseId, err)
	} else if course == nil {
		return nil, errs.New(errs.ErrQueryCourseStream, "course id %s not found while getting stream course", courseId)
	}

	courseStream, err := u.courseStreamRepository.GetByQuery(entity.CourseStream{FromCourseId: courseId})
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get stream course of course id: %s", courseStream, err)
	}

	return courseStream, nil
}

func (u *courseStreamUseCase) Get(id string) (*entity.CourseStream, error) {
	courseStream, err := u.courseStreamRepository.Get(id)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get stream course id %s", courseStream, err)
	}

	return courseStream, nil
}

func (u *courseStreamUseCase) Update(id string, comment string) error {
	courseStream, err := u.Get(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot validate stream course id %s to update", id, err)
	} else if courseStream == nil {
		return errs.New(errs.ErrCourseStreamNotFound, "stream course id %s to update not found", id)
	}

	err = u.courseStreamRepository.Update(id, &entity.CourseStream{Comment: comment})
	if err != nil {
		return errs.New(errs.ErrUpdateCourseStream, "cannot update stream course id %s", id)
	}

	return nil
}
