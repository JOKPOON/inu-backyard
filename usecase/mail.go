package usecase

import (
	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"github.com/team-inu/inu-backyard/internal/utils"

	"github.com/team-inu/inu-backyard/infrastructure/mail"
)

type MailUseCase struct {
	MailRepository entity.MailRepository
}

func NewMailUseCase(
	mailRepository entity.MailRepository,
) entity.MailUseCase {
	return &MailUseCase{
		MailRepository: mailRepository,
	}
}

func (u MailUseCase) SendForgotPasswordEmail(to string) error {
	otp := utils.GenerateRandomInt(6)
	err := u.MailRepository.SaveToken(to, otp)
	if err != nil {
		return errs.New(errs.SameCode, "cannot save token")
	}

	go mail.SentForgotPasswordMail(to, otp)

	return nil
}

func (u MailUseCase) ValidateResetPasswordToken(email, token string) error {
	userToken, err := u.MailRepository.GetToken(email)
	if err != nil {
		return errs.New(errs.SameCode, "token not found")
	}

	if userToken != token {
		return errs.New(errs.SameCode, "token not match")
	}

	return nil
}

func (u MailUseCase) GetToken(email string) (string, error) {
	return u.MailRepository.GetToken(email)
}

func (u MailUseCase) DeleteToken(email string) error {
	return u.MailRepository.DeleteToken(email)
}
