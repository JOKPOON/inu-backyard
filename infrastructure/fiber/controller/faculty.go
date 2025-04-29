package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type FacultyController struct {
	FacultyUseCase entity.FacultyUseCase
	Validator      validator.PayloadValidator
}

func NewFacultyController(validator validator.PayloadValidator, facultyUseCase entity.FacultyUseCase) *FacultyController {
	return &FacultyController{
		FacultyUseCase: facultyUseCase,
		Validator:      validator,
	}
}

func (c FacultyController) GetAll(ctx *fiber.Ctx) error {
	faculties, err := c.FacultyUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, faculties)
}

func (c FacultyController) GetById(ctx *fiber.Ctx) error {
	facultyId := ctx.Params("facultyId")

	faculty, err := c.FacultyUseCase.GetById(facultyId)
	if err != nil {
		return err
	}

	if faculty == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, faculty)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, faculty)
}

func (c FacultyController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateFacultyRequestPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.FacultyUseCase.Create(&entity.Faculty{
		Id:     ulid.Make().String(),
		NameTH: payload.NameTH,
		NameEN: payload.NameEN,
	})
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c FacultyController) Update(ctx *fiber.Ctx) error {
	facultyId := ctx.Params("facultyId")

	_, err := c.FacultyUseCase.GetById(facultyId)
	if err != nil {
		return err
	}

	var payload entity.UpdateFacultyRequestPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err = c.FacultyUseCase.Update(&entity.Faculty{
		Id:     facultyId,
		NameTH: payload.NameTH,
		NameEN: payload.NameEN,
	})
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c FacultyController) Delete(ctx *fiber.Ctx) error {
	facultyId := ctx.Params("facultyId")

	_, err := c.FacultyUseCase.GetById(facultyId)
	if err != nil {
		return err
	}

	err = c.FacultyUseCase.Delete(facultyId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
