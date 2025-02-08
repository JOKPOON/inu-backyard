package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type QuestionController struct {
	SurveyUseCase entity.SurveyUseCase
	Validator     validator.PayloadValidator
}

func NewQuestionController(validator validator.PayloadValidator, surveyUseCase entity.SurveyUseCase) *QuestionController {
	return &QuestionController{
		SurveyUseCase: surveyUseCase,
		Validator:     validator,
	}
}

// CreateQuestion handles adding a new question to a survey
func (c QuestionController) Create(ctx *fiber.Ctx) error {
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
func (c QuestionController) GetBySurveyId(ctx *fiber.Ctx) error {
	surveyId := ctx.Params("surveyId")

	questions, err := c.SurveyUseCase.GetQuestionsBySurveyId(surveyId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, questions)
}

// GetQuestion retrieves a question by ID
func (c QuestionController) GetQuestion(ctx *fiber.Ctx) error {
	id := ctx.Params("questionId")

	question, err := c.SurveyUseCase.GetQuestionById(id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, question)
}

// UpdateQuestion modifies an existing question
func (c QuestionController) Update(ctx *fiber.Ctx) error {
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
func (c QuestionController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("questionId")

	err := c.SurveyUseCase.DeleteQuestion(id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
