package entity

type Survey struct {
	Id         string     `json:"id" gorm:"primaryKey;type:char(255)"`
	IsComplete bool       `json:"is_complete"`
	CreateAt   string     `json:"create_at" gorm:"type:timestamp"`
	Questions  []Question `json:"questions" gorm:"foreignKey:SurveyId"`
	CourseId   string     `json:"course_id" gorm:"index"`
}

type Question struct {
	Id       string `json:"id" gorm:"primaryKey;type:char(255)"`
	Question string `json:"question"`
	POId     string `json:"po_id"`
	PLOId    string `json:"plo_id"`
	SOId     string `json:"so_id"`
	SurveyId string `json:"survey_id" gorm:"index"` // Foreign key for Survey

	Scores []QScore `json:"q_scores" gorm:"foreignKey:QuestionId;constraint:OnDelete:CASCADE;"`
}

type QScore struct {
	Id         string  `json:"id" gorm:"primaryKey;type:char(255)"`
	Score      float64 `json:"score"`
	QuestionId string  `json:"question_id" gorm:"index"` // Foreign key for Question
}
