package tasks

import (
	"context"
	"os"

	"github.com/hibiken/asynq"
)

var Server *asynq.Server

func CreateServer() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := redisHost + ":" + redisPort

	Server = asynq.NewServer(asynq.RedisClientOpt{
		Addr: redisAddr,
	}, asynq.Config{})

	err := Server.Run(asynq.HandlerFunc(handler))

	return err
}

func handler(ctx context.Context, t *asynq.Task) error {
	var err error

	switch t.Type() {
	case AUTO_DELETE_URL_TASK:
		err = HandleAutoDeleteUrlTask(ctx, t)
	}

	return err
}
