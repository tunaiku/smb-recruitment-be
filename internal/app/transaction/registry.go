package transaction

import (
	"log"

	"github.com/go-chi/chi"
	"github.com/go-pg/pg/v10"
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/handler"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/repository"
	"github.com/tunaiku/mobilebanking/internal/app/transaction/service"
	"go.uber.org/dig"
)

func Register(container *dig.Container) {

	container.Provide(func(db *pg.DB) domain.TransactionRepository {
		return repository.NewTransactionRepository(db)
	})
	container.Provide(func(repo domain.TransactionRepository) domain.MyTransactionService {
		return service.NewMyTransactionService(repo)
	})

	container.Provide(func(
		userSessionHelper domain.UserSessionHelper,
		transactionService domain.MyTransactionService,
		accountInformationService domain.AccountInformationService,
		transactionInformationService domain.TransactionInformationService,
		optCredetialManager domain.OtpCredentialManager,
		pinCredentialManger domain.PinCredentialManager,
	) *handler.TransactionEndpoint {
		return handler.NewTransactionEndpoint(
			userSessionHelper,
			transactionService,
			accountInformationService,
			transactionInformationService,
			optCredetialManager,
			pinCredentialManger,
		)
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
