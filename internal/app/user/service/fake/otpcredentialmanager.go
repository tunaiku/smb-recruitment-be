package fake

import (
	"fmt"

	"github.com/tunaiku/mobilebanking/internal/app/domain"
)

const (
	DefaultOtp = "111111"
)

type FakeOtpCredentialManager struct {
	repository domain.UserRepository
}

func NewFakeOtpCredentialManager(repository domain.UserRepository) *FakeOtpCredentialManager {
	return &FakeOtpCredentialManager{repository: repository}
}

func (fake *FakeOtpCredentialManager) Validate(userId string, credential string) error {
	user, err := fake.repository.LoadUser(userId)
	if err != nil {
		return err
	}
	if !user.ConfiguredTransactionCredential.IsOtpConfigured() {
		return domain.ErrOtpNotConfigured
	}
	if credential != DefaultOtp {
		return domain.ErrCredentialNotMatch
	}
	return nil
}

func (fake *FakeOtpCredentialManager) RequestNewOtp(userId string) error {
	user, err := fake.repository.LoadUser(userId)
	if err != nil {
		return err
	}
	if !user.ConfiguredTransactionCredential.IsOtpConfigured() {
		return domain.ErrOtpNotConfigured
	}
	fmt.Println("generating otp for userId ", userId)
	fmt.Println("your otp is ", DefaultOtp)
	return nil
}
