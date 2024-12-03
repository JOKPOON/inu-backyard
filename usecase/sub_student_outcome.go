package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

func (u StudentOutcomeUsecase) GetAllSubSO() ([]entity.SubStudentOutcome, error) {
	ssos, err := u.studentOutcomeRepo.GetAllSubSO()
	if err != nil {
		return nil, errs.New(errs.ErrQuerySubPLO, "cannot get all sub sos", err)
	}

	return ssos, nil
}

// func (u StudentOutcomeUsecase) GetSubPloByPloId(ploId string) ([]entity.SubStudentOutcome, error) {
// 	splos, err := u.programLearningOutcomeRepo.GetSubPloByPloId(ploId)
// 	if err != nil {
// 		return nil, errs.New(errs.ErrQuerySubPLO, "cannot get sub plos by plo id", err)
// 	}

// 	return splos, nil
// }

// func (u StudentOutcomeUsecase) GetSubPloByCode(code string, programme string, year int) (*entity.SubStudentOutcome, error) {
// 	splos, err := u.programLearningOutcomeRepo.GetSubPloByCode(code, programme, year)
// 	if err != nil {
// 		return nil, errs.New(errs.ErrQuerySubPLO, "cannot get sub plos by plo code", err)
// 	}

// 	return splos, nil
// }

func (u StudentOutcomeUsecase) GetSubSOById(id string) (*entity.SubStudentOutcome, error) {
	sso, err := u.studentOutcomeRepo.GetSubSOById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQuerySubPLO, "cannot get sub so by id %s", id, err)
	}

	return sso, nil
}

func (u StudentOutcomeUsecase) CreateSubSO(payload []entity.CreateSubStudentOutcome) error {
	soIds := []string{}
	for _, sso := range payload {
		soIds = append(soIds, sso.StudentOutcomeId)
	}

	soIds = slice.RemoveDuplicates(soIds)
	nonExistedSoIds, err := u.FilterNonExisted(soIds)
	if err != nil {
		return errs.New(errs.SameCode, "cannot find non existing so id while creating sub so")
	} else if len(nonExistedSoIds) > 0 {
		return errs.New(errs.ErrCreateSubPLO, "so ids not existed while creating sub so %v", nonExistedSoIds)
	}

	subSos := make([]*entity.SubStudentOutcome, 0, len(payload))
	for _, subPlo := range payload {
		subSos = append(subSos, &entity.SubStudentOutcome{
			Id:               ulid.Make().String(),
			Code:             subPlo.Code,
			DescriptionThai:  subPlo.DescriptionThai,
			DescriptionEng:   subPlo.DescriptionEng,
			StudentOutcomeId: subPlo.StudentOutcomeId,
		})
	}

	err = u.studentOutcomeRepo.CreateManySubSO(subSos)
	if err != nil {
		return errs.New(errs.ErrCreateSubPLO, "cannot create sub so", err)
	}

	return nil
}

func (u StudentOutcomeUsecase) UpdateSubSO(id string, subStudentOutcome *entity.UpdateSubStudentOutcomePayload) error {
	existSubSO, err := u.GetSubSOById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get subSO id %s to update", id, err)
	} else if existSubSO == nil {
		return errs.New(errs.ErrSubPLONotFound, "cannot get subSO id %s to update", id)
	}

	nonExistedSoIds, err := u.FilterNonExisted([]string{subStudentOutcome.StudentOutcomeId})
	if err != nil {
		return errs.New(errs.SameCode, "cannot find non existing so id while updating sub so")
	} else if len(nonExistedSoIds) > 0 {
		return errs.New(errs.ErrCreateSubPLO, "so ids not existed while updating sub so %v", nonExistedSoIds)
	}

	err = u.studentOutcomeRepo.UpdateSubSO(id, &entity.SubStudentOutcome{
		Code:             subStudentOutcome.Code,
		DescriptionThai:  subStudentOutcome.DescriptionThai,
		DescriptionEng:   subStudentOutcome.DescriptionEng,
		StudentOutcomeId: subStudentOutcome.StudentOutcomeId,
	})
	if err != nil {
		return errs.New(errs.ErrUpdateSubPLO, err.Error(), err)
	}

	return nil
}

func (u StudentOutcomeUsecase) DeleteSubSO(id string) error {
	existSubSO, err := u.GetSubSOById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get subSO id %s to delete", id, err)
	} else if existSubSO == nil {
		return errs.New(errs.ErrSubPLONotFound, "cannot get subSO id %s to delete", id)
	}

	err = u.studentOutcomeRepo.DeleteSubSO(id)
	if err != nil {
		return errs.New(errs.ErrDeleteSubPLO, "cannot delete sub so", err)
	}

	return nil
}

func (u StudentOutcomeUsecase) FilterNonExistedSubSO(ids []string) ([]string, error) {
	existedIds, err := u.studentOutcomeRepo.FilterExistedSubSO(ids)
	if err != nil {
		return nil, errs.New(errs.ErrQuerySubSO, "cannot query sub so", err)
	}

	nonExistedIds := slice.Subtraction(ids, existedIds)

	return nonExistedIds, nil
}
