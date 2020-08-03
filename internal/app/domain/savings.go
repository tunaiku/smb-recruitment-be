package domain

import (
	"math/big"
	"time"

	"github.com/micro/go-micro/v3/errors"
)

var (
	ErrTransactionDetailNotFound = errors.BadRequest("com.tunaiku.service.cbs", "transaction detail not found")
	ErrAccountNotFound           = errors.BadRequest("com.tunaiku.service.cbs", "account not found")
)

type TransactionDetail struct {
	Code          string
	MinimumAmount *big.Float
}

type TransactionPrivileges struct {
	Codes []string
}

type TransactionCreation struct {
	SourceAccount      string
	DestinationAccount string
	TransactionCode    string
	Amount             *big.Float
	Currency           string
	TransactionDate    *time.Time
}

type AccountInformationService interface {
	IsAccountExists(accountNumber string) bool
	GetTransactionPrivileges(accountNumber string) (TransactionPrivileges, error)
}

type TransactionInformationService interface {
	FindTransactionDetailByCode(code string) (TransactionDetail, error)
}

type TransactionService interface {
	CreateTransaction(transactionCreation TransactionCreation) error
}
