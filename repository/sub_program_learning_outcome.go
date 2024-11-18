package repository

import (
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

func (r programLearningOutcomeRepositoryGorm) GetAllSubPlo() ([]entity.SubProgramLearningOutcome, error) {
	var splos []entity.SubProgramLearningOutcome
	err := r.gorm.Find(&splos).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get subPLOs: %w", err)
	}

	return splos, err
}

func (r programLearningOutcomeRepositoryGorm) GetSubPloByPloId(ploId string) ([]entity.SubProgramLearningOutcome, error) {
	var splos []entity.SubProgramLearningOutcome
	err := r.gorm.Where("program_learning_outcome_id = ?", ploId).Find(&splos).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get subPLOs with plo id: %w", err)
	}

	return splos, err
}

func (r programLearningOutcomeRepositoryGorm) GetSubPloByCode(code string, programme string, year int) (*entity.SubProgramLearningOutcome, error) {
	var splo entity.SubProgramLearningOutcome
	err := r.gorm.Model(&splo).
		Select("sub_program_learning_outcome.*").
		Joins("LEFT JOIN program_learning_outcome on sub_program_learning_outcome.program_learning_outcome_id = program_learning_outcome.id").
		Where("sub_program_learning_outcome.code = ? AND program_learning_outcome.program_year = ? AND program_learning_outcome.programme_name = ?", code, year, programme).
		First(&splo).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get subPLOs with plo id: %w", err)
	}

	return &splo, err
}

func (r programLearningOutcomeRepositoryGorm) GetSubPLO(id string) (*entity.SubProgramLearningOutcome, error) {
	var splo entity.SubProgramLearningOutcome
	err := r.gorm.Where("id = ?", id).First(&splo).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get subPLO by id: %w", err)
	}

	return &splo, nil
}

func (r programLearningOutcomeRepositoryGorm) CreateSubPLO(subProgramLearningOutcome []entity.SubProgramLearningOutcome) error {
	err := r.gorm.Create(&subProgramLearningOutcome).Error
	if err != nil {
		return fmt.Errorf("cannot create subProgramLearningOutcome: %w", err)
	}
	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r programLearningOutcomeRepositoryGorm) UpdateSubPLO(id string, subProgramLearningOutcome *entity.SubProgramLearningOutcome) error {
	tx := r.gorm.Model(&entity.SubProgramLearningOutcome{}).
		Where("id = ?", id).
		Updates(subProgramLearningOutcome)

	if tx.Error != nil {
		return fmt.Errorf("cannot update subProgramLearningOutcome: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("subProgramLearningOutcome not found")
	}

	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r programLearningOutcomeRepositoryGorm) DeleteSubPLO(id string) error {
	tx := r.gorm.Delete(&entity.SubProgramLearningOutcome{Id: id})

	if tx.Error != nil {
		return fmt.Errorf("cannot delete subProgramLearningOutcome: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("subProgramLearningOutcome not found")
	}

	go cacheOutcomes(r.gorm, TabeeSelectorAllPloCourses)
	go cacheOutcomes(r.gorm, TabeeSelectorAllPoCourses)

	return nil
}

func (r programLearningOutcomeRepositoryGorm) FilterExistedSubPLO(ids []string) ([]string, error) {
	var existedIds []string

	err := r.gorm.Model(&entity.SubProgramLearningOutcome{}).
		Where("id IN ?", ids).
		Pluck("id", &existedIds).
		Error
	if err != nil {
		return nil, fmt.Errorf("cannot filter existed subPLOs: %w", err)
	}

	return existedIds, nil
}
