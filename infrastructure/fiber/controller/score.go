package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/middleware"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type ScoreController struct {
	ScoreUseCase entity.ScoreUseCase
	Validator    validator.PayloadValidator
}

func NewScoreController(validator validator.PayloadValidator, scoreUseCase entity.ScoreUseCase) *ScoreController {
	return &ScoreController{
		ScoreUseCase: scoreUseCase,
		Validator:    validator,
	}
}

func (c ScoreController) GetAll(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)

	var scores []entity.Score
	var err error

	if user.IsRoles([]entity.UserRole{entity.UserRoleHeadOfCurriculum, entity.UserRoleModerator, entity.UserRoleTABEEManager}) {
		scores, err = c.ScoreUseCase.GetAll()
	} else {
		scores, err = c.ScoreUseCase.GetByUserId(user.Id)
	}

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, scores)
}

func (c ScoreController) GetById(ctx *fiber.Ctx) error {
	scoreId := ctx.Params("scoreId")

	score, err := c.ScoreUseCase.GetById(scoreId)
	if err != nil {
		return err
	}

	if score == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, score)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, score)
}

func (c ScoreController) GetByAssignmentId(ctx *fiber.Ctx) error {
	assignmentId := ctx.Params("assignmentId")
	courseId := ctx.Query("courseId")

	assignmentScore, err := c.ScoreUseCase.GetByAssignmentId(assignmentId, courseId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignmentScore)
}

func (c ScoreController) CreateMany(ctx *fiber.Ctx) error {
	var payload entity.BulkCreateScoreRequestPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	user := ctx.Locals("user").(*entity.User)

	err := c.ScoreUseCase.CreateMany(
		user.Id,
		payload.AssignmentId,
		payload.StudentScores,
	)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c ScoreController) Delete(ctx *fiber.Ctx) error {
	scoreId := ctx.Params("scoreId")

	_, err := c.ScoreUseCase.GetById(scoreId)
	if err != nil {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	err = c.ScoreUseCase.Delete(*user, scoreId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c ScoreController) Update(ctx *fiber.Ctx) error {
	scoreId := ctx.Params("scoreId")

	_, err := c.ScoreUseCase.GetById(scoreId)
	if err != nil {
		return err
	}
	var payload entity.UpdateScoreRequestPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	err = c.ScoreUseCase.Update(*user, scoreId, payload.Score)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
