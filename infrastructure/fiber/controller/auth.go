package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"github.com/team-inu/inu-backyard/infrastructure/captcha"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/response"
	"github.com/team-inu/inu-backyard/internal/config"
	"github.com/team-inu/inu-backyard/internal/validator"
)

type AuthController struct {
	Config      config.AuthConfig
	Validator   validator.PayloadValidator
	Turnstile   captcha.Validator
	AuthUseCase entity.AuthUseCase
	UserUseCase entity.UserUseCase
}

func NewAuthController(
	validator validator.PayloadValidator,
	config config.AuthConfig,
	turnstile captcha.Validator,
	authUseCase entity.AuthUseCase,
	userUseCase entity.UserUseCase,
) *AuthController {
	return &AuthController{
		Config:      config,
		Validator:   validator,
		Turnstile:   turnstile,
		AuthUseCase: authUseCase,
		UserUseCase: userUseCase,
	}
}

func (c AuthController) Me(ctx *fiber.Ctx) error {
	user, err := c.UserUseCase.GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user == nil {
		return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, errs.New(0, "cannot get user from context"))
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, user)
}

func (c AuthController) SignIn(ctx *fiber.Ctx) error {
	var payload entity.SignInPayload
	if ok, err := c.Validator.Validate(&payload, ctx); !ok {
		return err
	}

	ipAddress := ctx.IP()
	userAgent := string(ctx.Context().UserAgent())

	cfToken := string(ctx.Request().Header.Peek("Cf-Token")[:])

	isTokenValid, err := c.Turnstile.Validate(cfToken, ipAddress)
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, errs.New(0, "cannot validate challenge token"))
	} else if !isTokenValid {
		return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, errs.New(0, "invalid challenge token"))

	}

	cookie, err := c.AuthUseCase.SignIn(payload, ipAddress, userAgent)
	if err != nil {
		return err
	}

	ctx.Cookie(cookie)

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"expired_at": cookie.Expires,
	})
}

func (c AuthController) SignOut(ctx *fiber.Ctx) error {
	sid := ctx.Cookies(c.Config.Session.CookieName)
	cookie, err := c.AuthUseCase.SignOut(sid)
	if err != nil {
		return err
	}
	ctx.Cookie(cookie)

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"signout_at": time.Now(),
	})
}
