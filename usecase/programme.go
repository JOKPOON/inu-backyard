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

func (u programmeUseCase) GetAllCourseOutcomeLinked(programmeId string) ([]entity.CourseOutcomes, error) {
	programme, err := u.GetById(programmeId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get programme id %s to get course outcome", programmeId, err)
	} else if programme == nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get programme id %s to get course outcome", programmeId, err)
	}

	resp, err := u.programmeRepo.GetAllCourseOutcomeLinked(programmeId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get course outcome by programme id %s", programmeId, err)
	}

	return resp, nil
}

func (u programmeUseCase) GetAllCourseLinkedPO(programmeIds []string) ([]entity.ProgrammeLinkedPO, error) {
	var programmeLinkedPOs []entity.ProgrammeLinkedPO
	for _, id := range programmeIds {
		programme, err := u.GetById(id)
		if err != nil {
			return nil, errs.New(errs.SameCode, "cannot get programme id %s to get course outcome", id, err)
		} else if programme == nil {
			return nil, errs.New(errs.ErrQueryProgramme, "cannot get programme id %s to get course outcome", id, err)
		}

		resp, err := u.programmeRepo.GetAllCourseLinkedPO(id)
		if err != nil {
			return nil, errs.New(errs.ErrQueryProgramme, "cannot get course outcome by programme id %s", id, err)
		}

		resp.ProgrammeName = programme.NameTH + ", " + programme.NameEN
		resp.ProgrammeYear = programme.Year

		programmeLinkedPOs = append(programmeLinkedPOs, *resp)
	}

	return programmeLinkedPOs, nil
}

func (u programmeUseCase) GetAllCourseLinkedPLO(programmeIds []string) ([]entity.ProgrammeLinkedPLO, error) {
	var programmeLinkedPLOs []entity.ProgrammeLinkedPLO
	for _, id := range programmeIds {
		programme, err := u.GetById(id)
		if err != nil {
			return nil, errs.New(errs.SameCode, "cannot get programme id %s to get course outcome", id, err)
		} else if programme == nil {
			return nil, errs.New(errs.ErrQueryProgramme, "cannot get programme id %s to get course outcome", id, err)
		}

		resp, err := u.programmeRepo.GetAllCourseLinkedPLO(id)
		if err != nil {
			return nil, errs.New(errs.ErrQueryProgramme, "cannot get course outcome by programme id %s", id, err)
		}

		resp.ProgrammeName = programme.NameTH + ", " + programme.NameEN
		resp.ProgrammeYear = programme.Year

		programmeLinkedPLOs = append(programmeLinkedPLOs, *resp)
	}

	return programmeLinkedPLOs, nil
}

func (u programmeUseCase) GetAllCourseLinkedSO(programmeId []string) ([]entity.ProgrammeLinkedSO, error) {
	var programmeLinkedSOs []entity.ProgrammeLinkedSO
	for _, id := range programmeId {
		programme, err := u.GetById(id)
		if err != nil {
			return nil, errs.New(errs.SameCode, "cannot get programme id %s to get course outcome", programmeId, err)
		} else if programme == nil {
			return nil, errs.New(errs.ErrQueryProgramme, "cannot get programme id %s to get course outcome", programmeId, err)
		}

		resp, err := u.programmeRepo.GetAllCourseLinkedSO(id)
		if err != nil {
			return nil, errs.New(errs.ErrQueryProgramme, "cannot get course outcome by programme id %s", programmeId, err)
		}
		//TODO:
		resp.ProgrammeName = programme.NameTH + ", " + programme.NameEN
		resp.ProgrammeYear = programme.Year

		programmeLinkedSOs = append(programmeLinkedSOs, *resp)
	}

	return programmeLinkedSOs, nil
}

