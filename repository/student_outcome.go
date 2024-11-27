package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

type StudentOutcomeRepositoryGorm struct {
	gorm *gorm.DB
}

func NewStudentOutcomeRepositoryGorm(gorm *gorm.DB) entity.StudentOutcomeRepository {
	return &StudentOutcomeRepositoryGorm{gorm: gorm}
}

func (r StudentOutcomeRepositoryGorm) GetAll() ([]entity.StudentOutcome, error) {
	var sos []entity.StudentOutcome
	err := r.gorm.Preload("SubStudentOutcomes").Find(&sos).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return sos, nil
}

func (r StudentOutcomeRepositoryGorm) GetById(id string) (*entity.StudentOutcome, error) {
	var so entity.StudentOutcome
	err := r.gorm.Preload("SubStudentOutcomes").Where("id =?", id).First(&so).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &so, nil
}

func (r StudentOutcomeRepositoryGorm) Create(so *entity.StudentOutcome) error {
	err := r.gorm.Create(so).Error
	if err != nil {
		return fmt.Errorf("cannot create student_outcome: %w", err)
	}

	return nil
}

func (r StudentOutcomeRepositoryGorm) CreateMany(sos []*entity.StudentOutcome) error {
	err := r.gorm.Create(&sos).Error
	if err != nil {
		return fmt.Errorf("cannot create student_outcome: %w", err)
	}

	return nil
}

func (r StudentOutcomeRepositoryGorm) Update(id string, so *entity.StudentOutcome) error {
	err := r.gorm.Model(&entity.StudentOutcome{}).Where("id =?", id).Updates(so).Error
	if err != nil {
		return fmt.Errorf("cannot update student_outcome: %w", err)
	}

	return nil
}

func (r StudentOutcomeRepositoryGorm) Delete(id string) error {
	err := r.gorm.Delete(&entity.StudentOutcome{Id: id}).Error

	if err != nil {
		return fmt.Errorf("cannot delete student_outcome: %w", err)
	}

	return nil
}

func (r StudentOutcomeRepositoryGorm) FilterExisted(ids []string) ([]string, error) {
	var existedIds []string

	err := r.gorm.Raw("SELECT id FROM `student_outcome` WHERE id in ?", ids).Scan(&existedIds).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query student_outcome: %w", err)
	}

	return existedIds, nil
}
