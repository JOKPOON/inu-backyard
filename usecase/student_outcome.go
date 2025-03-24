package usecase

import (
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

type StudentOutcomeUsecase struct {
	studentOutcomeRepo entity.StudentOutcomeRepository
	programmeUseCase   entity.ProgrammeUseCase
}

func NewStudentOutcomeUseCase(studentOutcomeRepo entity.StudentOutcomeRepository, programmeUseCase entity.ProgrammeUseCase) entity.StudentOutcomeUseCase {
	return &StudentOutcomeUsecase{studentOutcomeRepo: studentOutcomeRepo, programmeUseCase: programmeUseCase}
}

func (u StudentOutcomeUsecase) GetAll(programId string) ([]entity.StudentOutcome, error) {
	plos, err := u.studentOutcomeRepo.GetAll(programId)
	if err != nil {
		return nil, errs.New(errs.ErrQuerySO, "cannot get all SOs", err)
	}

	return plos, nil
}

func (u StudentOutcomeUsecase) GetById(id string) (*entity.StudentOutcome, error) {
	plo, err := u.studentOutcomeRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQuerySubSO, "cannot get student outcome by id %s", id, err)
	}

	return plo, nil
}

func (u StudentOutcomeUsecase) Create(payload []entity.CreateStudentOutcome) error {
	sos := []*entity.StudentOutcome{}
	ssos := []*entity.SubStudentOutcome{}
	for _, so := range payload {
		id := ulid.Make().String()
		sos = append(sos, &entity.StudentOutcome{
			Id:              id,
			Code:            so.Code,
			DescriptionThai: so.DescriptionThai,
			DescriptionEng:  so.DescriptionEng,
			ProgramId:       so.ProgramId,
		})

		for _, sso := range so.SubStudentOutcomes {
			ssos = append(ssos, &entity.SubStudentOutcome{
				Id:               ulid.Make().String(),
				Code:             sso.Code,
				DescriptionThai:  sso.DescriptionThai,
				DescriptionEng:   sso.DescriptionEng,
				StudentOutcomeId: id,
			})
		}

	}

	err := u.studentOutcomeRepo.CreateMany(sos)
	if err != nil {
		return err
	}

	if len(ssos) > 0 {
		err = u.studentOutcomeRepo.CreateManySubSO(ssos)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u StudentOutcomeUsecase) Update(id string, payload *entity.UpdateStudentOutcomePayload) error {
	existedSO, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get student outcome id %s to update", id, err)
	} else if existedSO == nil {
		return errs.New(errs.ErrSONotFound, "cannot get student outcome id %s to update", id)
	}

	err = u.studentOutcomeRepo.Update(id, &entity.StudentOutcome{
		Code:            payload.Code,
		DescriptionThai: payload.DescriptionThai,
		DescriptionEng:  payload.DescriptionEng,
		ProgramId:       payload.ProgramId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (u StudentOutcomeUsecase) Delete(id string) error {
	err := u.studentOutcomeRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (u StudentOutcomeUsecase) FilterNonExisted(ids []string) ([]string, error) {
	existedIds, err := u.studentOutcomeRepo.FilterExisted(ids)
	if err != nil {
		return nil, errs.New(errs.ErrQuerySO, "cannot query so", err)
	}

	nonExistedIds := slice.Subtraction(ids, existedIds)

	return nonExistedIds, nil
}
