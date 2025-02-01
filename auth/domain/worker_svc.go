package domain

import (
	"context"
)

type TaskProducer interface {
	EnqueueSendMailTask(ctx context.Context, params SendMailParams) error
}

type TaskConsumer interface {
	RegisterHandlers()
	Run()
	GracefulStop(ctx context.Context) error
}
