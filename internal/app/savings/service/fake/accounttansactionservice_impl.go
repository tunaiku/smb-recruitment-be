package fake

import (
	"fmt"

	"github.com/tunaiku/mobilebanking/internal/app/domain"
)

type FakeAccountTransactionService struct {
}

func NewFakeTransactionService() *FakeAccountTransactionService {
	return &FakeAccountTransactionService{}
}

func (fake *FakeAccountTransactionService) CreateTransaction(transactionCreation domain.TransactionCreation) error {
	fmt.Println("creating transaction...")
	return nil
}
