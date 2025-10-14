package worker

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	RegisterTask(task Task)
}

type RedisTaskProcessor struct {
	server       *asynq.Server
	taskHandlers map[string]asynq.Handler
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt) TaskProcessor {
	logger := NewLogger()

	redis.SetLogger(logger)

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("process task failed")
			}),
			Logger: logger,
		},
	)

	return &RedisTaskProcessor{
		server: server,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	for taskName, taskHandler := range p.taskHandlers {
		mux.HandleFunc(taskName, taskHandler.ProcessTask)
	}

	return p.server.Start(mux)
}

func (p *RedisTaskProcessor) Shutdown() {
	p.server.Shutdown()
}

func (p *RedisTaskProcessor) RegisterTask(task Task) {
	p.taskHandlers[task.Name] = asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		return task.Handler(ctx, t.Payload())
	})
}
