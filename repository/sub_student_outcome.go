package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
)

func (r StudentOutcomeRepositoryGorm) CreateSubSO(subStudentOutcome []*entity.SubStudentOutcome) error {
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
