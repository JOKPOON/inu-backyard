package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

type courseLearningOutcomeUseCase struct {
	courseLearningOutcomeRepo     entity.CourseLearningOutcomeRepository
	courseUseCase                 entity.CourseUseCase
	programOutcomeUseCase         entity.ProgramOutcomeUseCase
	programLearningOutcomeUseCase entity.ProgramLearningOutcomeUseCase
	studentOutcomeUseCase         entity.StudentOutcomeUseCase
}

func NewCourseLearningOutcomeUseCase(
	courseLearningOutcomeRepo entity.CourseLearningOutcomeRepository,
	courseUseCase entity.CourseUseCase,
	programOutcomeUseCase entity.ProgramOutcomeUseCase,
	programLearningOutcomeUseCase entity.ProgramLearningOutcomeUseCase,
	studentOutcomeUseCase entity.StudentOutcomeUseCase,
) entity.CourseLearningOutcomeUseCase {
	return &courseLearningOutcomeUseCase{
		courseLearningOutcomeRepo:     courseLearningOutcomeRepo,
		courseUseCase:                 courseUseCase,
		programOutcomeUseCase:         programOutcomeUseCase,
		programLearningOutcomeUseCase: programLearningOutcomeUseCase,
		studentOutcomeUseCase:         studentOutcomeUseCase,
	}
}

func (u courseLearningOutcomeUseCase) GetAll() ([]entity.CourseLearningOutcome, error) {
	clos, err := u.courseLearningOutcomeRepo.GetAll()
	if err != nil {
		return nil, errs.New(errs.ErrQueryCLO, "cannot get all CLOs", err)
	}

	return clos, nil
}

func (u courseLearningOutcomeUseCase) GetById(id string) (*entity.CourseLearningOutcome, error) {
	clo, err := u.courseLearningOutcomeRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQueryCLO, "cannot get CLO by id %s", id, err)
	}

	return clo, nil
}

func (u courseLearningOutcomeUseCase) GetByCourseId(courseId string) ([]entity.CourseLearningOutcomeWithPO, error) {
	course, err := u.courseUseCase.GetById(courseId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get course id %s while querying clo", courseId, err)
	} else if course == nil {
		return nil, errs.New(errs.ErrCourseNotFound, "course id %s not found while querying clo", courseId)
	}

	clo, err := u.courseLearningOutcomeRepo.GetByCourseId(courseId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryCLO, "cannot get CLO by course id %s", courseId, err)
	}

	return clo, nil
}

