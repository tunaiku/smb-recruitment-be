package domain

import (
	"time"
)

type TransactionState int

const (
	UnknownTransactionStatus TransactionState = iota
	WaitAuthorization
	Failed
	Success
)

type AuthorizationMethod int

const (
	UnknownAuthorizationMethod AuthorizationMethod = iota
	OtpAuthorization
	PinAuthorization
)

type Transaction struct {
	ID                  string
	UserID              string
	State               TransactionState
	AuthorizationMethod AuthorizationMethod
	TransactionCode     string
	Amount              float64
	SourceAccount       string
	DestinationAccount  string
	CreatedAt           time.Time
}

type MyTransactionService interface {
	CreateTransaction(tx *Transaction) (string, error)
	ReadTransaction(id string) (*Transaction, error)
	UpdateTransaction(tx *Transaction) error
}

type TransactionRepository interface {
	SaveTransaction(tx *Transaction) error
	ReadTransaction(ID string) (*Transaction, error)
}
