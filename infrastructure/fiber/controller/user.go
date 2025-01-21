package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type UserController struct {
	UserUseCase entity.UserUseCase
	AuthUseCase entity.AuthUseCase
	Validator   validator.PayloadValidator
}

func NewUserController(validator validator.PayloadValidator, userUseCase entity.UserUseCase, authUseCase entity.AuthUseCase) *UserController {
	return &UserController{
		UserUseCase: userUseCase,
		AuthUseCase: authUseCase,
		Validator:   validator,
	}
}

func (c UserController) GetAll(ctx *fiber.Ctx) error {
	pageIndex := ctx.Query("pageIndex")
	pageSize := ctx.Query("pageSize")
	query := ctx.Query("query")

	users, err := c.UserUseCase.GetAll(query, pageIndex, pageSize)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, users)
}

func (c UserController) GetById(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	user, err := c.UserUseCase.GetById(userId)
	if err != nil {
		return err
	}

	if user == nil {
		return response.NewSuccessResponse(ctx, fiber.StatusNotFound, user)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, user)
}

func (c UserController) Create(ctx *fiber.Ctx) error {
	var payload entity.CreateUserPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	// user := middleware.GetUserFromCtx(ctx)
	// if !user.IsRoles([]entity.UserRole{entity.UserRoleHeadOfCurriculum}) {
	// 	return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, nil)
	// }

	err := c.UserUseCase.Create(payload)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c UserController) CreateMany(ctx *fiber.Ctx) error {
	var payload entity.CreateBulkUserPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	newUsers := make([]entity.User, 0, len(payload.Users))
	for _, user := range payload.Users {
		newUsers = append(newUsers, entity.User{
			TitleTHShort:       user.TitleTHShort,
			TitleENShort:       user.TitleENShort,
			TitleTH:            user.TitleTH,
			TitleEN:            user.TitleEN,
			FirstNameTH:        user.FirstNameTH,
			LastNameTH:         user.LastNameTH,
			FirstNameEN:        user.FirstNameEN,
			LastNameEN:         user.LastNameEN,
			Email:              user.Email,
			AcademicPositionTH: user.AcademicPositionTH,
			AcademicPositionEN: user.AcademicPositionEN,
			Role:               user.Role,
			DegreeTH:           user.DegreeTH,
			DegreeEN:           user.DegreeEN,
			Password:           user.Password,
			Tel:                user.Tel,
		})
	}

	err := c.UserUseCase.CreateMany(newUsers)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, nil)
}

func (c UserController) Update(ctx *fiber.Ctx) error {
	var payload entity.UpdateUserPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	targetUserId := ctx.Params("userId")

	// err := c.UserUseCase.CheckUserRole(ctx, targetUserId, entity.UserRoleHeadOfCurriculum)
	// if err != nil {
	// 	return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, nil)
	// }

	err := c.UserUseCase.Update(targetUserId, &entity.User{
		TitleTHShort:       payload.TitleTHShort,
		TitleENShort:       payload.TitleENShort,
		TitleTH:            payload.TitleTH,
		TitleEN:            payload.TitleEN,
		FirstNameTH:        payload.FirstNameTH,
		LastNameTH:         payload.LastNameTH,
		FirstNameEN:        payload.FirstNameEN,
		LastNameEN:         payload.LastNameEN,
		AcademicPositionTH: payload.AcademicPositionTH,
		AcademicPositionEN: payload.AcademicPositionEN,
		Email:              payload.Email,
		Role:               payload.Role,
		DegreeTH:           payload.DegreeTH,
		DegreeEN:           payload.DegreeEN,
		Tel:                payload.Tel,
	})
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c UserController) Delete(ctx *fiber.Ctx) error {
	targetUserId := ctx.Params("userId")

	err := c.UserUseCase.CheckUserRole(ctx, targetUserId, entity.UserRoleHeadOfCurriculum)
	if err != nil {
		return err
	}

	err = c.UserUseCase.Delete(targetUserId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c UserController) ChangePassword(ctx *fiber.Ctx) error {
	var payload entity.ChangePasswordPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	targetUserId := ctx.Params("userId")

	err := c.UserUseCase.CheckUserRole(ctx, targetUserId, entity.UserRoleHeadOfCurriculum)
	if err != nil {
		return err
	}

	err = c.AuthUseCase.ChangePassword(targetUserId, payload.OldPassword, payload.NewPassword)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
