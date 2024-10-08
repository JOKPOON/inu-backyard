package entity

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Session struct {
	Id        string    `gorm:"primaryKey;type:char(255)"`
	UserId    string    `json:"userId" gorm:"column:user_id"`
	IpAddress string    `json:"ipAddress" db:"ip_address"`
	UserAgent string    `json:"userAgent" db:"user_agent"`
	ExpiredAt time.Time `json:"expiredAt" db:"expired_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`

	User User `gorm:"foreignKey:UserId"`
}

type SessionRepository interface {
	Create(session *Session) error
	Get(id string) (*Session, error)
	Delete(id string) error
	DeleteByUserId(userId string) error
	DeleteDuplicates(userId string, ipAddress string, userAgent string) error
}

type SessionUseCase interface {
	Sign(id string) string
	Unsign(header string) (string, error)
	Create(userId string, ipAddress string, userAgent string) (*fiber.Cookie, error)
	Get(header string) (*Session, error)
	Destroy(id string) (*fiber.Cookie, error)
	DestroyByUserId(userId string) (*fiber.Cookie, error)
	Validate(header string) (*Session, error)
}
