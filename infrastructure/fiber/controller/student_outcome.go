package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type StudentOutcomeController struct {
	StudentOutcomeUsecase entity.StudentOutcomeUseCase
	Validator             validator.PayloadValidator
}

func NewStudentOutcomeController(validator validator.PayloadValidator, studentOutcomeUsecase entity.StudentOutcomeUseCase) *StudentOutcomeController {
	return &StudentOutcomeController{
		StudentOutcomeUsecase: studentOutcomeUsecase,
		Validator:             validator,
	}
}

func (c StudentOutcomeController) GetAll(ctx *fiber.Ctx) error {
	programId := ctx.Query("program_id")
	plos, err := c.StudentOutcomeUsecase.GetAll(programId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, plos)
}

func (c StudentOutcomeController) GetById(ctx *fiber.Ctx) error {
	soId := ctx.Params("soId")

	so, err := c.StudentOutcomeUsecase.GetById(soId)
	if err != nil {
		return err
	}

	if so == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, so)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, so)
}

func (c StudentOutcomeController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateStudentOutcomePayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.StudentOutcomeUsecase.Create(payload.StudentOutcomes)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c StudentOutcomeController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateStudentOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("soId")

	err := c.StudentOutcomeUsecase.Update(id, &payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c StudentOutcomeController) Delete(ctx *fiber.Ctx) error {
	soId := ctx.Params("soId")

	err := c.StudentOutcomeUsecase.Delete(soId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
