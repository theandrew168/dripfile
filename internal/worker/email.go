package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/internal/task"
)

func (w *Worker) HandleEmailSend(ctx context.Context, t *asynq.Task) error {
	var info task.EmailSendInfo
	err := json.Unmarshal(t.Payload(), &info)
	if err != nil {
		return err
	}

	err = w.mailer.SendEmail(
		info.FromName,
		info.FromEmail,
		info.ToName,
		info.ToEmail,
		info.Subject,
		info.Body,
	)
	if err != nil {
		return err
	}

	return nil
}
