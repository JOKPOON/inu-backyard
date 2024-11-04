package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type AssignmentController struct {
	AssignmentUseCase entity.AssignmentUseCase
	Validator         validator.PayloadValidator
}

func NewAssignmentController(validator validator.PayloadValidator, assignmentUseCase entity.AssignmentUseCase) *AssignmentController {
	return &AssignmentController{
		AssignmentUseCase: assignmentUseCase,
		Validator:         validator,
	}
}

func (c AssignmentController) GetById(ctx *fiber.Ctx) error {
	assignmentId := ctx.Params("assignmentId")

	assignment, err := c.AssignmentUseCase.GetById(assignmentId)

	if err != nil {
		return err
	}

	if assignment == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, assignment)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignment)
}

func (c AssignmentController) GetByCourseId(ctx *fiber.Ctx) error {
	var payload entity.GetAssignmentsByCourseIdPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	assignments, err := c.AssignmentUseCase.GetByCourseId(payload.CourseId)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignments)
}

func (c AssignmentController) GetAll(ctx *fiber.Ctx) error {
	assignments, err := c.AssignmentUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignments)
}

func (c AssignmentController) GetByGroupId(ctx *fiber.Ctx) error {
	assignmentGroupId := ctx.Params("assignmentGroupId")

	assignmentGroups, err := c.AssignmentUseCase.GetByGroupId(assignmentGroupId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignmentGroups)
}

func (c AssignmentController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateAssignmentPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.AssignmentUseCase.Create(payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c AssignmentController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateAssignmentPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("assignmentId")

	err := c.AssignmentUseCase.Update(id, payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c AssignmentController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("assignmentId")

	err := c.AssignmentUseCase.Delete(id)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c AssignmentController) CreateLinkCourseLearningOutcome(ctx *fiber.Ctx) error {
	assignmentId := ctx.Params("assignmentId")
	var payload entity.CreateLinkCourseLearningOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.AssignmentUseCase.CreateLinkCourseLearningOutcome(assignmentId, payload.CourseLearningOutcomeIds)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c AssignmentController) DeleteLinkCourseLearningOutcome(ctx *fiber.Ctx) error {
	assignmentId := ctx.Params("assignmentId")
	cloId := ctx.Params("cloId")

	err := c.AssignmentUseCase.DeleteLinkCourseLearningOutcome(assignmentId, cloId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
