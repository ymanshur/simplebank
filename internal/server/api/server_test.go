package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ymanshur/simplebank/config"
	"github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/pkg/util"
	"github.com/ymanshur/simplebank/pkg/worker"
)

func newTestServer(repo repo.Repo, taskDistributor worker.TaskDistributor) (*Server, error) {
	config := config.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	return NewServer(config, repo, taskDistributor)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
