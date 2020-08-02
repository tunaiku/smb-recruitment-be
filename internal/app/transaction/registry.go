package transaction

import (
	"log"

	"github.com/go-chi/chi"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/handler"
	"go.uber.org/dig"
)

func Register(container *dig.Container) {
	container.Provide(func(userSessionHelper domain.UserSessionHelper) *handler.TransactionEndpoint {
		return handler.NewTransactionEndpoint(userSessionHelper)
	})
}

func Invoke(container *dig.Container) {
	err := container.Invoke(func(router chi.Router, endpoint *handler.TransactionEndpoint) {
		log.Println("invoke transaction startup ...")
		endpoint.BindRoutes(router)
	})
	if err != nil {
		log.Fatal(err)
	}
}
