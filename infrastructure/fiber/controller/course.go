package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/middleware"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/request"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
	"github.com/team-inu/inu-backyard/usecase"
)

type courseController struct {
	courseUseCase   entity.CourseUseCase
	importerUseCase usecase.ImporterUseCase
	Validator       validator.PayloadValidator
}

func NewCourseController(validator validator.PayloadValidator, courseUseCase entity.CourseUseCase, importerUseCase usecase.ImporterUseCase) *courseController {
	return &courseController{
		courseUseCase:   courseUseCase,
		importerUseCase: importerUseCase,
		Validator:       validator,
	}
}

func (c courseController) GetAll(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)

	var courses []entity.Course
	var err error

	if user.IsRoles([]entity.UserRole{entity.UserRoleHeadOfCurriculum, entity.UserRoleModerator, entity.UserRoleTABEEManager}) {
		courses, err = c.courseUseCase.GetAll()
	} else {
		courses, err = c.courseUseCase.GetByUserId(user.Id)
	}

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, courses)
}

func (c courseController) GetById(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	course, err := c.courseUseCase.GetById(courseId)
	if err != nil {
		return err
	}

	if course == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, course)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, course)
}

func (c courseController) GetByUserId(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	courses, err := c.courseUseCase.GetByUserId(userId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, courses)
}

func (c courseController) Create(ctx *fiber.Ctx) error {
	var payload request.CreateCourseRequestPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	err := c.courseUseCase.Create(
		*user,
		payload.SemesterId,
		payload.UserId,
		payload.Name,
		payload.Code,
		payload.Curriculum,
		payload.Description,
		payload.ExpectedPassingCloPercentage,
		payload.AcademicYear,
		payload.GraduateYear,
		payload.ProgramYear,
		payload.CriteriaGrade,
	)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c courseController) Update(ctx *fiber.Ctx) error {
	var payload request.UpdateCourseRequestPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("courseId")

	user := middleware.GetUserFromCtx(ctx)

	err := c.courseUseCase.Update(
		*user,
		id,
		payload.Name,
		payload.Code,
		payload.Curriculum,
		payload.Description,
		payload.ExpectedPassingCloPercentage,
		payload.AcademicYear,
		payload.GraduateYear,
		payload.ProgramYear,
		payload.CriteriaGrade,
		*payload.IsPortfolioCompleted,
	)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c courseController) Delete(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	user := middleware.GetUserFromCtx(ctx)

	err := c.importerUseCase.UpdateOrCreate(
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
