package repository

import (
	"time"

	"github.com/team-inu/inu-backyard/entity"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"github.com/team-inu/inu-backyard/internal/utils/session"
)

type MailRepository struct {
	Session session.SessionM
}

func NewMailRepository(session session.SessionM) entity.MailRepository {
	return &MailRepository{Session: session}
}

func (r MailRepository) GetToken(email string) (string, error) {
	token, valid := r.Session.GetSessionData(email)
	if !valid {
		return "", errs.New(errs.SameCode, "token not found")
	}

	return token, nil
}

func (r MailRepository) SaveToken(email string, token string) error {
	r.Session.CreateSession(email, time.Duration(time.Minute*15), token)
	return nil
}

func (r MailRepository) DeleteToken(email string) error {
	r.Session.RemoveSession(email)
	return nil
}
