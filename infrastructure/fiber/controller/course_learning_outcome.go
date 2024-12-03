package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type CourseLearningOutcomeController struct {
	CourseLearningOutcomeUseCase entity.CourseLearningOutcomeUseCase
	Validator                    validator.PayloadValidator
}

func NewCourseLearningOutcomeController(validator validator.PayloadValidator, courseLearningOutcomeUseCase entity.CourseLearningOutcomeUseCase) *CourseLearningOutcomeController {
	return &CourseLearningOutcomeController{
		CourseLearningOutcomeUseCase: courseLearningOutcomeUseCase,
		Validator:                    validator,
	}
}

func (c CourseLearningOutcomeController) GetAll(ctx *fiber.Ctx) error {
	clos, err := c.CourseLearningOutcomeUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, clos)

}

func (c CourseLearningOutcomeController) GetById(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")

	clo, err := c.CourseLearningOutcomeUseCase.GetById(cloId)
	if err != nil {
		return err
	}

	if clo == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, clo)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, clo)
}

func (c CourseLearningOutcomeController) GetByCourseId(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	clos, err := c.CourseLearningOutcomeUseCase.GetByCourseId(courseId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, clos)
}

func (c CourseLearningOutcomeController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateCourseLearningOutcomePayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.CourseLearningOutcomeUseCase.Create(payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c CourseLearningOutcomeController) CreateLinkSubProgramLearningOutcome(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")
	var payload entity.CreateLinkSubProgramLearningOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.CourseLearningOutcomeUseCase.CreateLinkSubProgramLearningOutcome(cloId, payload.SubProgramLearningOutcomeIds)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c CourseLearningOutcomeController) CreateLinkSubStudentOutcome(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")
	var payload entity.CreateLinkSubStudentOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.CourseLearningOutcomeUseCase.CreateLinkSubStudentOutcome(cloId, payload.SubStudentOutcomeIds)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c CourseLearningOutcomeController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateCourseLearningOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("cloId")

	err := c.CourseLearningOutcomeUseCase.Update(id, payload)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c CourseLearningOutcomeController) Delete(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")

	err := c.CourseLearningOutcomeUseCase.Delete(cloId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c CourseLearningOutcomeController) DeleteLinkSubProgramLearningOutcome(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")
	subPloId := ctx.Params("sploId")

	err := c.CourseLearningOutcomeUseCase.DeleteLinkSubProgramLearningOutcome(cloId, subPloId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c CourseLearningOutcomeController) DeleteLinkSubStudentOutcome(ctx *fiber.Ctx) error {
	cloId := ctx.Params("cloId")
	subSoId := ctx.Params("ssoId")

	err := c.CourseLearningOutcomeUseCase.DeleteLinkSubStudentOutcome(cloId, subSoId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
