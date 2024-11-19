package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
)

func (u assignmentUseCase) GetAllGroup() ([]entity.AssignmentGroup, error) {
	assignmentGroup, err := u.assignmentRepo.GetAllGroup()
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get assignment group", assignmentGroup, err)
	}

	return assignmentGroup, nil
}

func (u assignmentUseCase) GetGroupByGroupId(assignmentGroupId string) (*entity.AssignmentGroup, error) {
	assignmentGroup, err := u.assignmentRepo.GetGroupByGroupId(assignmentGroupId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get assignment group id %s", assignmentGroup, err)
	}

	return assignmentGroup, nil
}

func (u assignmentUseCase) GetGroupByCourseId(courseId string, withAssignment bool) ([]entity.AssignmentGroup, error) {
	course, err := u.courseUseCase.GetById(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course id %s while get assignments", course, err)
	} else if course == nil {
		return nil, errs.New(errs.ErrCourseNotFound, "course id %s not found while getting assignments", courseId, err)
	}

	assignmentGroup, err := u.assignmentRepo.GetGroupByQuery(entity.AssignmentGroup{CourseId: courseId}, withAssignment)
	if err != nil {
		return nil, errs.New(errs.ErrQueryAssignment, "cannot get assignment group by course id %s", courseId, err)
	}

	return assignmentGroup, nil
}

func (u assignmentUseCase) CreateGroup(payload entity.CreateAssignmentGroupPayload) error {
	course, err := u.courseUseCase.GetById(payload.CourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot validate course id %s while creating assignment group", payload.CourseId, err)
	} else if course == nil {
		return errs.New(errs.ErrCourseNotFound, "course id %s now found while creating assignment group", payload.CourseId)
	}

	assignment := entity.AssignmentGroup{
		Id:       ulid.Make().String(),
		Name:     payload.Name,
		CourseId: payload.CourseId,
		Weight:   payload.Weight,
	}

	err = u.assignmentRepo.CreateGroup(&assignment)
	if err != nil {
		return errs.New(errs.ErrCreateAssignment, "cannot create assignment group", err)
	}

	return nil
}

func (u assignmentUseCase) UpdateGroup(assignmentGroupId string, payload entity.UpdateAssignmentGroupPayload) error {
	assignmentGroup, err := u.GetGroupByGroupId(assignmentGroupId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot validate assignment group id %s to update", assignmentGroupId, err)
	} else if assignmentGroup == nil {
		return errs.New(errs.ErrAssignmentNotFound, "assignment group id %s to update not found", assignmentGroupId)
	}

	err = u.assignmentRepo.UpdateGroup(assignmentGroupId, &entity.AssignmentGroup{
		Name:   payload.Name,
		Weight: payload.Weight,
	})
	if err != nil {
		return errs.New(errs.ErrUpdateAssignment, "cannot update assignment group id %s", assignmentGroupId)
	}

	return nil
}

func (u assignmentUseCase) DeleteGroup(assignmentGroupId string) error {
	assignmentGroup, err := u.GetGroupByGroupId(assignmentGroupId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot validate assignment group id %s to delete", assignmentGroupId, err)
	} else if assignmentGroup == nil {
		return errs.New(errs.ErrAssignmentNotFound, "assignment group id %s not found while deleting", assignmentGroupId)
	}

	err = u.assignmentRepo.DeleteGroup(assignmentGroupId)
	if err != nil {
		return errs.New(errs.ErrDeleteAssignment, "cannot delete assignment group", err)
	}

	return nil
}
