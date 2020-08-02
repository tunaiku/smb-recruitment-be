package domain

import (
	"context"

	"github.com/micro/go-micro/v3/errors"
)

var (
	ErrUnauthorized = errors.Unauthorized("com.tunaiku.service.mbanking", "invalid credential")
)

type AuthenticationResult struct {
	AccessToken string `json:"access_token"`
}

type AuthenticationService interface {
	Authenticate(username string, password string) (AuthenticationResult, error)
}

type UserSession struct {
	*User
}

type UserSessionHelper interface {
	GetFromContext(ctx context.Context) (session UserSession, err error)
}
