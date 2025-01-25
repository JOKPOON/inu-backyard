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

func (u programmeUseCase) GetByName(name string) ([]entity.Programme, error) {
	programme, err := u.programmeRepo.GetByName(name)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get programme by name %s", name, err)
	}

	return programme, nil
}

func (u programmeUseCase) GetByNameAndYear(name string, year string) (*entity.Programme, error) {
	programme, err := u.programmeRepo.GetByNameAndYear(name, year)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get programme by name %s", name, err)
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
	existProgramme, err := u.GetByNameAndYear(payload.Name, payload.Year)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get programme name %s to update", payload.Name+", "+payload.Year, err)
	} else if existProgramme != nil {
		return errs.New(errs.ErrDupName, "cannot create duplicate programme name %s", payload.Name+", "+payload.Year)
	}

	json, err := json.Marshal(payload.Structure)
	if err != nil {
		return errs.New(errs.ErrCreateProgramme, "cannot marshal programme structure", err)
	}

	programme := &entity.Programme{
		Id:        ulid.Make().String(),
		Name:      payload.Name,
		Year:      payload.Year,
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

	err = u.programmeRepo.Update(id, &entity.Programme{Name: programme.Name})
	if err != nil {
		return errs.New(errs.ErrUpdateProgramme, "cannot update programme by id %s", programme.Name, err)
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

func (u programmeUseCase) FilterNonExisted(names []string) ([]string, error) {
	existedNames, err := u.programmeRepo.FilterExisted(names)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot query programmes", err)
	}

	nonExistedIds := slice.Subtraction(names, existedNames)

	return nonExistedIds, nil
}
