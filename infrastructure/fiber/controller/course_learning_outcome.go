package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type courseLearningOutcomeController struct {
	courseLearningOutcomeUseCase entity.CourseLearningOutcomeUseCase
	Validator                    validator.PayloadValidator
}

func NewCourseLearningOutcomeController(validator validator.PayloadValidator, courseLearningOutcomeUseCase entity.CourseLearningOutcomeUseCase) *courseLearningOutcomeController {
	return &courseLearningOutcomeController{
		courseLearningOutcomeUseCase: courseLearningOutcomeUseCase,
		Validator:                    validator,
	}
}

func (c courseLearningOutcomeController) GetAll(ctx *fiber.Ctx) error {
	clos, err := c.courseLearningOutcomeUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, clos)

}

func (c courseLearningOutcomeController) GetById(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")

	clo, err := c.courseLearningOutcomeUseCase.GetById(cloId)
	if err != nil {
		return err
	}

	if clo == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, clo)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, clo)
}

func (c courseLearningOutcomeController) GetByCourseId(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	clos, err := c.courseLearningOutcomeUseCase.GetByCourseId(courseId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, clos)
}

func (c courseLearningOutcomeController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateCourseLearningOutcomePayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.courseLearningOutcomeUseCase.Create(payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c courseLearningOutcomeController) CreateLinkSubProgramLearningOutcome(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")
	var payload entity.CreateLinkSubProgramLearningOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.courseLearningOutcomeUseCase.CreateLinkSubProgramLearningOutcome(cloId, payload.SubProgramLearningOutcomeId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c courseLearningOutcomeController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateCourseLearningOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("cloId")

	err := c.courseLearningOutcomeUseCase.Update(id, payload)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c courseLearningOutcomeController) Delete(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")

	err := c.courseLearningOutcomeUseCase.Delete(cloId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c courseLearningOutcomeController) DeleteLinkSubProgramLearningOutcome(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")
	subPloId := ctx.Params("sploId")

	err := c.courseLearningOutcomeUseCase.DeleteLinkSubProgramLearningOutcome(cloId, subPloId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
