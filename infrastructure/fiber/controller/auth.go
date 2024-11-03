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
	config      config.AuthConfig
	validator   validator.PayloadValidator
	turnstile   captcha.Validator
	authUseCase entity.AuthUseCase
	userUseCase entity.UserUseCase
}

func NewAuthController(
	validator validator.PayloadValidator,
	config config.AuthConfig,
	turnstile captcha.Validator,
	authUseCase entity.AuthUseCase,
	userUseCase entity.UserUseCase,
) *AuthController {
	return &AuthController{
		config:      config,
		validator:   validator,
		turnstile:   turnstile,
		authUseCase: authUseCase,
		userUseCase: userUseCase,
	}
}

func (c AuthController) Me(ctx *fiber.Ctx) error {
	user, err := c.userUseCase.GetUserFromCtx(ctx)
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
	if ok, err := c.validator.Validate(&payload, ctx); !ok {
		return err
	}

	ipAddress := ctx.IP()
	userAgent := string(ctx.Context().UserAgent())

	cfToken := string(ctx.Request().Header.Peek("Cf-Token")[:])

	isTokenValid, err := c.turnstile.Validate(cfToken, ipAddress)
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, errs.New(0, "cannot validate challenge token"))
	} else if !isTokenValid {
		return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, errs.New(0, "invalid challenge token"))

	}

	cookie, err := c.authUseCase.SignIn(payload, ipAddress, userAgent)
	if err != nil {
		return err
	}

	ctx.Cookie(cookie)

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"expired_at": cookie.Expires,
	})
}

func (c AuthController) SignOut(ctx *fiber.Ctx) error {
	sid := ctx.Cookies(c.config.Session.CookieName)
	cookie, err := c.authUseCase.SignOut(sid)
	if err != nil {
		return err
	}
	ctx.Cookie(cookie)

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"signout_at": time.Now(),
	})
}
