package usecase

import (
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
)

type surveyUseCase struct {
	surveyRepo entity.SurveyRepository
}

func NewSurveyUseCase(surveyRepo entity.SurveyRepository) entity.SurveyUseCase {
	return &surveyUseCase{
		surveyRepo: surveyRepo,
	}
}

func (u surveyUseCase) GetById(id string) (*entity.Survey, error) {
	survey, err := u.surveyRepo.GetById(id)
	if err != nil {
		return nil, errs.New(errs.ErrQuerySurvey, "cannot get survey by id %s", id, err)
	}
	return survey, nil
}

func (u surveyUseCase) GetAll() ([]entity.Survey, error) {
	surveys, err := u.surveyRepo.GetAll()
	if err != nil {
		return nil, errs.New(errs.ErrQuerySurvey, "cannot get all surveys", err)
	}
	return surveys, nil
}

func (u surveyUseCase) Create(request *entity.CreateSurveyRequest) error {
	for idx := range request.Questions {
		request.Questions[idx].Id = ulid.Make().String()
	}

	survey := &entity.Survey{
		Id:          ulid.Make().String(),
		Title:       request.Title,
		Description: request.Description,
		CourseId:    request.CourseId,
		IsComplete:  request.IsComplete,
		Questions:   request.Questions,
		CreateAt:    time.Now(),
	}

	err := u.surveyRepo.Create(survey)
	if err != nil {
		return errs.New(errs.ErrCreateSurvey, "cannot create survey", err)
	}

	return nil
}

func (u surveyUseCase) Update(id string, request *entity.UpdateSurveyRequest) error {
	existSurvey, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get survey id %s to update", id, err)
	} else if existSurvey == nil {
		return errs.New(errs.ErrSurveyNotFound, "survey id %s not found to update", id)
	}

	survey := &entity.Survey{
		Id:          id,
		Title:       request.Title,
		Description: request.Description,
		IsComplete:  request.IsComplete,
		Questions:   request.Questions,
		CourseId:    request.CourseId,
	}

	err = u.surveyRepo.Update(survey)
	if err != nil {
		return errs.New(errs.ErrUpdateSurvey, "cannot update survey by id %s", id, err)
	}

	return nil
}

func (u surveyUseCase) Delete(id string) error {
	existSurvey, err := u.GetById(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get survey id %s to delete", id, err)
	} else if existSurvey == nil {
		return errs.New(errs.ErrSurveyNotFound, "survey id %s not found to delete", id)
	}

	err = u.surveyRepo.Delete(id)
	if err != nil {
		return errs.New(errs.ErrDeleteSurvey, "cannot delete survey", err)
	}

	return nil
}

func (u surveyUseCase) GetByCourseId(courseId string) (*entity.Survey, error) {
	surveys, err := u.surveyRepo.GetByCourseId(courseId)
	if err != nil {
		return nil, errs.New(errs.ErrQuerySurvey, "cannot get survey by course id %s", courseId, err)
	}
	return surveys, nil
}

func (u surveyUseCase) GetQuestionsBySurveyId(surveyId string) ([]entity.Question, error) {
	questions, err := u.surveyRepo.GetQuestionsBySurveyId(surveyId)
	if err != nil {
		return nil, errs.New(errs.ErrQuerySurvey, "cannot get questions by survey id %s", surveyId, err)
	}
	return questions, nil
}

func (u surveyUseCase) GetQuestionById(questionId string) (*entity.Question, error) {
	question, err := u.surveyRepo.GetQuestionById(questionId)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get question by id %s", questionId, err)
	}

	return question, nil
}

func (u surveyUseCase) CreateQuestion(surveyId string, question *entity.CreateQuestionRequest) error {
	existSurvey, err := u.GetById(surveyId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get survey id %s to create question", surveyId, err)
	} else if existSurvey == nil {
		return errs.New(errs.ErrSurveyNotFound, "survey id %s not found to create question", surveyId)
	}

	questionToCreate := &entity.Question{
		Id:       ulid.Make().String(),
		Question: question.Question,
		POId:     question.POId,
		PLOId:    question.PLOId,
		SOId:     question.SOId,
		SurveyId: surveyId,
	}

	err = u.surveyRepo.AddQuestion(questionToCreate)
	if err != nil {
		return errs.New(errs.SameCode, "cannot create question", err)
	}

	return nil
}

func (u surveyUseCase) UpdateQuestion(questionId string, question *entity.UpdateQuestionRequest) error {
	existQuestion, err := u.surveyRepo.GetQuestionById(questionId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get question id %s to update", questionId, err)
	} else if existQuestion == nil {
		return errs.New(errs.SameCode, "question id %s not found to update", questionId)
	}

	questionToUpdate := &entity.Question{
		Id:       questionId,
		Question: question.Question,
		POId:     question.POId,
		PLOId:    question.PLOId,
		SOId:     question.SOId,
	}

	err = u.surveyRepo.UpdateQuestion(questionToUpdate)
	if err != nil {
		return errs.New(errs.SameCode, "cannot update question by id %s", questionId, err)
	}

	return nil
}

func (u surveyUseCase) DeleteQuestion(questionId string) error {
	existQuestion, err := u.surveyRepo.GetQuestionById(questionId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get question id %s to delete", questionId, err)

	} else if existQuestion == nil {
		return errs.New(errs.SameCode, "question id %s not found to delete", questionId)
	}

	err = u.surveyRepo.RemoveQuestion(questionId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot delete question by id %s", questionId, err)
	}

	return nil
}

func (u surveyUseCase) GetSurveysWithCourseAndOutcomes() ([]entity.SurveyWithCourseAndOutcomes, error) {
	surveys, err := u.surveyRepo.GetSurveysWithCourseAndOutcomes()
	if err != nil {
		return nil, errs.New(errs.ErrQuerySurvey, "cannot get surveys with course and outcomes", err)
	}
	return surveys, nil
}
