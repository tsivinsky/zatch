package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"url-shortener/pkg/db"

	"github.com/hibiken/asynq"
)

const (
	AUTO_DELETE_URL_TASK = "urls:auto_delete"
)

type AutoDeleteUrlPayload struct {
	UrlId uint
}

func NewAutoDeleteUrlTask(urlId uint) (*asynq.Task, error) {
	payload, err := json.Marshal(AutoDeleteUrlPayload{
		UrlId: urlId,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(AUTO_DELETE_URL_TASK, payload), nil
}

func HandleAutoDeleteUrlTask(ctx context.Context, t *asynq.Task) error {
	var payload AutoDeleteUrlPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	err := db.Db.Delete(db.Url{}, "id", payload.UrlId).Error

	return err
}
