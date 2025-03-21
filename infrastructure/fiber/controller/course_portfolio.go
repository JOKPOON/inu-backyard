package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type CoursePortfolioController struct {
	CoursePortfolioUseCase entity.CoursePortfolioUseCase
	Validator              validator.PayloadValidator
}

func NewCoursePortfolioController(validator validator.PayloadValidator, coursePortfolioUseCase entity.CoursePortfolioUseCase) *CoursePortfolioController {
	return &CoursePortfolioController{
		CoursePortfolioUseCase: coursePortfolioUseCase,
		Validator:              validator,
	}
}

func (c CoursePortfolioController) Generate(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	coursePortfolio, err := c.CoursePortfolioUseCase.Generate(courseId)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, coursePortfolio)
}

func (c CoursePortfolioController) GetCloPassingStudentsByCourseId(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	records, err := c.CoursePortfolioUseCase.GetCloPassingStudentsByCourseId(courseId)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}

func (c CoursePortfolioController) GetStudentOutcomeStatusByCourseId(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	records, err := c.CoursePortfolioUseCase.GetStudentOutcomesStatusByCourseId(courseId)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}

func (c CoursePortfolioController) GetAllProgramLearningOutcomeCourses(ctx *fiber.Ctx) error {
	records, err := c.CoursePortfolioUseCase.GetAllProgramLearningOutcomeCourses()
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}

func (c CoursePortfolioController) GetAllProgramOutcomeCourses(ctx *fiber.Ctx) error {
	records, err := c.CoursePortfolioUseCase.GetAllProgramOutcomeCourses()
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}

func (c CoursePortfolioController) Update(ctx *fiber.Ctx) error {
	var payload entity.SaveCoursePortfolioPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	courseId := ctx.Params("courseId")

	fmt.Println(payload)
	err := c.CoursePortfolioUseCase.UpdateCoursePortfolio(courseId, payload.CourseSummary, payload.CourseDevelopment)
	if err != nil {
		return err
	}

	return nil
}

func (c CoursePortfolioController) GetOutcomesByStudentId(ctx *fiber.Ctx) error {
	studentId := ctx.Params("studentId")

	records, err := c.CoursePortfolioUseCase.GetOutcomesByStudentId(studentId)
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, records)
}

func (c CoursePortfolioController) GetCourseResult(ctx *fiber.Ctx) error {
	courseId := ctx.Params("courseId")

	result, err := c.CoursePortfolioUseCase.CalculateGradeDistribution(courseId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, result)
}
