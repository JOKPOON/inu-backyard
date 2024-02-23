package request

import "github.com/team-inu/inu-backyard/entity"

type CreateScoreRequestPayload struct {
	StudentId    string  `json:"studentId" validate:"required"`
	Score        float64 `json:"score" validate:"required"`
	LecturerId   string  `json:"lecturerId" validate:"required"`
	AssignmentId string  `json:"assignmentId" validate:"required"`
}

type BulkCreateScoreRequestPayload struct {
	StudentScores []entity.StudentScore `json:"studentScores" validate:"dive"`
	LecturerId    string                `json:"lecturerId" validate:"required"`
	AssignmentId  string                `json:"assignmentId" validate:"required"`
}

type UpdateScoreRequestPayload struct {
	Score float64 `json:"score" validate:"required"`
}
