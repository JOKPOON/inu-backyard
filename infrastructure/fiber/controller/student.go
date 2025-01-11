package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type StudentController struct {
	StudentUseCase entity.StudentUseCase
	Validator      validator.PayloadValidator
}

func NewStudentController(validator validator.PayloadValidator, studentUseCase entity.StudentUseCase) *StudentController {
	return &StudentController{
		StudentUseCase: studentUseCase,
		Validator:      validator,
	}
}

func (c StudentController) GetAll(ctx *fiber.Ctx) error {
	students, err := c.StudentUseCase.GetAll()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, students)
}

func (c StudentController) GetById(ctx *fiber.Ctx) error {
	studentId := ctx.Params("studentId")

	student, err := c.StudentUseCase.GetById(studentId)

	if err != nil {
		return err
	}

	if student == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, student)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, student)
}

func (c StudentController) GetStudents(ctx *fiber.Ctx) error {
	var payload entity.GetStudentsByParamsPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	students, err := c.StudentUseCase.GetByParams(&entity.Student{
		ProgrammeId:    payload.ProgrammeId,
		DepartmentName: payload.DepartmentName,
		Year:           payload.Year,
	}, -1, -1)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, students)
}

func (c StudentController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateStudentPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.StudentUseCase.CreateMany([]entity.CreateStudentPayload{payload})
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c StudentController) CreateMany(ctx *fiber.Ctx) error {
	var payload entity.CreateBulkStudentsPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.StudentUseCase.CreateMany(payload.Students)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c StudentController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateStudentPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("studentId")
	err := c.StudentUseCase.Update(id, &payload)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c StudentController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("studentId")

	err := c.StudentUseCase.Delete(id)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c StudentController) GetAllSchools(ctx *fiber.Ctx) error {
	schools, err := c.StudentUseCase.GetAllSchools()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, map[string]interface{}{
		"schools": schools,
	})
}
func (c StudentController) GetAllAdmissions(ctx *fiber.Ctx) error {
	admissions, err := c.StudentUseCase.GetAllAdmissions()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, map[string]interface{}{
		"admissions": admissions,
	})
}
