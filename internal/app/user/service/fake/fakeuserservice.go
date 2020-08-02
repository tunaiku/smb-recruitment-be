package fake

import "github.com/tunaiku/mobilebanking/internal/app/domain"

type FakeUserService struct {
	repository domain.UserRepository
}

func NewFakeUserService(repository domain.UserRepository) *FakeUserService {
	return &FakeUserService{repository: repository}
}

func (fake *FakeUserService) FindUser(userId string) (domain.FindUserResult, error) {
	user, err := fake.repository.LoadUser(userId)
	if err != nil {
		return domain.FindUserResult{}, err
	}
	return domain.FindUserResult{
		AccountReference: user.AccountReference,
		ID:               user.ID,
		JoinAt:           user.JoinDate,
		Name:             user.Name,
	}, nil
}
