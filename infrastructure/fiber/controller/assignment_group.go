package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
)

func (c AssignmentController) GetAllGroup(ctx *fiber.Ctx) error {
	assignmentGroups, err := c.AssignmentUseCase.GetAllGroup()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignmentGroups)
}

func (c AssignmentController) GetGroupByCourseId(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")
	groupId := ctx.Query("groupId")
	withAssignment := ctx.Query("withAssignment")

	assignmentGroups, err := c.AssignmentUseCase.GetGroupByCourseId(courseId, groupId, withAssignment == "true")
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignmentGroups)
}

func (c AssignmentController) CreateGroup(ctx *fiber.Ctx) error {
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

func (c AssignmentController) UpdateGroup(ctx *fiber.Ctx) error {
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

func (c AssignmentController) DeleteGroup(ctx *fiber.Ctx) error {
	id := ctx.Params("assignmentGroupId")

	err := c.AssignmentUseCase.DeleteGroup(id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
