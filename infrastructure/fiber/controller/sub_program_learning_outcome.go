package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	request "github.com/team-inu/inu-backyard/infrastructure/fiber/request"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type SubProgramLearningOutcomeController struct {
	ProgramLearningOutcomeUseCase entity.ProgramLearningOutcomeUseCase
	Validator                     validator.PayloadValidator
}

func NewSubProgramLearningOutcomeController(validator validator.PayloadValidator, programLearningOutcomeUseCase entity.ProgramLearningOutcomeUseCase) *SubProgramLearningOutcomeController {
	return &SubProgramLearningOutcomeController{
		ProgramLearningOutcomeUseCase: programLearningOutcomeUseCase,
		Validator:                     validator,
	}
}

func (c SubProgramLearningOutcomeController) GetAll(ctx *fiber.Ctx) error {
	splos, err := c.ProgramLearningOutcomeUseCase.GetAllSubPlo()
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, splos)
}

func (c SubProgramLearningOutcomeController) GetById(ctx *fiber.Ctx) error {
	sploId := ctx.Params("sploId")

	splo, err := c.ProgramLearningOutcomeUseCase.GetSubPLO(sploId)
	if err != nil {
		return err
	}

	if splo == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, splo)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, splo)
}

func (c SubProgramLearningOutcomeController) Create(ctx *fiber.Ctx) error {
	var payload request.CreateSubProgramLearningOutcomePayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	subPlos := make([]entity.CreateSubProgramLearningOutcomeDto, 0, len(payload.SubProgramLearningOutcomes))
	for _, subPlo := range payload.SubProgramLearningOutcomes {
		subPlos = append(subPlos, entity.CreateSubProgramLearningOutcomeDto(subPlo))
	}

	err := c.ProgramLearningOutcomeUseCase.CreateSubPLO(subPlos)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c SubProgramLearningOutcomeController) Update(ctx *fiber.Ctx) error {
	var payload request.UpdateSubProgramLearningOutcomePayload

	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	id := ctx.Params("sploId")

	err := c.ProgramLearningOutcomeUseCase.UpdateSubPLO(id, &entity.SubProgramLearningOutcome{
		Code:                     payload.SubProgramLearningOutcomes[0].Code,
		DescriptionThai:          payload.SubProgramLearningOutcomes[0].DescriptionThai,
		DescriptionEng:           *payload.SubProgramLearningOutcomes[0].DescriptionEng,
		ProgramLearningOutcomeId: payload.SubProgramLearningOutcomes[0].ProgramLearningOutcomeId,
	})

	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c SubProgramLearningOutcomeController) Delete(ctx *fiber.Ctx) error {
	sploId := ctx.Params("sploId")

	err := c.ProgramLearningOutcomeUseCase.DeleteSubPLO(sploId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
