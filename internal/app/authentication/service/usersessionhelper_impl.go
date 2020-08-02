package service

import (
	"context"

	"github.com/go-chi/jwtauth"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
)

type UserSessionHelperImpl struct {
	userRepository domain.UserRepository
}

func NewUserSessionHelperImpl(userRepository domain.UserRepository) *UserSessionHelperImpl {
	return &UserSessionHelperImpl{userRepository: userRepository}
}

func (helper UserSessionHelperImpl) GetFromContext(ctx context.Context) (session domain.UserSession, err error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return domain.UserSession{}, err
	}
	userID := claims["sub"].(string)
	user, err := helper.userRepository.LoadUser(userID)
	if err != nil {
		return domain.UserSession{}, err
	}
	return domain.UserSession{User: user}, err
}
