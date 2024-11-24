package entity

import "time"

type ProgramImprovement struct {
	Id              string    `json:"id" gorm:"type:char(255);primaryKey"`
	ProgrammeName   string    `json:"programme_name" gorm:"type:char(255);primaryKey"`
	IssueIdentified string    `json:"issue_identified" gorm:"type:text"`
	ActionTaken     string    `json:"action_taken" gorm:"type:text"`
	Result          string    `json:"result" gorm:"type:text"`
	Date            time.Time `json:"date"`

	Programme Programme `gorm:"foreignKey:ProgrammeName"`
}

type CreateProgramImprovementRequestPayload struct {
	ProgrammeName   string    `json:"programme_name" validate:"required"`
	IssueIdentified string    `json:"issue_identified" validate:"required"`
	ActionTaken     string    `json:"action_taken" validate:"required"`
	Result          string    `json:"result" validate:"required"`
	Date            time.Time `json:"date" validate:"required"`
}
