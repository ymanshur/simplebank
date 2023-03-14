package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/ymanshur/simplebank/db/sqlc"
	util2 "github.com/ymanshur/simplebank/pkg/util"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util2.Config{
		TokenSymmetricKey:   util2.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
