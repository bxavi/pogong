package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var testDb *sql.DB
var testQueries *Queries = New(testDb)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/postgres?sslmode=disable"
)

func TestMain(m *testing.M) {
	con, err := sql.Open(dbDriver, dbSource)
	if err != nil || con == nil {
		log.Fatal("TestMain no access to be db")
	}
	testQueries = New(con)
	os.Exit(m.Run())
}

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
func RandomBigInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
func RandomMoney() int64 {
	return rand.Int63n(RandomBigInt(1, 1000))
}
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(10))
}

func createRandomAccountParams() CreateAccountsParams {
	return CreateAccountsParams{
		Email:    RandomEmail(),
		Password: RandomString(10),
	}
}

func createRandomAccount() Account {
	p := createRandomAccountParams()
	a, _ := testQueries.CreateAccounts(context.Background(), p)
	return *a
}

func TestCreateAccount(t *testing.T) {
	p := createRandomAccountParams()
	a, err := testQueries.CreateAccounts(context.Background(), p)
	require.NoError(t, err)
	require.NotEmpty(t, a)
	require.Equal(t, p.Email, a.Email)
	require.Equal(t, p.Password, a.Password)
}
func TestDeleteAccount(t *testing.T) {
	p := createRandomAccount()
	err := testQueries.DeleteAccounts(context.Background(), p.ID)
	require.NoError(t, err)
	d, err := testQueries.GetAccounts(context.Background(), p.ID)
	require.Error(t, err)
	require.Empty(t, d)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount()
	}

	l, err := testQueries.ListAccounts(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, l)

	for _, a := range l {
		require.NotEmpty(t, a)
	}
}

func TestUpdateAccount(t *testing.T) {
	a := createRandomAccount()

	args := UpdateAccountsParams{
		ID:       a.ID,
		Email:    RandomEmail(),
		Password: RandomString(30),
	}

	b, err := testQueries.UpdateAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	require.Equal(t, b.Email, args.Email)
	require.Equal(t, b.Password, args.Password)
}
