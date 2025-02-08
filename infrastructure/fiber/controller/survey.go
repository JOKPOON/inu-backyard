package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type SurveyController struct {
	SurveyUseCase entity.SurveyUseCase
	Validator     validator.PayloadValidator
}

func NewSurveyController(validator validator.PayloadValidator, surveyUseCase entity.SurveyUseCase) *SurveyController {
	return &SurveyController{
		SurveyUseCase: surveyUseCase,
		Validator:     validator,
	}
}

// GetAll retrieves all surveys
func (c SurveyController) GetAll(ctx *fiber.Ctx) error {
	surveys, err := c.SurveyUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, surveys)
}

// GetById retrieves a survey by ID
func (c SurveyController) GetById(ctx *fiber.Ctx) error {
	surveyId := ctx.Params("surveyId")

	survey, err := c.SurveyUseCase.GetById(surveyId)
	if err != nil {
		return err
	}

	if survey == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, survey)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, survey)
}

// Create handles creating a new survey
func (c SurveyController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateSurveyRequest

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.SurveyUseCase.Create(&payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

// Update modifies an existing survey
func (c SurveyController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateSurveyRequest

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("surveyId")
	err := c.SurveyUseCase.Update(id, &payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

// Delete removes a survey by ID
func (c SurveyController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("surveyId")

	err := c.SurveyUseCase.Delete(id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

// CreateQuestion handles adding a new question to a survey
func (c SurveyController) CreateQuestion(ctx *fiber.Ctx) error {
	surveyId := ctx.Params("surveyId")
	var payload entity.CreateQuestionRequest

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.SurveyUseCase.CreateQuestion(surveyId, &payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

// GetBySurveyId retrieves all questions from a survey
func (c SurveyController) GetQuestionBySurveyId(ctx *fiber.Ctx) error {
	surveyId := ctx.Params("surveyId")

	questions, err := c.SurveyUseCase.GetQuestionsBySurveyId(surveyId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, questions)
}

// GetQuestionById retrieves a question by ID
func (c SurveyController) GetQuestionById(ctx *fiber.Ctx) error {
	id := ctx.Params("questionId")

	question, err := c.SurveyUseCase.GetQuestionById(id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, question)
}

// UpdateQuestion modifies an existing question
func (c SurveyController) UpdateQuestion(ctx *fiber.Ctx) error {
	var payload entity.UpdateQuestionRequest

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("questionId")
	err := c.SurveyUseCase.UpdateQuestion(id, &payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

// DeleteQuestion removes a question from a survey
func (c SurveyController) DeleteQuestion(ctx *fiber.Ctx) error {
	id := ctx.Params("questionId")

	err := c.SurveyUseCase.DeleteQuestion(id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c SurveyController) GetSurveysWithCourseAndOutcomes(ctx *fiber.Ctx) error {
	surveys, err := c.SurveyUseCase.GetSurveysWithCourseAndOutcomes()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, surveys)
}