func (u courseLearningOutcomeUseCase) Create(payload entity.CreateCourseLearningOutcomePayload) error {
	if payload.ExpectedPassingAssignmentPercentage > 100 || payload.ExpectedPassingAssignmentPercentage < 0 {
		return errs.New(errs.ErrCreateCLO, "expected passing assignment percentage must be between 0 and 100")
	}

	course, err := u.courseUseCase.GetById(payload.CourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get course id %s while creating clo", payload.CourseId, err)
	} else if course == nil {
		return errs.New(errs.ErrCourseNotFound, "course id %s not found while creating clo", payload.CourseId)
	}

	if payload.ProgramOutcomeId != "" {
		po, err := u.programOutcomeUseCase.GetById(payload.ProgramOutcomeId)
		if err != nil {
			return errs.New(errs.SameCode, "cannot get program outcome id %s while creating clo", payload.ProgramOutcomeId, err)
		} else if po == nil {
			return errs.New(errs.ErrCourseNotFound, "program outcome id %s not found while creating clo", payload.ProgramOutcomeId)
		}
	}

	if len(payload.SubProgramLearningOutcomeIds) > 0 {
		nonExistedSubPloIds, err := u.programLearningOutcomeUseCase.FilterNonExistedSubPLO(payload.SubProgramLearningOutcomeIds)
		if err != nil {
			return errs.New(errs.SameCode, "cannot get non existed sub plo ids while creating clo")
		} else if len(nonExistedSubPloIds) != 0 {
			return errs.New(errs.ErrCreateEnrollment, "there are non exist sub plo %v", nonExistedSubPloIds)
		}
	}

	if len(payload.SubStudentOutcomeIds) > 0 {
		nonExistedSubSoIds, err := u.studentOutcomeUseCase.FilterNonExistedSubSO(payload.SubStudentOutcomeIds)
		if err != nil {
			return errs.New(errs.SameCode, "cannot get non existed sub so ids while creating clo")
		} else if len(nonExistedSubSoIds) != 0 {
			return errs.New(errs.ErrCreateEnrollment, "there are non exist sub so %v", nonExistedSubSoIds)
		}
	}

	subPlos := []*entity.SubProgramLearningOutcome{}
	for _, ploId := range payload.SubProgramLearningOutcomeIds {
		subPlos = append(subPlos, &entity.SubProgramLearningOutcome{
			Id: ploId,
		})
	}

	subSos := []*entity.SubStudentOutcome{}
	for _, soId := range payload.SubStudentOutcomeIds {
		subSos = append(subSos, &entity.SubStudentOutcome{
			Id: soId,
		})
	}

	clo := entity.CourseLearningOutcome{
		Id:                                  ulid.Make().String(),
		Code:                                payload.Code,
		Description:                         payload.Description,
		Status:                              payload.Status,
		ExpectedPassingAssignmentPercentage: payload.ExpectedPassingAssignmentPercentage,
		ExpectedPassingStudentPercentage:    payload.ExpectedPassingStudentPercentage,
		CourseId:                            payload.CourseId,
		ProgramOutcomeId:                    payload.ProgramOutcomeId,
		SubProgramLearningOutcomes:          subPlos,
		SubStudentOutcomes:                  subSos,
	}

	err = u.courseLearningOutcomeRepo.Create(&clo)
	if err != nil {
		return errs.New(errs.ErrCreateCLO, "cannot create CLO", err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) CreateLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeIds []string) error {
	existCourseLearningOutcome, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get courseLearningOutcome id %s to link subPLO", id, err)
	}

	if existCourseLearningOutcome == nil {
		return errs.New(errs.ErrCLONotFound, "cannot get courseLearningOutcome id %s to link subPLO", id)
	}

	nonExistedSubPloIds, err := u.programLearningOutcomeUseCase.FilterNonExistedSubPLO(subProgramLearningOutcomeIds)

	if err != nil {
		return errs.New(errs.SameCode, "cannot get non existed sub plo ids while linking clo and sub plo")
	} else if len(nonExistedSubPloIds) != 0 {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist sub plo %v", nonExistedSubPloIds)
	}

	err = u.courseLearningOutcomeRepo.CreateLinkSubProgramLearningOutcome(id, subProgramLearningOutcomeIds)

	if err != nil {
		return errs.New(errs.ErrCreateCLO, "cannot link CLO and subPLO", err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) CreateLinkSubStudentOutcome(id string, subStudentOutcomeIds []string) error {
	existCourseLearningOutcome, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get courseLearningOutcome id %s to link subPLO", id, err)
	} else if existCourseLearningOutcome == nil {
		return errs.New(errs.ErrCLONotFound, "cannot get courseLearningOutcome id %s to link subPLO", id)
	}

	nonExistedSubStudentOutcomeIds, err := u.studentOutcomeUseCase.FilterNonExistedSubSO(subStudentOutcomeIds)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get non existed sub student outcome ids while linking clo and sub student outcome")
	} else if len(nonExistedSubStudentOutcomeIds) != 0 {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist sub student outcome %v", nonExistedSubStudentOutcomeIds)
	}

	err = u.courseLearningOutcomeRepo.CreateLinkSubStudentOutcome(id, subStudentOutcomeIds)
	if err != nil {
		return errs.New(errs.SameCode, "cannot link CLO and sub student outcome", err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) Update(id string, payload entity.UpdateCourseLearningOutcomePayload) error {
	existCourseLearningOutcome, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get courseLearningOutcome id %s to update", id, err)
	} else if existCourseLearningOutcome == nil {
		return errs.New(errs.ErrCLONotFound, "cannot get courseLearningOutcome id %s to update", id)
	}

	existedProgramOutcome, err := u.programOutcomeUseCase.GetById(payload.ProgramOutcomeId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get program outcome id %s to update clo", payload.ProgramOutcomeId, err)
	} else if existedProgramOutcome == nil {
		return errs.New(errs.ErrPONotFound, "program outcome id %s not found while updating clo", payload.ProgramOutcomeId)
	}

	err = u.courseLearningOutcomeRepo.Update(id, &entity.CourseLearningOutcome{
		Code:                                payload.Code,
		Description:                         payload.Description,
		ExpectedPassingAssignmentPercentage: payload.ExpectedPassingAssignmentPercentage,
		ExpectedPassingStudentPercentage:    payload.ExpectedPassingStudentPercentage,
		Status:                              payload.Status,
		ProgramOutcomeId:                    payload.ProgramOutcomeId,
	})

	if err != nil {
		return errs.New(errs.ErrUpdateCLO, "cannot update courseLearningOutcome by id %s", id, err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) Delete(id string) error {
	clo, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get clo id %s to delete", id, err)
	} else if clo == nil {
		return errs.New(errs.ErrAssignmentNotFound, "cannot get clo id %s to delete", id)
	}

	err = u.courseLearningOutcomeRepo.Delete(id)
	if err != nil {
		return errs.New(errs.ErrDeleteCLO, "cannot delete CLO", err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) DeleteLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeId string) error {
	err := u.courseLearningOutcomeRepo.DeleteLinkSubProgramLearningOutcome(id, subProgramLearningOutcomeId)
	if err != nil {
		return errs.New(errs.ErrUnLinkSubPLO, "cannot delete link CLO and subPLO", err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) DeleteLinkSubStudentOutcome(id string, subStudentOutcomeId string) error {
	err := u.courseLearningOutcomeRepo.DeleteLinkSubStudentOutcome(id, subStudentOutcomeId)
	if err != nil {
		return errs.New(errs.ErrUnlinkSubSO, "cannot delete link CLO and sub student outcome", err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) FilterNonExisted(ids []string) ([]string, error) {
	existedIds, err := u.courseLearningOutcomeRepo.FilterExisted(ids)
	if err != nil {
		return nil, errs.New(errs.ErrQueryCLO, "cannot query clo", err)
	}

	nonExistedIds := slice.Subtraction(ids, existedIds)

	return nonExistedIds, nil
}
