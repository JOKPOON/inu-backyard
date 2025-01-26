package usecase

import (
	"encoding/json"

	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

type programmeUseCase struct {
	programmeRepo entity.ProgrammeRepository
}

func NewProgrammeUseCase(programmeRepo entity.ProgrammeRepository) entity.ProgrammeUseCase {
	return &programmeUseCase{programmeRepo: programmeRepo}
}

func (u programmeUseCase) GetAll() ([]entity.Programme, error) {
	programme, err := u.programmeRepo.GetAll()
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get all programme", err)
	}

	return programme, nil
}

func (u programmeUseCase) GetByName(namesTH string, nameEN string) ([]entity.Programme, error) {
	programme, err := u.programmeRepo.GetByName(namesTH, nameEN)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get programme by name %s", namesTH+", "+nameEN, err)
	}

	return programme, nil
}

func (u programmeUseCase) GetByNameAndYear(nameTH string, nameEN string, year string) (*entity.Programme, error) {
	programme, err := u.programmeRepo.GetByNameAndYear(nameTH, nameEN, year)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get programme by name %s", nameTH+", "+nameEN, err)
	}

	return programme, nil
}

func (u programmeUseCase) GetById(id string) (*entity.Programme, error) {
	programme, err := u.programmeRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get programme by id %s", id, err)
	}

	return programme, nil
}

func (u programmeUseCase) Create(payload entity.CreateProgrammePayload) error {
	existProgramme, err := u.GetByNameAndYear(payload.NameTH, payload.NameEN, payload.Year)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get programme name %s to update", payload.NameTH+", "+payload.NameEN, err)
	} else if existProgramme != nil {
		return errs.New(errs.ErrDupName, "cannot create duplicate programme name %s", payload.NameTH+", "+payload.NameEN)
	}

	json, err := json.Marshal(payload.Structure)
	if err != nil {
		return errs.New(errs.ErrCreateProgramme, "cannot marshal programme structure", err)
	}

	programme := &entity.Programme{
		Id:            ulid.Make().String(),
		NameTH:        payload.NameTH,
		NameEN:        payload.NameEN,
		DegreeTH:      payload.DegreeTH,
		DegreeEN:      payload.DegreeEN,
		DegreeShortTH: payload.DegreeShortTH,
		DegreeShortEN: payload.DegreeShortEN,
		Year:          payload.Year,

		Structure: json,
	}

	err = u.programmeRepo.Create(programme)
	if err != nil {
		return errs.New(errs.ErrCreateProgramme, "cannot create programme", err)
	}

	return nil
}

func (u programmeUseCase) Update(id string, programme *entity.UpdateProgrammePayload) error {
	existProgramme, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get programme name %s to update", id, err)
	} else if existProgramme == nil {
		return errs.New(errs.ErrProgrammeNotFound, "cannot get programme name %s to update", id)
	}

	err = u.programmeRepo.Update(id, &entity.Programme{
		NameTH:        programme.NameTH,
		NameEN:        programme.NameEN,
		DegreeTH:      programme.DegreeTH,
		DegreeEN:      programme.DegreeEN,
		DegreeShortTH: programme.DegreeShortTH,
		DegreeShortEN: programme.DegreeShortEN,
		Year:          programme.Year,

		Structure: existProgramme.Structure,
	})
	if err != nil {
		return errs.New(errs.ErrUpdateProgramme, "cannot update programme by id %s", id, err)
	}

	return nil
}

func (u programmeUseCase) Delete(id string) error {
	programme, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get programme id %s name delete", id, err)
	} else if programme == nil {
		return errs.New(errs.ErrProgrammeNotFound, "cannot get programme name %s to delete", id)
	}

	err = u.programmeRepo.Delete(id)

	if err != nil {
		return errs.New(errs.ErrDeleteProgramme, "cannot delete programme by name %s", id, err)
	}

	return nil
}

func (u programmeUseCase) FilterNonExisted(namesTH []string, namesEN []string) ([]string, error) {
	existedNames, err := u.programmeRepo.FilterExisted(namesTH, namesEN)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot query programmes", err)
	}

	nonExistedIds := slice.Subtraction(namesTH, existedNames)

	return nonExistedIds, nil
}
