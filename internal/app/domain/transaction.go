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
