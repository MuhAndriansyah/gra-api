package tasks

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail,
		opts ...asynq.Option) error
	Close() error
}

type RedisTaskDestributor struct {
	client *asynq.Client
}

func (r *RedisTaskDestributor) Close() error {
	return r.client.Close()
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)

	return &RedisTaskDestributor{
		client: client,
	}
}
