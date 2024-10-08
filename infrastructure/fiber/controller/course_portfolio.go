package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/request"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type coursePortfolioController struct {
	coursePortfolioUseCase entity.CoursePortfolioUseCase
	Validator              validator.PayloadValidator
}

func NewCoursePortfolioController(validator validator.PayloadValidator, coursePortfolioUseCase entity.CoursePortfolioUseCase) *coursePortfolioController {
	return &coursePortfolioController{
		coursePortfolioUseCase: coursePortfolioUseCase,
		Validator:              validator,
	}
}

func (c coursePortfolioController) Generate(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	coursePortfolio, err := c.coursePortfolioUseCase.Generate(courseId)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, coursePortfolio)
}

func (c coursePortfolioController) GetCloPassingStudentsByCourseId(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	records, err := c.coursePortfolioUseCase.GetCloPassingStudentsByCourseId(courseId)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}

func (c coursePortfolioController) GetStudentOutcomeStatusByCourseId(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	records, err := c.coursePortfolioUseCase.GetStudentOutcomesStatusByCourseId(courseId)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}

func (c coursePortfolioController) GetAllProgramLearningOutcomeCourses(ctx *fiber.Ctx) error {
	records, err := c.coursePortfolioUseCase.GetAllProgramLearningOutcomeCourses()
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}

func (c coursePortfolioController) GetAllProgramOutcomeCourses(ctx *fiber.Ctx) error {
	records, err := c.coursePortfolioUseCase.GetAllProgramOutcomeCourses()
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}

func (c coursePortfolioController) Update(ctx *fiber.Ctx) error {
	var payload request.SaveCoursePortfolioPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	courseId := ctx.Params("courseId")

	fmt.Println(payload)
	err := c.coursePortfolioUseCase.UpdateCoursePortfolio(courseId, payload.CourseSummary, payload.CourseDevelopment)
	if err != nil {
		return err
	}

	return nil
}

func (c coursePortfolioController) GetOutcomesByStudentId(ctx *fiber.Ctx) error {
	studentId := ctx.Params("studentId")

	records, err := c.coursePortfolioUseCase.GetOutcomesByStudentId(studentId)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}
