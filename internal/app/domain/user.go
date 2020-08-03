package domain

import (
	"time"

	"github.com/micro/go-micro/v3/errors"
)

var (
	ErrCredentialNotMatch   error = errors.BadRequest("com.tunaiku.service.mbanking", "invalid credential")
	ErrOtpNotConfigured     error = errors.BadRequest("com.tunaiku.service.mbanking", "otp credential not configured on the user")
	ErrOtpAlreadyConfigured error = errors.BadRequest("com.tunaiku.service.mbanking", "otp credential not configured on the user")
	ErrPinNotConfigured     error = errors.BadRequest("com.tunaiku.service.mbanking", "pin credential not nonfigured on the user")
	ErrUserNotFound         error = errors.BadRequest("com.tunaiku.service.mbanking", "user not found")
)

//OtpCredential Represent user's otp credential
type OtpCredential struct {
	PhoneNumber string
}

//PinCredential Represent user's pin credential
type PinCredential struct {
	Pin string
}

//ConfiguredCredential Represent user credential
type ConfiguredCredential struct {
	Pin *PinCredential
	Otp *OtpCredential
}

//IsPinConfigured it would return true if user configure pin
func (c *ConfiguredCredential) IsPinConfigured() bool {
	return c.Pin != nil
}

//IsOtpConfigured it would return true if user configure otp
func (c *ConfiguredCredential) IsOtpConfigured() bool {
	return c.Otp != nil
}

type User struct {
	ID                              string
	Name                            string
	AccountReference                string
	JoinDate                        time.Time
	Username                        string
	Password                        string
	ConfiguredTransactionCredential *ConfiguredCredential
}

type FindUserResult struct {
	ID               string
	Name             string
	AccountReference string
	JoinAt           time.Time
}

type UserRepository interface {
	LoadUser(id string) (*User, error)
	LoadByUsername(username string) (*User, error)
}

type UserService interface {
	FindUser(userId string) (FindUserResult, error)
}

type UserCredentialValidator interface {
	Validate(userId string, credential string) error
}

type PinCredentialManager interface {
	UserCredentialValidator
}

type OtpCredentialManager interface {
	UserCredentialValidator
	RequestNewOtp(userId string) error
}
