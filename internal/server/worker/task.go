package worker

import "context"

type Task interface {
	Name() string
	Handler(ctx context.Context, payload []byte) error
}
