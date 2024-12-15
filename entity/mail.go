package entity

type MailUseCase interface {
	SendForgotPasswordEmail(to string) error
	ValidateResetPasswordToken(email, token string) error
	GetToken(email string) (string, error)
	DeleteToken(email string) error
}

type MailRepository interface {
	SaveToken(email string, token string) error
	GetToken(email string) (string, error)
	DeleteToken(email string) error
}
