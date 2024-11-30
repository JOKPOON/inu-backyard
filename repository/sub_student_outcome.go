package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

func (r StudentOutcomeRepositoryGorm) GetAllSubSO() ([]entity.SubStudentOutcome, error) {
	var ssos []entity.SubStudentOutcome
	err := r.gorm.Find(&ssos).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get subSOs: %w", err)
	}

	return ssos, err
}

func (r StudentOutcomeRepositoryGorm) GetSubSOById(id string) (*entity.SubStudentOutcome, error) {
	var sso entity.SubStudentOutcome
	err := r.gorm.Preload("StudentOutcome").Where("id = ?", id).First(&sso).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get subSO by id: %w", err)
	}

	return &sso, nil
}

func (r StudentOutcomeRepositoryGorm) UpdateSubSO(id string, subSO *entity.SubStudentOutcome) error {
	tx := r.gorm.Model(&entity.SubStudentOutcome{}).
		Where("id = ?", id).
		Updates(subSO)

	if tx.Error != nil {
		return fmt.Errorf("cannot update subSO: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("subSO not found")
	}

	//TODO: abet
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r StudentOutcomeRepositoryGorm) DeleteSubSO(id string) error {
	tx := r.gorm.Delete(&entity.SubStudentOutcome{Id: id})

	if tx.Error != nil {
		return fmt.Errorf("cannot delete subSO: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("subSO not found")
	}

	//TODO: abet
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r StudentOutcomeRepositoryGorm) CreateSubSO(subStudentOutcome *entity.SubStudentOutcome) error {
	err := r.gorm.Create(&subStudentOutcome).Error
	if err != nil {
		return fmt.Errorf("cannot create subStudentOutcome: %w", err)
	}
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r StudentOutcomeRepositoryGorm) CreateManySubSO(subStudentOutcome []*entity.SubStudentOutcome) error {
	err := r.gorm.Create(&subStudentOutcome).Error
	if err != nil {
		return fmt.Errorf("cannot create subStudentOutcome: %w", err)
	}
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r StudentOutcomeRepositoryGorm) FilterExistedSubSO(ids []string) ([]string, error) {
	var existedIds []string

	err := r.gorm.Model(&entity.SubStudentOutcome{}).
		Where("id IN ?", ids).
		Pluck("id", &existedIds).
		Error
	if err != nil {
		return nil, fmt.Errorf("cannot filter existed subSOs: %w", err)
	}

	return existedIds, nil
}
