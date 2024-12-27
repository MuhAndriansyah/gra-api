package worker

import (
	"backend-layout/internal/config"
	"backend-layout/internal/tasks"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/hibiken/asynq"
)

type RedisTaskProcessor struct {
	server *asynq.Server
}

func NewTaskProcessor() *RedisTaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)
	srv := asynq.NewServer(asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%d", config.LoadRedisConfig().Host, config.LoadRedisConfig().Port),
	}, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 10,
			"default":  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).
				Bytes("payload", task.Payload()).Msg("process task failed")
		}),
		Logger: logger,
	})

	return &RedisTaskProcessor{
		server: srv,
	}

}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TaskSendVerifyEmail, tasks.HandlerVerifyEmail)

	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}
