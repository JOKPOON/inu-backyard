package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type GradeController struct {
	GradeUseCase entity.GradeUseCase
	Validator    validator.PayloadValidator
}

func NewGradeController(validator validator.PayloadValidator, gradeUseCase entity.GradeUseCase) *GradeController {
	return &GradeController{
		GradeUseCase: gradeUseCase,
		Validator:    validator,
	}
}

func (c GradeController) GetAll(ctx *fiber.Ctx) error {
	studentId := ctx.Query("studentId")
	if studentId != "" {
		grade, err := c.GradeUseCase.GetByStudentId(studentId)
		if err != nil {
			return err
		}
		return response.NewSuccessResponse(ctx, fiber.StatusOK, grade)
	}

	grades, err := c.GradeUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, grades)
}

func (c GradeController) GetById(ctx *fiber.Ctx) error {
	gradeId := ctx.Params("gradeId")

	grade, err := c.GradeUseCase.GetById(gradeId)
	if err != nil {
		return err
	}

	if grade == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, grade)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, grade)
}

func (c GradeController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateGradePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.GradeUseCase.Create(payload.StudentId, payload.Year, payload.Grade)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c GradeController) CreateMany(ctx *fiber.Ctx) error {
	var payload entity.CreateManyGradesPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.GradeUseCase.CreateMany(payload.StudentGrade, payload.Year, payload.SemesterSequence)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c GradeController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateGradePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("gradeId")

	err := c.GradeUseCase.Update(id, &entity.Grade{
		StudentId: payload.StudentId,
		Grade:     payload.Grade,
	})
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c GradeController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("gradeId")

	err := c.GradeUseCase.Delete(id)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
