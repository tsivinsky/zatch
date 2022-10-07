package tasks

import (
	"os"

	"github.com/hibiken/asynq"
)

var Asynq *asynq.Client

func Connect() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := redisHost + ":" + redisPort

	Asynq = asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	return nil
}
