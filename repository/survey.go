package repository

import (
	"fmt"
	"strings"

	"github.com/team-inu/inu-backyard/entity"
	"gorm.io/gorm"
)

type SurveyRepositoryGorm struct {
	gorm *gorm.DB
}

func NewSurveyRepositoryGorm(gorm *gorm.DB) entity.SurveyRepository {
	return &SurveyRepositoryGorm{gorm: gorm}
}

func (r *SurveyRepositoryGorm) Create(survey *entity.Survey) error {
	err := r.gorm.Create(&survey).Error
	if err != nil {
		return fmt.Errorf("cannot create survey: %w", err)
	}

	return nil
}

func (r *SurveyRepositoryGorm) Delete(id string) error {
	err := r.gorm.Where("id = ?", id).Delete(&entity.Survey{}).Error
	if err != nil {
		return fmt.Errorf("cannot delete survey by id: %w", err)
	}

	return nil
}

func (r *SurveyRepositoryGorm) GetAll() ([]entity.Survey, error) {
	var surveys []entity.Survey
	err := r.gorm.Preload("Questions.Scores").Find(&surveys).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query surveys: %w", err)
	}

	return surveys, nil
}

func (r *SurveyRepositoryGorm) GetById(id string) (*entity.Survey, error) {
	var survey entity.Survey
	err := r.gorm.Preload("Questions.Scores").Where("id = ?", id).First(&survey).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query survey by id: %w", err)
	}

	return &survey, nil
}

func (r *SurveyRepositoryGorm) GetByCourseId(courseID string) (*entity.Survey, error) {
	var survey entity.Survey
	err := r.gorm.Preload("Questions.Scores").Where("course_id = ?", courseID).First(&survey).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot query survey by course id: %w", err)
	}

	return &survey, nil
}

func (r *SurveyRepositoryGorm) Update(survey *entity.Survey) error {
	err := r.gorm.Updates(&survey).Error
	if err != nil {
		return fmt.Errorf("cannot update survey: %w", err)
	}

	return nil
}

func (r *SurveyRepositoryGorm) FilterExisted(ids []string) ([]string, error) {
	var existedIds []string
	err := r.gorm.Raw("SELECT id FROM `surveys` WHERE id in ?", ids).Scan(&existedIds).Error
	if err != nil {
		return nil, fmt.Errorf("cannot query existing surveys: %w", err)
	}

	return existedIds, nil
}

func (r *SurveyRepositoryGorm) AddQuestion(question *entity.Question) error {
	err := r.gorm.Create(&question).Error
	if err != nil {
		return fmt.Errorf("cannot add question: %w", err)
	}
	return nil
}

func (r *SurveyRepositoryGorm) RemoveQuestion(id string) error {
	err := r.gorm.Where("id = ?", id).Delete(&entity.Question{}).Error
	if err != nil {
		return fmt.Errorf("cannot remove question: %w", err)
	}
	return nil
}

func (r *SurveyRepositoryGorm) GetQuestionById(id string) (*entity.Question, error) {
	var question entity.Question
	err := r.gorm.Where("id = ?", id).First(&question).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("cannot get question: %w", err)
	}

	return &question, nil
}

func (r *SurveyRepositoryGorm) GetQuestionsBySurveyId(surveyID string) ([]entity.Question, error) {
	var questions []entity.Question
	err := r.gorm.Where("survey_id = ?", surveyID).Find(&questions).Error
	if err != nil {
		return nil, fmt.Errorf("cannot get questions: %w", err)
	}
	return questions, nil
}

func (r *SurveyRepositoryGorm) UpdateQuestion(question *entity.Question) error {
	err := r.gorm.Updates(&question).Error
	if err != nil {
		return fmt.Errorf("cannot update question: %w", err)
	}
	return nil
}

func (r *SurveyRepositoryGorm) GetPOsByCourseId(courseID string) (map[string][]string, error) {

	return nil, nil
}

func (r *SurveyRepositoryGorm) GetSurveysWithCourseAndOutcomes() ([]entity.SurveyWithCourseAndOutcomes, error) {
	var surveys []entity.SurveyWithCourseAndOutcomes

	rows, err := r.gorm.Raw(`
	SELECT
		s.id AS survey_id,
		s.title AS survey_title,
		s.description,
		s.is_complete,
		c.id AS course_id,
		c.name AS course_name,
		c.code AS course_code,
		c.academic_year,
	    IFNULL(GROUP_CONCAT(DISTINCT q.po_id), '') AS pos,
	    IFNULL(GROUP_CONCAT(DISTINCT q.plo_id), '') AS plos,
	    IFNULL(GROUP_CONCAT(DISTINCT q.so_id), '') AS sos
	FROM
		survey s
	LEFT JOIN course c ON
		s.course_id = c.id
	LEFT JOIN question q ON
		s.id = q.survey_id
	GROUP BY
		s.id,
		c.id;
	`).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var survey entity.SurveyWithCourseAndOutcomes
		var poStr, ploStr, soStr string
		// Scan data from SQL query
		if err := rows.Scan(
			&survey.SurveyId, &survey.SurveyTitle, &survey.Description, &survey.IsComplete,
			&survey.CourseId, &survey.CourseName, &survey.CourseCode, &survey.AcademicYear,
			&poStr, &ploStr, &soStr); err != nil {
			return nil, err
		}

		// Convert comma-separated strings to lists (removing duplicates)
		survey.POs = splitAndFilter(poStr)
		survey.PLOs = splitAndFilter(ploStr)
		survey.SOs = splitAndFilter(soStr)

		surveys = append(surveys, survey)
	}

	return surveys, nil
}

func splitAndFilter(s string) []string {
	if s == "" {
		return []string{}
	}
	items := strings.Split(s, ",")
	unique := make(map[string]bool)
	var result []string

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" && !unique[item] {
			unique[item] = true
			result = append(result, item)
		}
	}
	return result
}
