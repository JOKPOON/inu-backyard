package controller

import (
	"fmt"
	"strconv"

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

	err := c.CoursePortfolioUseCase.UpdateCoursePortfolio(courseId, payload.Implementation, payload.EducationOutcomes, payload.ContinuousDevelopment)
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

func (c CoursePortfolioController) GetCourseCloAssessment(ctx *fiber.Ctx) error {
	programmeId := ctx.Params("programmeId")
	toSerm, err := strconv.Atoi(ctx.Query("to"))
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}
	if toSerm == 0 {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}
	fromSerm, err := strconv.Atoi(ctx.Query("from"))
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}
	if fromSerm == 0 {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}

	if toSerm < fromSerm {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}

	file, err := c.CoursePortfolioUseCase.GetCourseCloAssessment(programmeId, fromSerm, toSerm)
	if err != nil {
		return err
	}

	ctx.Set("Content-Type", file.FileType)
	ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.FileName))

	// Send file from disk
	return ctx.SendFile(file.FilePath)
}

func (c CoursePortfolioController) GetCourseLinkedOutcomes(ctx *fiber.Ctx) error {
	programmeId := ctx.Params("programmeId")
	toSerm, err := strconv.Atoi(ctx.Query("to"))
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}
	if toSerm == 0 {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}
	fromSerm, err := strconv.Atoi(ctx.Query("from"))
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}
	if fromSerm == 0 {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}

	if toSerm < fromSerm {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}

	file, err := c.CoursePortfolioUseCase.GetCourseLinkedOutcomes(programmeId, fromSerm, toSerm)
	if err != nil {
		return err
	}

	ctx.Set("Content-Type", file.FileType)
	ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.FileName))

	// Send file from disk
	return ctx.SendFile(file.FilePath)
}

func (c CoursePortfolioController) GetCourseOutcomesSuccessRate(ctx *fiber.Ctx) error {
	programmeId := ctx.Params("programmeId")
	toSerm, err := strconv.Atoi(ctx.Query("to"))
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}
	if toSerm == 0 {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}
	fromSerm, err := strconv.Atoi(ctx.Query("from"))
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}
	if fromSerm == 0 {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}

	if toSerm < fromSerm {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, nil)
	}

	file, err := c.CoursePortfolioUseCase.GetCourseOutcomesSuccessRate(programmeId, fromSerm, toSerm)
	if err != nil {
		return err
	}

	// Set headers to indicate download
	ctx.Set("Content-Type", file.FileType)
	ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.FileName))

	// Send file from disk
	return ctx.SendFile(file.FilePath)
}
