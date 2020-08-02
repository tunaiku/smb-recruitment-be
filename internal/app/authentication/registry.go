package authentication

import (
	"log"

	"github.com/go-chi/chi"
	"github.com/tunaiku/mobilebanking/internal/app/authentication/handler"
	"github.com/tunaiku/mobilebanking/internal/app/authentication/service"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"go.uber.org/dig"
)

func Register(container *dig.Container) {
	container.Provide(func(userRepository domain.UserRepository) domain.AuthenticationService {
		return service.NewAuthenticationServiceImpl(userRepository)
	})

	container.Provide(func(authenticationService domain.AuthenticationService) *handler.AuthenticationEndpoint {
		return handler.NewAuthenticationEndpoint(authenticationService)
	})

	container.Provide(func(userRepository domain.UserRepository) domain.UserSessionHelper {
		return service.NewUserSessionHelperImpl(userRepository)
	})
}

func Invoke(container *dig.Container) {
	err := container.Invoke(func(router chi.Router, endpoint *handler.AuthenticationEndpoint) {
		log.Println("invoke authentication startup ...")
		endpoint.BindRoutes(router)
	})
	if err != nil {
		log.Fatalln(err)
	}
}
