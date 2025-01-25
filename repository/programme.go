package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

type programmeRepositoryGorm struct {
	gorm *gorm.DB
}

func NewProgrammeRepositoryGorm(gorm *gorm.DB) entity.ProgrammeRepository {
	return &programmeRepositoryGorm{gorm}
}

func (r programmeRepositoryGorm) GetAll() ([]entity.Programme, error) {
	var programs []entity.Programme

	err := r.gorm.Find(&programs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get programs: %w", err)
	}

	return programs, nil
}

func (r programmeRepositoryGorm) GetById(id string) (*entity.Programme, error) {
	var programme *entity.Programme

	err := r.gorm.Where("id = ?", id).First(&programme).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get programme by id: %w", err)
	}

	return programme, nil
}

func (r programmeRepositoryGorm) GetByName(name string) ([]entity.Programme, error) {
	var programme []entity.Programme

	err := r.gorm.Find(&programme, "name = ?", name).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get programme by id: %w", err)
	}

	return programme, nil
}

func (r programmeRepositoryGorm) GetByNameAndYear(name string, year string) (*entity.Programme, error) {
	var programme *entity.Programme

	err := r.gorm.Where("name = ? AND year = ?", name, year).First(&programme).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get programme by id: %w", err)
	}

	return programme, nil
}

func (r programmeRepositoryGorm) Create(programme *entity.Programme) error {
	err := r.gorm.Create(&programme).Error
	if err != nil {
		return fmt.Errorf("cannot create programme: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) Update(id string, programme *entity.Programme) error {
	err := r.gorm.Model(&entity.Programme{}).Where("id = ?", id).Updates(programme).Error
	if err != nil {
		return fmt.Errorf("cannot update programme: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) Delete(id string) error {
	err := r.gorm.Where("id = ?", id).Delete(&entity.Programme{}).Error
	if err != nil {
		return fmt.Errorf("cannot delete programme: %w", err)
	}

	return nil
}

func (r programmeRepositoryGorm) FilterExisted(names []string) ([]string, error) {
	var existedNames []string

	err := r.gorm.Raw("SELECT name FROM `programme` WHERE name in ?", names).Scan(&existedNames).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query programme: %w", err)
	}

	return existedNames, nil
}
