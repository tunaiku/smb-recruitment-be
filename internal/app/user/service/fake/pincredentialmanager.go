package fake

import "github.com/tunaiku/mobilebanking/internal/app/domain"

const (
	DefaultPin string = "123456"
)

type FakePinCredentialManager struct {
	repository domain.UserRepository
}

func NewFakePinCredentialManager(repository domain.UserRepository) *FakePinCredentialManager {
	return &FakePinCredentialManager{repository: repository}
}

func (fake *FakePinCredentialManager) Validate(userId string, credential string) error {
	user, err := fake.repository.LoadUser(userId)
	if err != nil {
		return err
	}
	if !user.ConfiguredTransactionCredential.IsPinConfigured() {
		return domain.ErrPinNotConfigured
	}
	if credential != DefaultPin {
		return domain.ErrCredentialNotMatch
	}
	return nil
}
