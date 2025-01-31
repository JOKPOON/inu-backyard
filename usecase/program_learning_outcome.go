package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

type programLearningOutcomeUseCase struct {
	programLearningOutcomeRepo entity.ProgramLearningOutcomeRepository
	programmeUseCase           entity.ProgrammeUseCase
}

func NewProgramLearningOutcomeUseCase(
	programLearningOutcomeRepo entity.ProgramLearningOutcomeRepository,
	programmeUseCase entity.ProgrammeUseCase,
) entity.ProgramLearningOutcomeUseCase {
	return &programLearningOutcomeUseCase{
		programLearningOutcomeRepo: programLearningOutcomeRepo,
		programmeUseCase:           programmeUseCase,
	}
}

func (u programLearningOutcomeUseCase) GetAll() ([]entity.ProgramLearningOutcome, error) {
	plos, err := u.programLearningOutcomeRepo.GetAll()
	if err != nil {
		return nil, errs.New(errs.ErrQueryPLO, "cannot get all PLOs", err)
	}

	return plos, nil
}

func (u programLearningOutcomeUseCase) GetById(id string) (*entity.ProgramLearningOutcome, error) {
	plo, err := u.programLearningOutcomeRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQueryPLO, "cannot get PLO by id %s", id, err)
	}

	return plo, nil
}

func (u programLearningOutcomeUseCase) Create(payload []entity.CreateProgramLearningOutcome) error {
	plos := make([]entity.ProgramLearningOutcome, 0, len(payload))
	subPlos := make([]entity.SubProgramLearningOutcome, 0)

	for _, plo := range payload {
		id := ulid.Make().String()

		plos = append(plos, entity.ProgramLearningOutcome{
			Id:              id,
			Code:            plo.Code,
			DescriptionThai: plo.DescriptionThai,
			DescriptionEng:  plo.DescriptionEng,
		})

		for _, subPlo := range plo.SubProgramLearningOutcomes {
			subPlos = append(subPlos, entity.SubProgramLearningOutcome{
				Id:                       ulid.Make().String(),
				Code:                     subPlo.Code,
				DescriptionThai:          subPlo.DescriptionThai,
				DescriptionEng:           subPlo.DescriptionEng,
				ProgramLearningOutcomeId: id,
			})
		}
	}

	err := u.programLearningOutcomeRepo.CreateMany(plos)
	if err != nil {
		return errs.New(errs.ErrCreatePLO, "cannot create PLO", err)
	}

	if len(subPlos) > 0 {
		err = u.programLearningOutcomeRepo.CreateSubPLO(subPlos)
		if err != nil {
			return errs.New(errs.ErrCreatePLO, "cannot create sub plo", err)
		}
	}

	return nil
}

func (u programLearningOutcomeUseCase) Update(id string, programLearningOutcome *entity.ProgramLearningOutcome) error {
	nonExistedPloIds, err := u.FilterNonExisted([]string{id})
	if err != nil {
		return errs.New(errs.SameCode, "cannot get programLearningOutcome id %s to update")
	} else if len(nonExistedPloIds) > 0 {
		return errs.New(errs.ErrCreateSubPLO, "plo id not existed while updating plo %v", nonExistedPloIds)
	}

	err = u.programLearningOutcomeRepo.Update(id, programLearningOutcome)
	if err != nil {
		return errs.New(errs.ErrUpdatePLO, "cannot update programLearningOutcome by id %s", programLearningOutcome.Id, err)
	}

	return nil
}

func (u programLearningOutcomeUseCase) Delete(id string) error {
	nonExistedPloIds, err := u.FilterNonExisted([]string{id})
	if err != nil {
		return errs.New(errs.SameCode, "cannot get programLearningOutcome id %s to delete")
	} else if len(nonExistedPloIds) > 0 {
		return errs.New(errs.ErrCreateSubPLO, "plo id not existed while deleting plo %v", nonExistedPloIds)
	}

	splos, err := u.GetSubPloByPloId(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get subProgramLearningOutcome related to this programLearningOutcome")
	} else if len(splos) > 0 {
		return errs.New(errs.ErrCreateSubPLO, "splo related to this plo still exist %v", splos[0].Id)
	}

	err = u.programLearningOutcomeRepo.Delete(id)
	if err != nil {
		return errs.New(errs.ErrDeletePLO, "cannot delete PLO", err)
	}

	return nil
}

func (u programLearningOutcomeUseCase) FilterNonExisted(ids []string) ([]string, error) {
	existedIds, err := u.programLearningOutcomeRepo.FilterExisted(ids)
	if err != nil {
		return nil, errs.New(errs.ErrQueryPLO, "cannot query plo", err)
	}

	nonExistedIds := slice.Subtraction(ids, existedIds)

	return nonExistedIds, nil
}
