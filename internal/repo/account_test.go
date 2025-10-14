package repo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/ymanshur/simplebank/db/sqlc"
	"github.com/ymanshur/simplebank/internal/common"
	"github.com/ymanshur/simplebank/pkg/util"
)

func createRandomAccount(t *testing.T) db.Account {
	user := createRandomUser(t)

	arg := db.CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: common.RandomCurrency(),
	}

	account, err := testRepo.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestRepo_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestRepo_GetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testRepo.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestRepo_UpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := db.UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testRepo.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestRepo_DeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testRepo.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testRepo.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, account2)
}

func TestRepo_ListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := db.ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testRepo.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
