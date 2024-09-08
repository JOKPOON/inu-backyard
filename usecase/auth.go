package usecase

import (
	"github.com/gofiber/fiber/v2"
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"golang.org/x/crypto/bcrypt"
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

func (u authUseCase) SignIn(email string, password string, ipAddress string, userAgent string) (*fiber.Cookie, error) {

	user, err := u.userUserCase.GetByEmail(email)
	if err != nil {
		return nil, errs.New(errs.SameCode, "cannot get user data to sign in", err)
	} else if user == nil {
		return nil, errs.New(errs.ErrUserNotFound, "password or email is incorrect")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errs.New(errs.ErrUserPassword, "password or email is incorrect")
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

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errs.New(errs.ErrUserPassword, "old password is incorrect")
	}

	newBcryptPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errs.New(errs.ErrCreateUser, "cannot create user", err)
	}

	newHashedPassword := string(newBcryptPassword)

	u.userUserCase.Update(userId, &entity.User{
		Password: newHashedPassword,
	})

	return nil
}
