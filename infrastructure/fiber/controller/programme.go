package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type ProgrammeController struct {
	ProgrammeUseCase entity.ProgrammeUseCase
	Validator        validator.PayloadValidator
}

func NewProgrammeController(validator validator.PayloadValidator, programmeUseCase entity.ProgrammeUseCase) *ProgrammeController {
	return &ProgrammeController{
		ProgrammeUseCase: programmeUseCase,
		Validator:        validator,
	}
}

func (c ProgrammeController) GetAll(ctx *fiber.Ctx) error {
	programmes, err := c.ProgrammeUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, programmes)
}

func (c ProgrammeController) GetByName(ctx *fiber.Ctx) error {
	name := ctx.Params("programmeName")

	programme, err := c.ProgrammeUseCase.GetByName(name)
	if err != nil {
		return err
	}

	if programme == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, programme)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, programme)
}

func (c ProgrammeController) GetByNameAndYear(ctx *fiber.Ctx) error {
	name := ctx.Params("programmeName")
	year := ctx.Params("year")

	programme, err := c.ProgrammeUseCase.GetByNameAndYear(name, year)
	if err != nil {
		return err
	}

	if programme == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, programme)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, programme)
}

func (c ProgrammeController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateProgrammePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.ProgrammeUseCase.Create(payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c ProgrammeController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateProgrammePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	name := ctx.Params("programmeName")

	err := c.ProgrammeUseCase.Update(name, &payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c ProgrammeController) Delete(ctx *fiber.Ctx) error {
	name := ctx.Params("programmeName")

	err := c.ProgrammeUseCase.Delete(name)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
