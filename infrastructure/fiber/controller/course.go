package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/middleware"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
	"github.com/team-inu/inu-backyard/usecase"
)

type CourseController struct {
	CourseUseCase   entity.CourseUseCase
	ImporterUseCase usecase.ImporterUseCase
	Validator       validator.PayloadValidator
}

func NewCourseController(validator validator.PayloadValidator, courseUseCase entity.CourseUseCase, importerUseCase usecase.ImporterUseCase) *CourseController {
	return &CourseController{
		CourseUseCase:   courseUseCase,
		ImporterUseCase: importerUseCase,
		Validator:       validator,
	}
}

func (c CourseController) GetAll(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)

	var courses *entity.GetAllCourseResponse
	var err error

	query := ctx.Query("query")
	year := ctx.Query("year")
	program := ctx.Query("program")

	println(query, year, program)

	if user.IsRoles([]entity.UserRole{entity.UserRoleHeadOfCurriculum, entity.UserRoleModerator, entity.UserRoleTABEEManager}) {
		courses, err = c.CourseUseCase.GetAll(query, year, program)
	} else {
		courses, err = c.CourseUseCase.GetByUserId(user.Id, query, year, program)
	}

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, courses)
}

func (c CourseController) GetById(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	course, err := c.CourseUseCase.GetById(courseId)
	if err != nil {
		return err
	}

	if course == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, course)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, course)
}

func (c CourseController) GetByUserId(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")
	query := ctx.Query("query")
	year := ctx.Query("year")
	program := ctx.Query("program")

	courses, err := c.CourseUseCase.GetByUserId(userId, query, year, program)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, courses)
}

func (c CourseController) GetStudentsPassingCLOs(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	resp, err := c.CourseUseCase.GetStudentsPassingCLOs(courseId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, resp)
}

func (c CourseController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateCoursePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	err := c.CourseUseCase.Create(
		*user,
		payload,
	)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c CourseController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateCoursePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("courseId")

	user := middleware.GetUserFromCtx(ctx)

	err := c.CourseUseCase.Update(
		*user,
		id,
		payload,
	)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c CourseController) Delete(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	user := middleware.GetUserFromCtx(ctx)

	err := c.ImporterUseCase.UpdateOrCreate(
		courseId,
		user.Id,
		make([]string, 0),
		make([]usecase.ImportCourseLearningOutcome, 0),
		make([]usecase.ImportAssignmentGroup, 0),
		true,
	)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
