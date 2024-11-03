package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"github.com/team-inu/inu-backyard/internal/utils"
)

type authUseCase struct {
	sessionUseCase entity.SessionUseCase
	userUserCase   entity.UserUseCase
}

func NewAuthUseCase(
	sessionUseCase entity.SessionUseCase,
	userUseCase entity.UserUseCase,
) entity.AuthUseCase {
	return &authUseCase{
		sessionUseCase: sessionUseCase,
		userUserCase:   userUseCase,
	}
}

func (u authUseCase) Authenticate(header string) (*entity.User, error) {
	session, err := u.sessionUseCase.Validate(header)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot authenticate user", err)
	}

	user, err := u.userUserCase.GetBySessionId(session.Id)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get user to authenticate", err)
	}
	return user, nil
}

func (u authUseCase) SignIn(payload entity.SignInPayload, ipAddress string, userAgent string) (*fiber.Cookie, error) {
	user, err := u.userUserCase.GetByEmail(payload.Email)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get user data to sign in", err)
	} else if user == nil {
		return nil, errs.New(errs.ErrUserNotFound, "password or email is incorrect")
	}

	err = utils.CheckPassword(user.Password, payload.Password)
	if err != nil {
		return nil, err
	}

	cookie, err := u.sessionUseCase.Create(user.Id, ipAddress, userAgent)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot create session to sign in", err)
	}
	return cookie, nil
}

func (u authUseCase) SignOut(header string) (*fiber.Cookie, error) {
	session, err := u.sessionUseCase.Validate(header)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot validate session to sign out", err)
	}

	cookie, err := u.sessionUseCase.Destroy(session.Id)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot destroy session to sign out", err)
	}
	return cookie, nil
}

func (u authUseCase) ChangePassword(userId string, oldPassword string, newPassword string) error {
	user, err := u.userUserCase.GetById(userId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get user data to sign in", err)
	} else if user == nil {
		return errs.New(errs.ErrUserNotFound, "cannot find target user")
	}

	err = utils.CheckPassword(user.Password, oldPassword)
	if err != nil {
		return err
	}

	newHashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	u.userUserCase.Update(userId, &entity.User{
		Password: newHashedPassword,
	})

	return nil
}
