package repository

import (
	"errors"
	"fmt"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

func (r assignmentRepositoryGorm) GetAllGroup() ([]entity.AssignmentGroup, error) {
	var assignmentGroup []entity.AssignmentGroup

	err := r.gorm.Preload("Assignments").Find(&assignmentGroup).Error
	if err != nil {
		return nil, errors.New("cannot get assignment group")
	}

	return assignmentGroup, nil
}

func (r assignmentRepositoryGorm) GetGroupByQuery(query entity.AssignmentGroup, withAssignment bool) ([]entity.AssignmentGroup, error) {
	var assignmentGroup []entity.AssignmentGroup
	var err error

	if withAssignment {
		err = r.gorm.Preload("Assignments").Where(query).Find(&assignmentGroup).Error
	} else {
		err = r.gorm.Where(query).Find(&assignmentGroup).Error
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get assignment group by query: %w", err)
	}

	return assignmentGroup, nil
}

func (r assignmentRepositoryGorm) GetGroupByGroupId(assignmentGroupId string) (*entity.AssignmentGroup, error) {
	var assignmentGroup entity.AssignmentGroup
	err := r.gorm.Where("id = ?", assignmentGroupId).First(&assignmentGroup).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query to get assignment group by id: %w", err)
	}

	return &assignmentGroup, nil
}

func (r assignmentRepositoryGorm) CreateGroup(assignmentGroup *entity.AssignmentGroup) error {
	err := r.gorm.Create(&assignmentGroup).Error
	if err != nil {
		return fmt.Errorf("cannot create assignment group: %w", err)
	}

	return nil
}

func (r assignmentRepositoryGorm) UpdateGroup(assignmentGroupId string, assignmentGroup *entity.AssignmentGroup) error {
	err := r.gorm.Model(&entity.AssignmentGroup{}).Where("id = ?", assignmentGroupId).Updates(assignmentGroup).Error
	if err != nil {
		return fmt.Errorf("cannot update assignment group: %w", err)
	}

	return nil
}

func (r assignmentRepositoryGorm) DeleteGroup(assignmentGroupId string) error {
	err := r.gorm.Delete(&entity.AssignmentGroup{Id: assignmentGroupId}).Error

	if err != nil {
		return fmt.Errorf("cannot delete assignment group: %w", err)
	}

	return nil
}
