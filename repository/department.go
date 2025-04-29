package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

type DepartmentRepositoryGorm struct {
	gorm *gorm.DB
}

func NewDepartmentRepositoryGorm(gorm *gorm.DB) entity.DepartmentRepository {
	return &DepartmentRepositoryGorm{gorm: gorm}
}

func (r DepartmentRepositoryGorm) Create(department *entity.Department) error {
	err := r.gorm.Create(&department).Error
	if err != nil {
		return fmt.Errorf("cannot create department: %w", err)
	}

	return nil
}

func (r DepartmentRepositoryGorm) Delete(id string) error {
	err := r.gorm.Where("id = ?", id).Delete(&entity.Department{}).Error
	if err != nil {
		return fmt.Errorf("cannot delete department by name: %w", err)
	}

	return nil
}

func (r DepartmentRepositoryGorm) GetAll() ([]entity.Department, error) {
	var departments []entity.Department
	err := r.gorm.Find(&departments).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query to get department by name: %w", err)
	}

	return departments, nil
}

func (r *DepartmentRepositoryGorm) GetById(id string) (*entity.Department, error) {
	var department *entity.Department

	err := r.gorm.Where("id = ?", id).First(&department).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get department by name: %w", err)
	}

	return department, nil
}

func (r *DepartmentRepositoryGorm) Update(department *entity.Department) error {
	//find old department by name
	var oldDepartment *entity.Department
	err := r.gorm.Where("id = ?", department.Id).First(&oldDepartment).Error
	if err != nil {
		return fmt.Errorf("cannot get department while updating department: %w", err)
	}

	//update old department with new name
	err = r.gorm.Model(&oldDepartment).Updates(department).Error
	if err != nil {
		return fmt.Errorf("cannot update department by name: %w", err)
	}

	return nil
}

func (r *DepartmentRepositoryGorm) FilterExisted(id []string) ([]string, error) {
	var existedIds []string

	err := r.gorm.Raw("SELECT id FROM `department` WHERE id in ?", id).Scan(&existedIds).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query department: %w", err)
	}

	return existedIds, nil
}
