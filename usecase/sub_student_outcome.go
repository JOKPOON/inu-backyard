package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

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

	err = u.studentOutcomeRepo.CreateSubSO(subSos)
	if err != nil {
		return errs.New(errs.ErrCreateSubPLO, "cannot create sub plo", err)
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
