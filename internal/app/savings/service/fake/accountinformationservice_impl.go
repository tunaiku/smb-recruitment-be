package fake

import (
	"math/big"

	"github.com/tunaiku/mobilebanking/internal/app/domain"
)

type FakeAccountInformationService struct {
}

func NewFakeAccountInformationService() *FakeAccountInformationService {
	return &FakeAccountInformationService{}
}

var trxDetails = map[string]domain.TransactionDetail{
	"T001": {
		Code:          "T001",
		MinimumAmount: big.NewFloat(2000.0),
	},
	"T002": {
		Code:          "T002",
		MinimumAmount: big.NewFloat(3000.0),
	},
}

var accountPrivileges = map[string][]string{
	"10001": {"T001", "T002"},
	"10002": {"T001"},
}

func (impl *FakeAccountInformationService) FindTransactionDetailByCode(code string) (domain.TransactionDetail, error) {
	trx := trxDetails[code]
	if (trx == domain.TransactionDetail{}) {
		return domain.TransactionDetail{}, domain.ErrTransactionDetailNotFound
	}
	return trx, nil
}

func (impl *FakeAccountInformationService) IsAccountExists(accountNumber string) bool {
	return accountPrivileges[accountNumber] != nil
}

func (impl *FakeAccountInformationService) GetTransactionPrivileges(accountNumber string) (domain.TransactionPrivileges, error) {
	if !impl.IsAccountExists(accountNumber) {
		return domain.TransactionPrivileges{}, domain.ErrAccountNotFound
	}
	return domain.TransactionPrivileges{Codes: accountPrivileges[accountNumber]}, nil
}
