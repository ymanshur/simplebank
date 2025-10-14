package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ymanshur/simplebank/internal/common"
	"github.com/ymanshur/simplebank/pkg/token"
)

func addAuthorization(
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	role string,
	duration time.Duration,
) error {
	authToken, _, err := tokenMaker.CreateToken(username, role, duration)
	if err != nil {
		return err
	}

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, authToken)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
	return nil
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				err := addAuthorization(request, tokenMaker, authorizationTypeBearer, "user", common.DepositorRole, time.Minute)
				assert.NoError(t, err)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Expired Token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				err := addAuthorization(request, tokenMaker, authorizationTypeBearer, "user", common.DepositorRole, -time.Minute)
				assert.NoError(t, err)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server, err := newTestServer(nil, nil)
			assert.NoError(t, err)

			authPath := "/auth"
			server.router.GET(
				authPath,
				AuthMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)

			testCase.checkResponse(t, recorder)
		})
	}
}
