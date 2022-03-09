package handler

func (transactionEndpoint *TransactionEndpoint) IfTransactionIsAllowedForUser(userAcc, transactionCode string) (bool, error) {
	codes, err := transactionEndpoint.accountInformationService.GetTransactionPrivileges(userAcc)
	if err != nil {
		return false, err
	}
	isTransactionValid := false
	for _, code := range codes.Codes {
		if code == transactionCode {
			isTransactionValid = true
			break
		}
	}

	return isTransactionValid, nil
}
