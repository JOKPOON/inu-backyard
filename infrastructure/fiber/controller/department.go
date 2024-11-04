package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	request "github.com/team-inu/inu-backyard/infrastructure/fiber/request"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type DepartmentController struct {
	DepartmentUseCase entity.DepartmentUseCase
	Validator         validator.PayloadValidator
}

func NewDepartmentController(validator validator.PayloadValidator, departmentUseCase entity.DepartmentUseCase) *DepartmentController {
	return &DepartmentController{
		DepartmentUseCase: departmentUseCase,
		Validator:         validator,
	}
}

func (c DepartmentController) Create(ctx *fiber.Ctx) error {
	var payload request.CreateDepartmentRequestPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.DepartmentUseCase.Create(&entity.Department{
		Name:        payload.Name,
		FacultyName: payload.FacultyName,
	})

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c DepartmentController) Delete(ctx *fiber.Ctx) error {
	departmentName := ctx.Params("departmentName")

	_, err := c.DepartmentUseCase.GetByName(departmentName)

	if err != nil {
		return err
	}

	err = c.DepartmentUseCase.Delete(departmentName)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c DepartmentController) GetAll(ctx *fiber.Ctx) error {
	departments, err := c.DepartmentUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, departments)
}

func (c DepartmentController) GetByName(ctx *fiber.Ctx) error {
	departmentName := ctx.Params("departmentName")

	department, err := c.DepartmentUseCase.GetByName(departmentName)

	if err != nil {
		return err
	}

	if department == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, department)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, department)
}

func (c DepartmentController) Update(ctx *fiber.Ctx) error {
	departmentName := ctx.Params("departmentName")

	_, err := c.DepartmentUseCase.GetByName(departmentName)

	if err != nil {
		return err
	}

	var payload request.UpdateDepartmentRequestPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err = c.DepartmentUseCase.Update(&entity.Department{
		Name:        departmentName,
		FacultyName: payload.FacultyName,
	}, payload.NewName)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
