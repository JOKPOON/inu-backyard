package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils"
)

type courseLearningOutcomeUseCase struct {
	courseLearningOutcomeRepo        entity.CourseLearningOutcomeRepository
	courseUseCase                    entity.CourseUseCase
	programOutcomeUseCase            entity.ProgramOutcomeUseCase
	subProgramLearningOutcomeUseCase entity.SubProgramLearningOutcomeUseCase
}

func NewCourseLearningOutcomeUseCase(
	courseLearningOutcomeRepo entity.CourseLearningOutcomeRepository,
	courseUseCase entity.CourseUseCase,
	programOutcomeUseCase entity.ProgramOutcomeUseCase,
	subProgramLearningOutcomeUseCase entity.SubProgramLearningOutcomeUseCase,
) entity.CourseLearningOutcomeUseCase {
	return &courseLearningOutcomeUseCase{
		courseLearningOutcomeRepo:        courseLearningOutcomeRepo,
		courseUseCase:                    courseUseCase,
		programOutcomeUseCase:            programOutcomeUseCase,
		subProgramLearningOutcomeUseCase: subProgramLearningOutcomeUseCase,
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

func (u courseLearningOutcomeUseCase) GetByCourseId(courseId string) ([]entity.CourseLearningOutcome, error) {
	clo, err := u.courseLearningOutcomeRepo.GetByCourseId(courseId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryCLO, "cannot get CLO by course id %s", courseId, err)
	}

	return clo, nil
}

func (u courseLearningOutcomeUseCase) Create(dto entity.CreateCourseLearningOutcomeDto) error {
	if len(dto.SubProgramLearningOutcomeIds) == 0 {
		return errs.New(errs.ErrCreateCLO, "sub program learning outcome must not be empty when creating clo")
	}

	course, err := u.courseUseCase.GetById(dto.CourseId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get course id %s while creating clo", dto.CourseId, err)
	} else if course == nil {
		return errs.New(errs.ErrCourseNotFound, "course id %s not found while creating clo", dto.CourseId)
	}

	po, err := u.programOutcomeUseCase.GetById(dto.ProgramOutcomeId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get program outcome id %s while creating clo", dto.ProgramOutcomeId, err)
	} else if po == nil {
		return errs.New(errs.ErrCourseNotFound, "program outcome id %s not found while creating clo", dto.ProgramOutcomeId)
	}

	nonExistedSubPloIds, err := u.subProgramLearningOutcomeUseCase.FilterNonExisted(dto.SubProgramLearningOutcomeIds)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get non existed sub plo ids while creating clo")
	} else if len(nonExistedSubPloIds) != 0 {
		return errs.New(errs.ErrCreateEnrollment, "there are non exist sub plo")
	}

	subPlos := []*entity.SubProgramLearningOutcome{}
	for _, ploId := range dto.SubProgramLearningOutcomeIds {
		subPlos = append(subPlos, &entity.SubProgramLearningOutcome{
			Id: ploId,
		})
	}

	clo := entity.CourseLearningOutcome{
		Id:                                  ulid.Make().String(),
		Code:                                dto.Code,
		Description:                         dto.Description,
		Status:                              dto.Status,
		ExpectedPassingAssignmentPercentage: dto.ExpectedPassingAssignmentPercentage,
		ExpectedScorePercentage:             dto.ExpectedScorePercentage,
		ExpectedPassingStudentPercentage:    dto.ExpectedPassingStudentPercentage,
		ProgramOutcomeId:                    dto.ProgramOutcomeId,
		CourseId:                            dto.CourseId,
		SubProgramLearningOutcomes:          subPlos,
	}

	err = u.courseLearningOutcomeRepo.Create(&clo)
	if err != nil {
		return errs.New(errs.ErrCreateCLO, "cannot create CLO", err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) Update(id string, courseLearningOutcome *entity.CourseLearningOutcome) error {
	existCourseLearningOutcome, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get courseLearningOutcome id %s to update", id, err)
	} else if existCourseLearningOutcome == nil {
		return errs.New(errs.ErrCLONotFound, "cannot get courseLearningOutcome id %s to update", id)
	}

	err = u.courseLearningOutcomeRepo.Update(id, courseLearningOutcome)
	if err != nil {
		return errs.New(errs.ErrUpdateCLO, "cannot update courseLearningOutcome by id %s", courseLearningOutcome.Id, err)
	}

	return nil
}

func (u courseLearningOutcomeUseCase) Delete(id string) error {
	err := u.courseLearningOutcomeRepo.Delete(id)
	if err != nil {
		return errs.New(errs.ErrDeleteCLO, "cannot delete CLO", err)
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
