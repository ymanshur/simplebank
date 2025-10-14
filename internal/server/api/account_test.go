package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ymanshur/simplebank/internal/common"
	"github.com/ymanshur/simplebank/internal/repo"
	mockrepo "github.com/ymanshur/simplebank/internal/repo/mock"
	"github.com/ymanshur/simplebank/pkg/token"
	"github.com/ymanshur/simplebank/pkg/util"
)

func TestServer_GetAccount(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	// define table driven test sets
	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(mrepo *mockrepo.MockRepo)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK", // happy case
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},
			buildStubs: func(mrepo *mockrepo.MockRepo) {
				mrepo.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "Not Found",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},
			buildStubs: func(mrepo *mockrepo.MockRepo) {
				mrepo.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(repo.Account{}, repo.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal Error",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},
			buildStubs: func(mrepo *mockrepo.MockRepo) {
				mrepo.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(repo.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Invalid ID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},
			buildStubs: func(mrepo *mockrepo.MockRepo) {
				mrepo.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockrepo.NewMockRepo(ctrl)
			// build stubs
			testCase.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store, nil)

			recorder := httptest.NewRecorder()

			requestTarget := fmt.Sprintf("/v1/accounts/%d", testCase.accountID)
			request := httptest.NewRequest(http.MethodGet, requestTarget, nil)

			testCase.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			// check response
			testCase.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) repo.Account {
	return repo.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: common.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account repo.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount repo.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
