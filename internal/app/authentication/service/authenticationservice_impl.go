package service

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/v3/errors"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	authJwt "github.com/tunaiku/mobilebanking/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationServiceImpl struct {
	repository domain.UserRepository
}

func NewAuthenticationServiceImpl(repository domain.UserRepository) *AuthenticationServiceImpl {
	return &AuthenticationServiceImpl{repository: repository}
}

func (srv *AuthenticationServiceImpl) Authenticate(username string, password string) (domain.AuthenticationResult, error) {
	user, err := srv.repository.LoadByUsername(username)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			return domain.AuthenticationResult{}, errors.BadRequest("com.tunaiku.service.mbanking", "invalid credential")
		default:
			return domain.AuthenticationResult{}, errors.InternalServerError("com.tunaiku.service.mbanking", err.Error())
		}

	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return domain.AuthenticationResult{}, domain.ErrCredentialNotMatch
		default:
			return domain.AuthenticationResult{}, errors.InternalServerError("com.tunaiku.service.mbanking", err.Error())
		}
	}
	accessToken, err := mapToJwt(user)
	if err != nil {
		return domain.AuthenticationResult{}, errors.InternalServerError("com.tunaiku.service.mbanking", err.Error())
	}
	return domain.AuthenticationResult{AccessToken: accessToken}, nil
}

func mapToJwt(user *domain.User) (token string, err error) {
	token, err = authJwt.CreateTokenString(func() jwt.Claims {
		return jwt.StandardClaims{Subject: user.ID}
	})
	return
}
