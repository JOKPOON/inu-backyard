package entity

import (
	"time"
)

type SurveyRepository interface {
	Create(survey *Survey) error
	AddQuestion(question *Question) error
	Delete(id string) error
	RemoveQuestion(id string) error

	GetQuestionById(id string) (*Question, error)
	GetQuestionsBySurveyId(surveyID string) ([]Question, error)
	GetAll() ([]Survey, error)
	GetById(id string) (*Survey, error)
	GetByCourseId(courseId string) (*Survey, error)
	UpdateQuestion(question *Question) error
	Update(survey *Survey) error

	GetSurveysWithCourseAndOutcomes() ([]SurveyWithCourseAndOutcomes, error)
}

type SurveyUseCase interface {
	Create(request *CreateSurveyRequest) error
	Delete(id string) error
	GetAll() ([]Survey, error)
	GetById(id string) (*Survey, error)
	GetByCourseId(courseId string) (*Survey, error)
	Update(id string, request *UpdateSurveyRequest) error

	CreateQuestion(surveyId string, request *CreateQuestionRequest) error
	DeleteQuestion(id string) error
	GetQuestionById(id string) (*Question, error)
	GetQuestionsBySurveyId(surveyId string) ([]Question, error)
	UpdateQuestion(id string, request *UpdateQuestionRequest) error

	GetSurveysWithCourseAndOutcomes() ([]SurveyWithCourseAndOutcomes, error)
}

type Survey struct {
	Id          string     `json:"id" gorm:"primaryKey;type:char(255)"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsComplete  bool       `json:"is_complete"`
	CreateAt    time.Time  `json:"create_at" gorm:"autoCreateTime"`
	Questions   []Question `json:"questions" gorm:"foreignKey:SurveyId"`
	CourseId    string     `json:"course_id" gorm:"index"`
}

type Question struct {
	Id          string `json:"id" gorm:"primaryKey;type:char(255)"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Question    string `json:"question"`
	POId        string `json:"po_id"`
	PLOId       string `json:"plo_id"`
	SOId        string `json:"so_id"`
	SurveyId    string `json:"survey_id" gorm:"index"`

	Scores []QScore `json:"q_scores" gorm:"foreignKey:QuestionId;constraint:OnDelete:CASCADE;"`
}

type QScore struct {
	Id         string  `json:"id" gorm:"primaryKey;type:char(255)"`
	Score      float64 `json:"score"`
	QuestionId string  `json:"question_id" gorm:"index"`
}

type CreateSurveyRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsComplete  bool       `json:"is_complete"`
	CourseId    string     `json:"course_id"`
	Questions   []Question `json:"questions"`
}

type UpdateSurveyRequest struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	IsComplete  bool       `json:"is_complete,omitempty"`
	CourseId    string     `json:"course_id,omitempty"`
	Questions   []Question `json:"questions,omitempty"`
}

type CreateQuestionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Question    string `json:"question"`
	POId        string `json:"po_id"`
	PLOId       string `json:"plo_id"`
	SOId        string `json:"so_id"`
}

type UpdateQuestionRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Question    string `json:"question,omitempty"`
	POId        string `json:"po_id,omitempty"`
	PLOId       string `json:"plo_id,omitempty"`
	SOId        string `json:"so_id,omitempty"`
}

type SurveyWithCourseAndOutcomes struct {
	SurveyId     string   `json:"survey_id"`
	SurveyTitle  string   `json:"survey_title"`
	Description  string   `json:"description"`
	IsComplete   bool     `json:"is_complete"`
	CourseId     string   `json:"course_id"`
	CourseName   string   `json:"course_name"`
	CourseCode   string   `json:"course_code"`
	AcademicYear string   `json:"academic_year"`
	POs          []string `json:"pos"`
	PLOs         []string `json:"plos"`
	SOs          []string `json:"sos"`
}
