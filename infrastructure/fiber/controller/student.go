package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	request "github.com/team-inu/inu-backyard/infrastructure/fiber/request"
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
	var payload request.GetStudentsByParamsPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	students, err := c.StudentUseCase.GetByParams(&entity.Student{
		ProgrammeName:  payload.ProgrammeName,
		DepartmentName: payload.DepartmentName,
		Year:           payload.Year,
	}, -1, -1)

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, students)
}

func (c StudentController) Create(ctx *fiber.Ctx) error {
	var payload request.CreateStudentPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.StudentUseCase.CreateMany([]entity.Student{
		{
			Id:             payload.KmuttId,
			FirstName:      payload.FirstName,
			LastName:       payload.LastName,
			Email:          payload.Email,
			ProgrammeName:  payload.ProgrammeName,
			DepartmentName: payload.DepartmentName,
			GPAX:           *payload.GPAX,
			MathGPA:        *payload.MathGPA,
			EngGPA:         *payload.EngGPA,
			SciGPA:         *payload.SciGPA,
			School:         payload.School,
			City:           payload.City,
			Year:           payload.Year,
			Admission:      payload.Admission,
			Remark:         payload.Remark,
		},
	})
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c StudentController) CreateMany(ctx *fiber.Ctx) error {
	var payload request.CreateBulkStudentsPayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	newStudent := make([]entity.Student, 0, len(payload.Students))
	for _, student := range payload.Students {
		newStudent = append(newStudent, entity.Student{
			Id:             student.KmuttId,
			FirstName:      student.FirstName,
			LastName:       student.LastName,
			ProgrammeName:  student.ProgrammeName,
			DepartmentName: student.DepartmentName,
			GPAX:           *student.GPAX,
			MathGPA:        *student.MathGPA,
			EngGPA:         *student.EngGPA,
			SciGPA:         *student.SciGPA,
			School:         student.School,
			Year:           student.Year,
			Admission:      student.Admission,
			Remark:         student.Remark,
			City:           student.City,
		})
	}

	err := c.StudentUseCase.CreateMany(newStudent)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c StudentController) Update(ctx *fiber.Ctx) error {
	var payload request.UpdateStudentPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}
	fmt.Println("sadfadsfadfsafdsfdsadfs")

	id := ctx.Params("studentId")
	err := c.StudentUseCase.Update(id, &entity.Student{
		Id:             payload.KmuttId,
		FirstName:      payload.FirstName,
		LastName:       payload.LastName,
		ProgrammeName:  payload.ProgrammeName,
		DepartmentName: payload.DepartmentName,
		GPAX:           *payload.GPAX,
		MathGPA:        *payload.MathGPA,
		EngGPA:         *payload.EngGPA,
		SciGPA:         *payload.SciGPA,
		School:         payload.School,
		Year:           payload.Year,
		Admission:      payload.Admission,
		Remark:         *payload.Remark,
		City:           payload.City,
		Email:          payload.Email,
	})

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
