package mocks

import "github.com/stretchr/testify/mock"

type Turnstile struct {
	secret string
	mock.Mock
}

func NewTurnstile(secret string) *Turnstile {
	return &Turnstile{
		secret: secret,
	}
}

func (m *Turnstile) Validate(token string, ip string) (bool, error) {
	args := m.Called(token, ip)
	return args.Bool(0), args.Error(1)
}
