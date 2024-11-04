package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	request "github.com/team-inu/inu-backyard/infrastructure/fiber/request"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type ProgramOutcomeController struct {
	ProgramOutcomeUseCase entity.ProgramOutcomeUseCase
	Validator             validator.PayloadValidator
}

func NewProgramOutcomeController(validator validator.PayloadValidator, programOutcomeUseCase entity.ProgramOutcomeUseCase) *ProgramOutcomeController {
	return &ProgramOutcomeController{
		ProgramOutcomeUseCase: programOutcomeUseCase,
		Validator:             validator,
	}
}

func (c ProgramOutcomeController) GetAll(ctx *fiber.Ctx) error {
	pos, err := c.ProgramOutcomeUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, pos)
}

func (c ProgramOutcomeController) GetById(ctx *fiber.Ctx) error {
	poId := ctx.Params("poId")

	po, err := c.ProgramOutcomeUseCase.GetById(poId)
	if err != nil {
		return err
	}

	if po == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, po)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, po)
}

func (c ProgramOutcomeController) Create(ctx *fiber.Ctx) error {
	var payload request.CreateProgramOutcomePayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	pos := make([]entity.ProgramOutcome, 0, len(payload.ProgramOutcomes))
	for _, po := range payload.ProgramOutcomes {

		pos = append(pos, entity.ProgramOutcome{
			Code:        po.Code,
			Name:        po.Name,
			Description: po.Description,
		})
	}

	err := c.ProgramOutcomeUseCase.Create(pos)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c ProgramOutcomeController) Update(ctx *fiber.Ctx) error {
	var payload request.UpdateProgramOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("poId")

	err := c.ProgramOutcomeUseCase.Update(id, &entity.ProgramOutcome{
		Code:        payload.Code,
		Name:        payload.Name,
		Description: payload.Description,
	})

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c ProgramOutcomeController) Delete(ctx *fiber.Ctx) error {
	poId := ctx.Params("poId")

	_, err := c.ProgramOutcomeUseCase.GetById(poId)
	if err != nil {
		return err
	}

	err = c.ProgramOutcomeUseCase.Delete(poId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
