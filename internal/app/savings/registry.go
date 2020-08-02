package savings

import (
	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/savings/service/fake"
	"go.uber.org/dig"
)

func Register(container *dig.Container) {
	container.Provide(func() domain.AccountInformationService {
		return fake.NewFakeAccountInformationService()
	})
	container.Provide(func() domain.TransactionService {
		return fake.NewFakeTransactionService()
	})
}

func Invoke(container *dig.Container) {

}
