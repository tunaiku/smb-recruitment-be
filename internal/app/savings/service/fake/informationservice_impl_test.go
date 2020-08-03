package fake_test

import (
	"testing"

	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"github.com/tunaiku/mobilebanking/internal/app/savings/service/fake"
)

func TestFindTransactionDetailByCode_Should_ReturnTransactionDetail_When_TransactionIsAvailabelOnTheSystem(t *testing.T) {
	service := fake.NewFakeTransactionInformationService()
	detail, err := service.FindTransactionDetailByCode("T001")
	if err != nil {
		t.Fatal(err)
	}
	if detail.Code != "T001" && detail.MinimumAmount == nil {
		t.Fatal("transaction code should be `T001` and the minimun amount shouldn't be nil")
	}
}

func TestFindTransactionDetailByCode_Should_ReturnErrTransactionDetailNotFound_When_TheTransactionCodeIsInvalid(t *testing.T) {
	service := fake.NewFakeTransactionInformationService()
	_, err := service.FindTransactionDetailByCode("T003")
	if err == nil || err != domain.ErrTransactionDetailNotFound {
		t.Fatal("err should be `domain.ErrTransactionDetailNotFound`")
	}
}

func TestIsAccountExists_Should_ReturnTrue_When_TheAccountIsAvailableOnTheSystem(t *testing.T) {
	service := fake.NewFakeAccountInformationService()
	isExists := service.IsAccountExists("10001")
	if !isExists {
		t.Fatal("account should be exists")
	}
}

func TestIsAccountExists_Should_ReturnFalse_When_TheAccountNumberIsInvalid(t *testing.T) {
	service := fake.NewFakeAccountInformationService()
	isExists := service.IsAccountExists("10003")
	if isExists {
		t.Fatal("account shouldn't be exists")
	}
}

func TestGetTransactionPrivileges_ShouldReturn_TwoTransactionCode_When_TheAccountIsAvailableOnTheSystem(t *testing.T) {
	service := fake.NewFakeAccountInformationService()
	privileges, err := service.GetTransactionPrivileges("10001")
	if err != nil {
		t.Fatal(err)
	}

	if len(privileges.Codes) != 2 {
		t.Fatal("the length of allowed transaction should be 2")
	}
}

func TestGetTransactionPrivileges_Should_ReturnErrAccountNotFound_When_TheAccountNotIsInvalid(t *testing.T) {
	service := fake.NewFakeAccountInformationService()
	_, err := service.GetTransactionPrivileges("10003")
	if err == nil || err != domain.ErrAccountNotFound {
		t.Fatal("err should be `domain.ErrAccountNotFound`")
	}

}
