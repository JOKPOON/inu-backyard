package usecase

import (
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	slice "github.com/team-inu/inu-backyard/internal/utils/slice"
)

type DepartmentUseCase struct {
	DepartmentRepo entity.DepartmentRepository
}

func NewDepartmentUseCase(departmentRepository entity.DepartmentRepository) entity.DepartmentUseCase {
	return &DepartmentUseCase{DepartmentRepo: departmentRepository}
}

func (u DepartmentUseCase) Create(department *entity.Department) error {
	err := u.DepartmentRepo.Create(department)

	if err != nil {
		return errs.New(errs.ErrCreateDepartment, "cannot create department", err)
	}

	return nil
}

func (u DepartmentUseCase) Delete(name string) error {
	err := u.DepartmentRepo.Delete(name)

	if err != nil {
		return errs.New(errs.ErrDeleteDepartment, "cannot delete department by name %s", name, err)
	}

	return nil
}

func (u DepartmentUseCase) GetAll() ([]entity.Department, error) {
	departments, err := u.DepartmentRepo.GetAll()

	if err != nil {
		return nil, errs.New(errs.ErrQueryDepartment, "cannot get all departments", err)
	}

	return departments, nil
}

func (u DepartmentUseCase) GetById(id string) (*entity.Department, error) {
	department, err := u.DepartmentRepo.GetById(id)

	if err != nil {
		return nil, errs.New(errs.ErrQueryDepartment, "cannot get department by name", err)
	}

	return department, nil
}

func (u DepartmentUseCase) Update(department *entity.Department) error {
	err := u.DepartmentRepo.Update(department)

	if err != nil {
		return errs.New(errs.ErrUpdateDepartment, "cannot update department", err)
	}

	return nil
}

func (u DepartmentUseCase) FilterNonExisted(names []string) ([]string, error) {
	existedNames, err := u.DepartmentRepo.FilterExisted(names)
	if err != nil {
		return nil, errs.New(errs.ErrQueryDepartment, "cannot query departments", err)
	}

	nonExistedIds := slice.Subtraction(names, existedNames)

	return nonExistedIds, nil
}
