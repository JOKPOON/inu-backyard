package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
)

func (c assignmentController) GetGroupByCourseId(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	assignmentGroups, err := c.AssignmentUseCase.GetGroupByCourseId(courseId, false)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignmentGroups)
}

func (c assignmentController) CreateGroup(ctx *fiber.Ctx) error {
	var payload entity.CreateAssignmentGroupPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.AssignmentUseCase.CreateGroup(payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c assignmentController) UpdateGroup(ctx *fiber.Ctx) error {
	var payload entity.UpdateAssignmentGroupPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("assignmentGroupId")

	err := c.AssignmentUseCase.UpdateGroup(id, payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c assignmentController) DeleteGroup(ctx *fiber.Ctx) error {
	id := ctx.Params("assignmentGroupId")

	err := c.AssignmentUseCase.DeleteGroup(id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
