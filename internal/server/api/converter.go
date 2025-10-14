package api

import "github.com/ymanshur/simplebank/internal/ucase"

func convertUser(user *ucase.UserResponse) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func convertAccount(account *ucase.AccountResponse) AccountResponse {
	return AccountResponse{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

func convertAccounts(accounts []ucase.AccountResponse) []AccountResponse {
	var res []AccountResponse
	for _, account := range accounts {
		res = append(res, convertAccount(&account))
	}
	return res
}

func convertTransferResult(result *ucase.TransferResult) *TransferResult {
	return &TransferResult{
		Transfer: TransferResponse{
			ID:            result.Transfer.ID,
			FromAccountID: result.Transfer.FromAccountID,
			ToAccountID:   result.Transfer.ToAccountID,
			Amount:        result.Transfer.Amount,
			CreatedAt:     result.Transfer.CreatedAt,
		},
		FromAccount: AccountResponse{
			ID:        result.FromAccount.ID,
			Owner:     result.FromAccount.Owner,
			Balance:   result.FromAccount.Balance,
			Currency:  result.FromAccount.Currency,
			CreatedAt: result.FromAccount.CreatedAt,
		},
		ToAccount: AccountResponse{
			ID:        result.ToAccount.ID,
			Owner:     result.ToAccount.Owner,
			Balance:   result.ToAccount.Balance,
			Currency:  result.ToAccount.Currency,
			CreatedAt: result.ToAccount.CreatedAt,
		},
		FromEntry: EntryResponse{
			ID:        result.FromEntry.ID,
			AccountID: result.FromEntry.AccountID,
			Amount:    result.FromEntry.Amount,
			CreatedAt: result.FromEntry.CreatedAt,
		},
		ToEntry: EntryResponse{
			ID:        result.ToEntry.ID,
			AccountID: result.ToEntry.AccountID,
			Amount:    result.ToEntry.Amount,
			CreatedAt: result.ToEntry.CreatedAt,
		},
	}
}
