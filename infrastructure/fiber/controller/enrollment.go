package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type EnrollmentController struct {
	EnrollmentUseCase entity.EnrollmentUseCase
	Validator         validator.PayloadValidator
}

func NewEnrollmentController(validator validator.PayloadValidator, enrollmentUseCase entity.EnrollmentUseCase) *EnrollmentController {
	return &EnrollmentController{
		EnrollmentUseCase: enrollmentUseCase,
		Validator:         validator,
	}
}

func (c EnrollmentController) GetAll(ctx *fiber.Ctx) error {
	enrollments, err := c.EnrollmentUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, enrollments)
}

func (c EnrollmentController) GetById(ctx *fiber.Ctx) error {
	enrollmentId := ctx.Params("enrollmentId")

	enrollment, err := c.EnrollmentUseCase.GetById(enrollmentId)
	if err != nil {
		return err
	}

	if enrollment == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, enrollment)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, enrollment)
}

func (c EnrollmentController) GetByCourseId(ctx *fiber.Ctx) error {
	enrollmentId := ctx.Params("courseId")
	query := ctx.Query("query")

	enrollments, err := c.EnrollmentUseCase.GetByCourseId(enrollmentId, query)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, enrollments)
}

func (c EnrollmentController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateEnrollmentsPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.EnrollmentUseCase.CreateMany(payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c EnrollmentController) Update(ctx *fiber.Ctx) error {
	enrollmentId := ctx.Params("enrollmentId")

	var payload entity.UpdateEnrollmentPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.EnrollmentUseCase.Update(enrollmentId, payload.Status)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c EnrollmentController) Delete(ctx *fiber.Ctx) error {
	enrollmentId := ctx.Params("enrollmentId")

	err := c.EnrollmentUseCase.Delete(enrollmentId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
