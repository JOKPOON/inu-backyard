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
	programmeUseCase              entity.ProgrammeUseCase
	programOutcomeUseCase         entity.ProgramOutcomeUseCase
	programLearningOutcomeUseCase entity.ProgramLearningOutcomeUseCase
	studentOutcomeUseCase         entity.StudentOutcomeUseCase
}

func NewCourseLearningOutcomeUseCase(
	courseLearningOutcomeRepo entity.CourseLearningOutcomeRepository,
	courseUseCase entity.CourseUseCase,
	programmeUseCase entity.ProgrammeUseCase,
	programOutcomeUseCase entity.ProgramOutcomeUseCase,
	programLearningOutcomeUseCase entity.ProgramLearningOutcomeUseCase,
	studentOutcomeUseCase entity.StudentOutcomeUseCase,
) entity.CourseLearningOutcomeUseCase {
	return &courseLearningOutcomeUseCase{
		courseLearningOutcomeRepo:     courseLearningOutcomeRepo,
		courseUseCase:                 courseUseCase,
		programmeUseCase:              programmeUseCase,
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

	if len(payload.ProgramOutcomeIds) > 0 {
		for _, programOutcomeId := range payload.ProgramOutcomeIds {
			po, err := u.programOutcomeUseCase.GetById(programOutcomeId)
			if err != nil {
				return errs.New(errs.SameCode, "cannot get program outcome id %s while creating clo", programOutcomeId, err)
			} else if po == nil {
				return errs.New(errs.ErrCourseNotFound, "program outcome id %s not found while creating clo", programOutcomeId)
			}
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

	pos := []*entity.ProgramOutcome{}
	for _, poId := range payload.ProgramOutcomeIds {
		pos = append(pos, &entity.ProgramOutcome{
			Id: poId,
		})
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
		ProgramOutcomes:                     pos,
		SubProgramLearningOutcomes:          subPlos,
		SubStudentOutcomes:                  subSos,
	}

	err = u.courseLearningOutcomeRepo.Create(&clo)
	if err != nil {
		return errs.New(errs.ErrCreateCLO, "cannot create CLO", err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) CreateLinkProgramOutcome(id string, programOutcomeIds []string) error {
	existCourseLearningOutcome, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get courseLearningOutcome id %s to link programOutcome", id, err)
	} else if existCourseLearningOutcome == nil {
		return errs.New(errs.ErrCLONotFound, "cannot get courseLearningOutcome id %s to link programOutcome", id)
	}

	course, err := u.courseUseCase.GetById(existCourseLearningOutcome.CourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get course id %s while linking clo and program outcome", existCourseLearningOutcome.CourseId, err)
	} else if course == nil {
		return errs.New(errs.ErrCourseNotFound, "course id %s not found while linking clo and program outcome", existCourseLearningOutcome.CourseId)
	}

	allPOs, err := u.programmeUseCase.GetAllPO(course.ProgrammeId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get program outcome while linking clo and program outcome", err)
	}

	allPOIds := []string{}
	for _, po := range allPOs {
		allPOIds = append(allPOIds, po.Id)
	}

	useablePOIds := slice.Intersection(programOutcomeIds, allPOIds)
	if len(useablePOIds) != len(programOutcomeIds) {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist po %v", slice.Subtraction(programOutcomeIds, useablePOIds))
	}

	err = u.courseLearningOutcomeRepo.CreateLinkProgramOutcome(id, programOutcomeIds)
	if err != nil {
		return errs.New(errs.ErrCreateCLO, "cannot link CLO and program outcome", err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) CreateLinkSubProgramLearningOutcome(id string, subProgramLearningOutcomeIds []string) error {
	existCourseLearningOutcome, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get courseLearningOutcome id %s to link subPLO", id, err)
	} else if existCourseLearningOutcome == nil {
		return errs.New(errs.ErrCLONotFound, "cannot get courseLearningOutcome id %s to link subPLO", id)
	}

	course, err := u.courseUseCase.GetById(existCourseLearningOutcome.CourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get course id %s while linking clo and program outcome", existCourseLearningOutcome.CourseId, err)
	} else if course == nil {
		return errs.New(errs.ErrCourseNotFound, "course id %s not found while linking clo and program outcome", existCourseLearningOutcome.CourseId)
	}

	allPLOs, err := u.programmeUseCase.GetAllPLO(course.ProgrammeId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get program learning outcome while linking clo and subPLO", err)
	}

	allSPLOIds := []string{}
	for _, plo := range allPLOs {
		for _, splo := range plo.SubProgramLearningOutcomes {
			allSPLOIds = append(allSPLOIds, splo.Id)
		}
	}

	useablePLOIds := slice.Intersection(subProgramLearningOutcomeIds, allSPLOIds)
	if len(useablePLOIds) != len(subProgramLearningOutcomeIds) {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist sub plo %v", slice.Subtraction(subProgramLearningOutcomeIds, useablePLOIds))
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

	course, err := u.courseUseCase.GetById(existCourseLearningOutcome.CourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get course id %s while linking clo and program outcome", existCourseLearningOutcome.CourseId, err)
	} else if course == nil {
		return errs.New(errs.ErrCourseNotFound, "course id %s not found while linking clo and program outcome", existCourseLearningOutcome.CourseId)
	}

	allSOs, err := u.programmeUseCase.GetAllSO(course.ProgrammeId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get student outcome while linking clo and sub student outcome", err)
	}

	allSSOIds := []string{}
	for _, so := range allSOs {
		for _, sso := range so.SubStudentOutcomes {
			allSSOIds = append(allSSOIds, sso.Id)
		}
	}

	useableSOIds := slice.Intersection(subStudentOutcomeIds, allSSOIds)
	if len(useableSOIds) != len(subStudentOutcomeIds) {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist sub so %v", slice.Subtraction(subStudentOutcomeIds, useableSOIds))
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

	err = u.courseLearningOutcomeRepo.Update(id, &entity.CourseLearningOutcome{
		Code:                                payload.Code,
		Description:                         payload.Description,
		ExpectedPassingAssignmentPercentage: payload.ExpectedPassingAssignmentPercentage,
		ExpectedPassingStudentPercentage:    payload.ExpectedPassingStudentPercentage,
		Status:                              payload.Status,
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

func (u courseLearningOutcomeUseCase) DeleteLinkProgramOutcome(id string, programOutcomeId string) error {
	err := u.courseLearningOutcomeRepo.DeleteLinkProgramOutcome(id, programOutcomeId)
	if err != nil {
		return errs.New(errs.ErrUnLinkSubPLO, "cannot delete link CLO and program outcome", err)
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
