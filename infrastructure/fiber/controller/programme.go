package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type ProgrammeController struct {
	ProgrammeUseCase entity.ProgrammeUseCase
	Validator        validator.PayloadValidator
}

func NewProgrammeController(validator validator.PayloadValidator, programmeUseCase entity.ProgrammeUseCase) *ProgrammeController {
	return &ProgrammeController{
		ProgrammeUseCase: programmeUseCase,
		Validator:        validator,
	}
}

func (c ProgrammeController) GetAll(ctx *fiber.Ctx) error {
	nameTH := ctx.Query("nameTH")
	nameEN := ctx.Query("nameEN")
	year := ctx.Query("year")

	if (nameTH != "" || nameEN != "") && year != "" {
		programme, err := c.ProgrammeUseCase.GetByNameAndYear(nameTH, nameEN, year)
		if err != nil {
			return err
		}

		if programme == nil {
			return response.NewSuccessResponse(ctx, fiber.StatusNotFound, programme)
		}

		return response.NewSuccessResponse(ctx, fiber.StatusOK, programme)
	} else if nameTH != "" || nameEN != "" {
		programme, err := c.ProgrammeUseCase.GetByName(nameTH, nameEN)
		if err != nil {
			return err
		}

		if programme == nil {
			return response.NewSuccessResponse(ctx, fiber.StatusNotFound, programme)
		}

		return response.NewSuccessResponse(ctx, fiber.StatusOK, programme)
	} else {
		programmes, err := c.ProgrammeUseCase.GetAll()
		if err != nil {
			return err
		}

		return response.NewSuccessResponse(ctx, fiber.StatusOK, programmes)
	}
}

func (c ProgrammeController) GetById(ctx *fiber.Ctx) error {
	programmeId := ctx.Params("programmeId")

	programme, err := c.ProgrammeUseCase.GetById(programmeId)
	if err != nil {
		return err
	}

	if programme == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, programme)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, programme)
}

// func (c ProgrammeController) GetByName(ctx *fiber.Ctx) error {
// 	nameTH := ctx.Query("nameTH")
// 	nameEN := ctx.Query("nameEN")

// 	programme, err := c.ProgrammeUseCase.GetByName(nameTH, nameEN)
// 	if err != nil {
// 		return err
// 	}

// 	if programme == nil {
// 		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, programme)
// 	}

// 	return response.NewSuccessResponse(ctx, fiber.StatusOK, programme)
// }

// func (c ProgrammeController) GetByNameAndYear(ctx *fiber.Ctx) error {
// 	programme, err := c.ProgrammeUseCase.GetByNameAndYear(nameTH, nameEN, year)
// 	if err != nil {
// 		return err
// 	}

// 	if programme == nil {
// 		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, programme)
// 	}

// 	return response.NewSuccessResponse(ctx, fiber.StatusOK, programme)
// }

func (c ProgrammeController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateProgrammePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.ProgrammeUseCase.Create(payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c ProgrammeController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateProgrammePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	name := ctx.Params("programmeName")

	err := c.ProgrammeUseCase.Update(name, &payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c ProgrammeController) Delete(ctx *fiber.Ctx) error {
	programmeId := ctx.Params("programmeId")

	err := c.ProgrammeUseCase.Delete(programmeId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c ProgrammeController) CreateLinkWithPO(ctx *fiber.Ctx) error {
	programmeId := ctx.Params("programmeId")
	payload := entity.LinkProgrammePO{}

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.ProgrammeUseCase.CreateLinkWithPO(programmeId, payload.POIds)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c ProgrammeController) CreateLinkWithPLO(ctx *fiber.Ctx) error {
	programmeId := ctx.Params("programmeId")
	payload := entity.LinkProgrammePLO{}

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.ProgrammeUseCase.CreateLinkWithPLO(programmeId, payload.PLOIds)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c ProgrammeController) CreateLinkWithSO(ctx *fiber.Ctx) error {
	programmeId := ctx.Params("programmeId")
	payload := entity.LinkProgrammeSO{}

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	err := c.ProgrammeUseCase.CreateLinkWithSO(programmeId, payload.SOIds)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c ProgrammeController) GetAllCourseOutcomeLinked(ctx *fiber.Ctx) error {
	programmeId := ctx.Params("programmeId")

	allCourseOutcome, err := c.ProgrammeUseCase.GetAllCourseOutcomeLinked(programmeId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, allCourseOutcome)
}

func (c ProgrammeController) GetAllCourseLinkedPO(ctx *fiber.Ctx) error {
	payload := entity.ProgrammeIds{}
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	allCourseLinkedPO, err := c.ProgrammeUseCase.GetAllCourseLinkedPO(payload.ProgrammeIds)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, allCourseLinkedPO)
}

func (c ProgrammeController) GetAllCourseLinkedPLO(ctx *fiber.Ctx) error {
	payload := entity.ProgrammeIds{}
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	allCourseLinkedPLO, err := c.ProgrammeUseCase.GetAllCourseLinkedPLO(payload.ProgrammeIds)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, allCourseLinkedPLO)
}

func (c ProgrammeController) GetAllCourseLinkedSO(ctx *fiber.Ctx) error {
	payload := entity.ProgrammeIds{}
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	allCourseLinkedSO, err := c.ProgrammeUseCase.GetAllCourseLinkedSO(payload.ProgrammeIds)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, allCourseLinkedSO)
}
