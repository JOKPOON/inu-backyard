package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type SubStudentOutcomeController struct {
	StudentOutcomeUseCase entity.StudentOutcomeUseCase
	Validator             validator.PayloadValidator
}

func NewSubStudentOutcomeController(validator validator.PayloadValidator, studentOutcomeUseCase entity.StudentOutcomeUseCase) *SubStudentOutcomeController {
	return &SubStudentOutcomeController{
		StudentOutcomeUseCase: studentOutcomeUseCase,
		Validator:             validator,
	}
}

func (c SubStudentOutcomeController) GetAll(ctx *fiber.Ctx) error {
	ssos, err := c.StudentOutcomeUseCase.GetAllSubSO()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, ssos)
}

func (c SubStudentOutcomeController) GetById(ctx *fiber.Ctx) error {
	sploId := ctx.Params("ssoId")

	splo, err := c.StudentOutcomeUseCase.GetSubSOById(sploId)
	if err != nil {
		return err
	}

	if splo == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, splo)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, splo)
}

func (c SubStudentOutcomeController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateSubStudentOutcomePayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.StudentOutcomeUseCase.CreateSubSO(payload.SubStudentOutcomes)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c SubStudentOutcomeController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateSubStudentOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("ssoId")

	err := c.StudentOutcomeUseCase.UpdateSubSO(id, &payload)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c SubStudentOutcomeController) Delete(ctx *fiber.Ctx) error {
	ssoId := ctx.Params("ssoId")

	err := c.StudentOutcomeUseCase.DeleteSubSO(ssoId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
