package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	db "github.com/bxavi/pogong/db"
	mockdb "github.com/bxavi/pogong/db/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/bxavi/pogong/util"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateAccountParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateAccountParams) // checks that x is correct type
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.Password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%v)", e.arg, e.password)
}

func EqCreateUserParamsMatcher(arg db.CreateAccountParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateAccountApi(t *testing.T) {

	account, password := randomAccount(t) // returns account with hashed password and original random password (unhashed)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    account.Email,
				"password": account.Password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Email:    account.Email,
					Password: account.Password,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), EqCreateUserParamsMatcher(arg, password)).
					Times(1).
					Return(account, nil)
				// todo create session stub
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, *account)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"email":    account.Email,
				"password": account.Password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Duplicate",
			body: gin.H{
				"email":    account.Email,
				"password": account.Password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&db.Account{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"email":    "BadEmail",
				"password": account.Password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: gin.H{
				"email":    account.Email,
				"password": "short",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// marshalbody data to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomAccount(t *testing.T) (*db.Account, string) {
	password, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)
	return &db.Account{
		ID:        util.RandomBigInt(1, 1000),
		Email:     util.RandomEmail(),
		Password:  password,
		CreatedAt: time.Now().UTC(),
	}, password
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account.CreatedAt, gotAccount.CreatedAt)
	require.Equal(t, account.Email, gotAccount.Email)
	require.Equal(t, account.ID, gotAccount.ID)
}
