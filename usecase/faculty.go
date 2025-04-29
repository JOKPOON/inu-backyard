package usecase

import (
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
)

type FacultyUseCase struct {
	facultyRepo entity.FacultyRepository
}

func NewFacultyUseCase(facultyRepo entity.FacultyRepository) entity.FacultyUseCase {
	return &FacultyUseCase{facultyRepo: facultyRepo}
}

func (u FacultyUseCase) Create(faculty *entity.Faculty) error {

	err := u.facultyRepo.Create(faculty)

	if err != nil {
		return errs.New(errs.ErrCreateFaculty, "cannot create faculty", err)
	}

	return nil
}

func (u FacultyUseCase) Delete(id string) error {
	err := u.facultyRepo.Delete(id)

	if err != nil {
		return errs.New(errs.ErrDeleteFaculty, "cannot delete faculty by name %s", id, err)
	}

	return nil
}

func (u FacultyUseCase) GetAll() ([]entity.Faculty, error) {
	faculties, err := u.facultyRepo.GetAll()

	if err != nil {
		return nil, errs.New(errs.ErrQueryFaculty, "cannot get all faculty", err)
	}

	return faculties, nil
}

func (u FacultyUseCase) GetById(id string) (*entity.Faculty, error) {
	faculty, err := u.facultyRepo.GetById(id)

	if err != nil {
		return nil, errs.New(errs.ErrQueryFaculty, "cannot get faculty by id %s", id, err)
	}

	return faculty, nil
}

func (u FacultyUseCase) Update(faculty *entity.Faculty) error {
	err := u.facultyRepo.Update(faculty)

	if err != nil {
		return errs.New(errs.ErrUpdateFaculty, "cannot update faculty", err)
	}

	return nil
}
