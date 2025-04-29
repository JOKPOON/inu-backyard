package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"gorm.io/gorm"
)

type FacultyRepositoryGorm struct {
	gorm *gorm.DB
}

func NewFacultyRepositoryGorm(gorm *gorm.DB) entity.FacultyRepository {
	return &FacultyRepositoryGorm{gorm: gorm}
}

func (r FacultyRepositoryGorm) Create(faculty *entity.Faculty) error {
	err := r.gorm.Create(&faculty).Error
	if err != nil {
		return errs.New(errs.ErrCreateFaculty, "cannot create faculty", err)
	}

	return nil
}

func (r FacultyRepositoryGorm) Delete(id string) error {
	err := r.gorm.Delete(&entity.Faculty{
		Id: id,
	}).Error
	if err != nil {
		return errs.New(errs.ErrDeleteFaculty, "cannot delete faculty by name %s", id, err)
	}

	return nil
}

func (r FacultyRepositoryGorm) GetAll() ([]entity.Faculty, error) {
	var faculties []entity.Faculty
	err := r.gorm.Preload("Departments").
		Preload("Departments.Programmes").Find(&faculties).Error
	if err != nil {
		return nil, errs.New(errs.ErrQueryFaculty, "cannot get all faculties", err)
	}

	return faculties, nil
}

func (r *FacultyRepositoryGorm) GetById(id string) (*entity.Faculty, error) {
	var faculty *entity.Faculty

	err := r.gorm.Where("id = ?", id).First(&faculty).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get faculty: %w", err)
	}

	return faculty, nil
}

func (r *FacultyRepositoryGorm) Update(faculty *entity.Faculty) error {
	//find old faculty by name
	var oldFaculty *entity.Faculty
	err := r.gorm.Where("id = ?", faculty.Id).First(&oldFaculty).Error
	if err != nil {
		return errs.New(errs.ErrQueryFaculty, "cannot get faculty", err)
	}

	//update old faculty with new name
	err = r.gorm.Model(&oldFaculty).Updates(faculty).Error
	if err != nil {
		return errs.New(errs.ErrUpdateFaculty, "cannot update faculty", err)
	}

	return nil

}
