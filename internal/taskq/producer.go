package taskq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/phamduytien1805/internal/platform/mail"
	"github.com/phamduytien1805/package/config"
)

type TaskProducer interface {
	EnqueueSendMailTask(ctx context.Context, payload mail.SendMailParams, opts ...asynq.Option) error
}

type TaskProducerImpl struct {
	client *asynq.Client
	logger *slog.Logger
}

func NewTaskProducer(config *config.RedisConfig, logger *slog.Logger) TaskProducer {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.Addr, DB: config.DB})
	return &TaskProducerImpl{
		client: client,
		logger: logger,
	}
}

func (t *TaskProducerImpl) EnqueueSendMailTask(ctx context.Context, payload mail.SendMailParams, opts ...asynq.Option) error {
	task, err := t.createTask(SEND_EMAIL_TASK, payload, opts...)
	if err != nil {
		return err
	}
	taskInfo, err := t.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	t.logger.Info("task enqueued", "task", taskInfo)
	return nil
}

func (t *TaskProducerImpl) createTask(taskKey DefinedTask, payload interface{}, opts ...asynq.Option) (*asynq.Task, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(string(taskKey), jsonPayload, opts...)
	return task, nil

}
