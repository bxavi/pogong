package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bxavi/pogong/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) *Account {
	password, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)
	require.NotEmpty(t, password)

	args := CreateAccountParams{
		Email:    util.RandomEmail(),
		Password: password,
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NotEmpty(t, account)
	require.NoError(t, err)
	require.Equal(t, account.Email, args.Email)
	require.Equal(t, account.Password, args.Password)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}
func TestDeleteAccount(t *testing.T) {
	p := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), p.ID)
	require.NoError(t, err)
	d, err := testQueries.GetAccount(context.Background(), p.ID)
	require.Error(t, err)
	require.Empty(t, d)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	l, err := testQueries.ListAccount(context.Background(), ListAccountParams{
		Limit:  sql.NullInt32{Valid: false},
		Offset: sql.NullInt32{Valid: false},
	})
	require.NoError(t, err)
	require.NotEmpty(t, l)

	for _, a := range l {
		require.NotEmpty(t, a)
	}
}

func TestUpdateAccount(t *testing.T) {
	a := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:       a.ID,
		Email:    util.RandomEmail(),
		Password: util.RandomString(30),
	}

	b, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	require.Equal(t, b.Email, args.Email)
	require.Equal(t, b.Password, args.Password)
}
