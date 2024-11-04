package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/middleware"
	request "github.com/team-inu/inu-backyard/infrastructure/fiber/request"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
	"github.com/team-inu/inu-backyard/usecase"
)

type ImporterController struct {
	ImporterUseCase usecase.ImporterUseCase
	Validator       validator.PayloadValidator
}

func NewImporterController(validator validator.PayloadValidator, importerUseCase usecase.ImporterUseCase) ImporterController {
	return ImporterController{
		ImporterUseCase: importerUseCase,
		Validator:       validator,
	}
}

func (c ImporterController) Import(ctx *fiber.Ctx) error {
	var payload request.ImportCoursePayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	err := c.ImporterUseCase.UpdateOrCreate(
		payload.CourseId,
		user.Id,
		payload.StudentIds,
		payload.CourseLearningOutcomes,
		payload.AssignmentGroups,
		false,
	)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
