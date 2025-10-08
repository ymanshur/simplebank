package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ymanshur/simplebank/internal/common"
	db "github.com/ymanshur/simplebank/internal/repo"
	mockdb "github.com/ymanshur/simplebank/internal/repomock"
	"github.com/ymanshur/simplebank/pkg/util"
	"github.com/ymanshur/simplebank/pkg/worker"
	mockworker "github.com/ymanshur/simplebank/pkg/worker/mock"
	pb "github.com/ymanshur/simplebank/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// eqCreateUserTxParamsMatcher is used to represent the valid or expected arguments to createUser method.
type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	expected.arg.HashedPassword = arg.HashedPassword
	// there is no way to compare between functions
	if !reflect.DeepEqual(expected.arg.CreateUserParams, arg.CreateUserParams) {
		return false
	}

	err = arg.AfterCreate(expected.user)
	return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg.CreateUserParams, e.password)
}

// EqCreateUserTxParams returns a matcher that matches on CreateUserTxParams equality.
func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func TestRPC_CreateUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username:       user.Username,
						FullName:       user.FullName,
						Email:          user.Email,
						HashedPassword: password,
					},
				}

				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				requireResMatchUser(t, res, user)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()
			taskDistributor := mockworker.NewMockTaskDistributor(taskCtrl)

			tc.buildStubs(store, taskDistributor)

			// start test server and pass the request
			server := newTestServer(t, store, taskDistributor)
			res, err := server.CreateUser(context.Background(), tc.req)

			// check response
			// t has been shadowed by the t.Run arg
			// so the check response call for each case will be independent
			tc.checkResponse(t, res, err)
		})
	}
}

// randomUser generate random user instance
func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		Role:           common.DepositorRole,
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

// requireResMatchUser user response match assertions
func requireResMatchUser(t *testing.T, res *pb.CreateUserResponse, user db.User) {
	createdUser := res.GetUser()
	require.Equal(t, user.Username, createdUser.Username)
	require.Equal(t, user.FullName, createdUser.FullName)
	require.Equal(t, user.Email, createdUser.Email)
}
