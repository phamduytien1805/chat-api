package taskq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
	"github.com/phamduytien1805/auth/domain"
	"github.com/phamduytien1805/package/config"
)

type TaskConsumerImpl struct {
	server  *asynq.Server
	mux     *asynq.ServeMux
	mailSvc domain.MailService
	logger  *slog.Logger
}

func NewTaskConsumer(config *config.RedisConfig, logger *slog.Logger, mailSvc domain.MailService) domain.TaskConsumer {
	server := asynq.NewServer(asynq.RedisClientOpt{Addr: config.Addr, DB: config.DB}, asynq.Config{
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.Error("error processing task", "type", task.Type, "payload", task.Payload, "error", err)
		}),
	})
	mux := asynq.NewServeMux()
	return &TaskConsumerImpl{
		server:  server,
		mux:     mux,
		mailSvc: mailSvc,
		logger:  logger,
	}
}

func (tc *TaskConsumerImpl) RegisterHandlers() {
	tc.mux.HandleFunc(string(SEND_EMAIL_TASK), func(ctx context.Context, t *asynq.Task) error {
		var params domain.SendMailParams
		if err := json.Unmarshal(t.Payload(), &params); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
		}
		return tc.mailSvc.SendVerificationMail(ctx, params)
	})
}

func (tc *TaskConsumerImpl) Run() {
	go func() {
		err := tc.server.Run(tc.mux)
		tc.logger.Info("task consumer starting")
		if err != nil {
			tc.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (tc *TaskConsumerImpl) GracefulStop(ctx context.Context) error {
	tc.server.Shutdown()
	return nil
}
