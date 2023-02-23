package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func createRandomAccountParams() CreateAccountParams {
	return CreateAccountParams{
		Email:    RandomEmail(),
		Password: RandomString(10),
	}
}

func createRandomAccount() Account {
	p := createRandomAccountParams()
	a, _ := testQueries.CreateAccount(context.Background(), p)
	return *a
}

func TestCreateAccount(t *testing.T) {
	p := createRandomAccountParams()
	a, err := testQueries.CreateAccount(context.Background(), p)
	require.NoError(t, err)
	require.NotEmpty(t, a)
	require.Equal(t, p.Email, a.Email)
	require.Equal(t, p.Password, a.Password)
}
func TestDeleteAccount(t *testing.T) {
	p := createRandomAccount()
	err := testQueries.DeleteAccount(context.Background(), p.ID)
	require.NoError(t, err)
	d, err := testQueries.GetAccount(context.Background(), p.ID)
	require.Error(t, err)
	require.Empty(t, d)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount()
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
	a := createRandomAccount()

	args := UpdateAccountParams{
		ID:       a.ID,
		Email:    RandomEmail(),
		Password: RandomString(30),
	}

	b, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	require.Equal(t, b.Email, args.Email)
	require.Equal(t, b.Password, args.Password)
}
