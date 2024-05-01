package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
)

func (u assignmentUseCase) GetGroupByGroupId(assignmentGroupId string) (*entity.AssignmentGroup, error) {
	assignmentGroup, err := u.assignmentRepo.GetGroupByGroupId(assignmentGroupId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get assignment group id %s", assignmentGroup, err)
	}

	return assignmentGroup, nil
}

func (u assignmentUseCase) CreateGroup(name string, courseId string) error {

	course, err := u.courseUseCase.GetById(courseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot validate course id %s while creating assignment group", courseId, err)
	} else if course == nil {
		return errs.New(errs.ErrCourseNotFound, "course id %s now found while creating assignment group", courseId)
	}

	assignment := entity.AssignmentGroup{
		Id:       ulid.Make().String(),
		Name:     name,
		CourseId: courseId,
	}

	err = u.assignmentRepo.CreateGroup(&assignment)
	if err != nil {
		return errs.New(errs.ErrCreateAssignment, "cannot create assignment group", err)
	}

	return nil
}
