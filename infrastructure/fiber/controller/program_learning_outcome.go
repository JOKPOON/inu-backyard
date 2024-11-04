package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type ProgramLearningOutcomeController struct {
	ProgramLearningOutcomeUseCase entity.ProgramLearningOutcomeUseCase
	Validator                     validator.PayloadValidator
}

func NewProgramLearningOutcomeController(validator validator.PayloadValidator, programLearningOutcomeUseCase entity.ProgramLearningOutcomeUseCase) *ProgramLearningOutcomeController {
	return &ProgramLearningOutcomeController{
		ProgramLearningOutcomeUseCase: programLearningOutcomeUseCase,
		Validator:                     validator,
	}
}

func (c ProgramLearningOutcomeController) GetAll(ctx *fiber.Ctx) error {
	plos, err := c.ProgramLearningOutcomeUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, plos)
}

func (c ProgramLearningOutcomeController) GetById(ctx *fiber.Ctx) error {
	ploId := ctx.Params("ploId")

	plo, err := c.ProgramLearningOutcomeUseCase.GetById(ploId)
	if err != nil {
		return err
	}

	if plo == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, plo)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, plo)
}

func (c ProgramLearningOutcomeController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateProgramLearningOutcomePayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.ProgramLearningOutcomeUseCase.Create(payload.ProgramLearningOutcomes)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c ProgramLearningOutcomeController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateProgramLearningOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("ploId")

	err := c.ProgramLearningOutcomeUseCase.Update(id, &entity.ProgramLearningOutcome{
		Code:            payload.Code,
		DescriptionThai: payload.DescriptionThai,
		DescriptionEng:  *payload.DescriptionEng, // because description eng can be empty string
		ProgramYear:     payload.ProgramYear,
		ProgrammeName:   payload.ProgrammeName,
	})

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c ProgramLearningOutcomeController) Delete(ctx *fiber.Ctx) error {
	ploId := ctx.Params("ploId")

	err := c.ProgramLearningOutcomeUseCase.Delete(ploId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
