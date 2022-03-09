package service

import (
	"github.com/tunaiku/mobilebanking/internal/app/domain"
)

type MyTransactionService struct {
	transactionRepository domain.TransactionRepository
}

func NewMyTransactionService(repo domain.TransactionRepository) *MyTransactionService {
	return &MyTransactionService{
		transactionRepository: repo,
	}
}

func (ts *MyTransactionService) CreateTransaction(tx *domain.Transaction) (string, error) {
	err := ts.transactionRepository.SaveTransaction(tx)
	return tx.ID, err
}

func (ts *MyTransactionService) ReadTransaction(id string) (*domain.Transaction, error) {
	return ts.transactionRepository.ReadTransaction(id)
}

func (ts *MyTransactionService) UpdateTransaction(tx *domain.Transaction) error {
	err := ts.transactionRepository.SaveTransaction(tx)
	return err
}
