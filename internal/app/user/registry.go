package user

import (
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/user/repository/inmemory"
	"github.com/tunaiku/mobilebanking/internal/app/user/service/fake"
	"go.uber.org/dig"
)

func Register(container *dig.Container) {
	container.Provide(func() domain.UserRepository {
		return inmemory.NewInMemoryUserRepository()
	})
	container.Provide(func(userRepository domain.UserRepository) domain.UserService {
		return fake.NewFakeUserService(userRepository)
	})

	container.Provide(func(userRepository domain.UserRepository) domain.OtpCredentialManager {
		return fake.NewFakeOtpCredentialManager(userRepository)
	})

	container.Provide(func(userRepository domain.UserRepository) domain.PinCredentialManager {
		return fake.NewFakePinCredentialManager(userRepository)
	})
}

func Invoke(container *dig.Container) {

}
