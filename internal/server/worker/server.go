package worker

import (
	"github.com/hibiken/asynq"
	"github.com/ymanshur/simplebank/config"
	"github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pkg/mail"
	"github.com/ymanshur/simplebank/pkg/worker"
)

type Server struct {
	server worker.TaskProcessor
}

func NewServer(config config.Config, repo repo.Repo, mailer mail.EmailSender, redisOpt asynq.RedisClientOpt) *Server {
	processor := worker.NewRedisTaskProcessor(redisOpt)

	sendVerifyEmailTask := NewTaskSendVerifyEmail(config, repo, mailer, ucase.NewVerifyEmailUseCase(repo))
	processor.RegisterTask(worker.Task{
		Name:    sendVerifyEmailTask.Name(),
		Handler: sendVerifyEmailTask.Handler,
	})

	return &Server{
		server: processor,
	}
}

func (s *Server) Start() error {
	return s.server.Start()
}

func (s *Server) Shutdown() {
	s.server.Shutdown()
}