func (u programmeUseCase) CreateLinkWithPO(programmeId string, poIds []string) error {
	programme, err := u.GetById(programmeId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get programme id %s to create link with po", programmeId, err)
	} else if programme == nil {
		return errs.New(errs.ErrQueryProgramme, "cannot get programme id %s to create link with po", programmeId, err)
	}

	existedPOIds, err := u.FilterExistedPO(programmeId, poIds)
	if err != nil {
		return errs.New(errs.ErrQueryProgramme, "cannot query po by programme id %s", programmeId, err)
	}

	nonExistedPOIds := slice.Subtraction(poIds, existedPOIds)

	for _, poId := range nonExistedPOIds {
		err = u.programmeRepo.CreateLinkWithPO(programmeId, poId)
		if err != nil {
			return errs.New(errs.ErrCreateProgramme, "cannot create link with po by programme id %s", programmeId, err)
		}
	}

	return nil
}

func (u programmeUseCase) CreateLinkWithPLO(programmeId string, ploIds []string) error {
	programme, err := u.GetById(programmeId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get programme id %s to create link with plo", programmeId, err)
	} else if programme == nil {
		return errs.New(errs.ErrQueryProgramme, "cannot get programme id %s to create link with plo", programmeId, err)
	}

	existedPLOIds, err := u.FilterExistedPLO(programmeId, ploIds)
	if err != nil {
		return errs.New(errs.ErrQueryProgramme, "cannot query plo by programme id %s", programmeId, err)
	}

	nonExistedPLOIds := slice.Subtraction(ploIds, existedPLOIds)

	for _, ploId := range nonExistedPLOIds {
		err = u.programmeRepo.CreateLinkWithPLO(programmeId, ploId)
		if err != nil {
			return errs.New(errs.ErrCreateProgramme, "cannot create link with plo by programme id %s", programmeId, err)
		}
	}

	return nil
}

func (u programmeUseCase) CreateLinkWithSO(programmeId string, soIds []string) error {
	programme, err := u.GetById(programmeId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get programme id %s to create link with so", programmeId, err)
	} else if programme == nil {
		return errs.New(errs.ErrQueryProgramme, "cannot get programme id %s to create link with so", programmeId, err)
	}

	existedSOIds, err := u.FilterExistedSO(programmeId, soIds)
	if err != nil {
		return errs.New(errs.ErrQueryProgramme, "cannot query so by programme id %s", programmeId, err)
	}

	nonExistedSOIds := slice.Subtraction(soIds, existedSOIds)

	for _, soId := range nonExistedSOIds {
		err = u.programmeRepo.CreateLinkWithSO(programmeId, soId)
		if err != nil {
			return errs.New(errs.ErrCreateProgramme, "cannot create link with so by programme id %s", programmeId, err)
		}
	}

	return nil
}

func (u programmeUseCase) FilterExistedPO(programmeId string, poIds []string) ([]string, error) {
	existedPOIds, err := u.programmeRepo.FilterExistedPO(programmeId, poIds)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot query po by programme id %s", programmeId, err)
	}

	return existedPOIds, nil
}

func (u programmeUseCase) FilterExistedPLO(programmeId string, ploIds []string) ([]string, error) {
	existedPLOIds, err := u.programmeRepo.FilterExistedPLO(programmeId, ploIds)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot query plo by programme id %s", programmeId, err)
	}

	return existedPLOIds, nil
}

func (u programmeUseCase) FilterExistedSO(programmeId string, soIds []string) ([]string, error) {
	existedSOIds, err := u.programmeRepo.FilterExistedSO(programmeId, soIds)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot query so by programme id %s", programmeId, err)
	}

	return existedSOIds, nil
}

func (u programmeUseCase) GetAllPO(programmeId string) ([]entity.ProgramOutcome, error) {
	pos, err := u.programmeRepo.GetAllPO(programmeId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get po by programme id %s", programmeId, err)
	}

	return pos, nil
}

func (u programmeUseCase) GetAllPLO(programmeId string) ([]entity.ProgramLearningOutcome, error) {
	plos, err := u.programmeRepo.GetAllPLO(programmeId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get plo by programme id %s", programmeId, err)
	}

	return plos, nil
}

func (u programmeUseCase) GetAllSO(programmeId string) ([]entity.StudentOutcome, error) {
	sos, err := u.programmeRepo.GetAllSO(programmeId)
	if err != nil {
		return nil, errs.New(errs.ErrQueryProgramme, "cannot get so by programme id %s", programmeId, err)
	}

	return sos, nil
}
