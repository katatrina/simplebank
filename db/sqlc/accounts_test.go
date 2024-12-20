package db

import (
	"context"
	"testing"
	
	"github.com/katatrina/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	username := util.RandomOwner()
	hashedPassword, err := util.HashPassword("secret")
	
	arg := CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
		FullName:       username,
		Email:          username + "@gmail.com",
	}
	
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	
	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)
	
	return user
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	
	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 1)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	
	account2, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 1)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	
	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	
	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}
	
	arg := ListAccountsByOwnerParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}
	
	accounts, err := testStore.ListAccountsByOwner(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
