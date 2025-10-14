package worker

import "context"

type TaskHandler func(ctx context.Context, payload []byte) error

type Task struct {
	Name    string
	Handler TaskHandler
}
