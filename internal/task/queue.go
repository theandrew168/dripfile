package task

import (
	"github.com/hibiken/asynq"
)

type Queue struct {
	client *asynq.Client
}

func NewQueue(redisURL string) (*Queue, error) {
	opts, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, err
	}

	client := asynq.NewClient(opts)
	q := Queue{
		client: client,
	}
	return &q, nil
}

func (q *Queue) Push(t Task) error {
	asynqTask := asynq.NewTask(t.Kind, []byte(t.Info))
	_, err := q.client.Enqueue(asynqTask)
	if err != nil {
		return err
	}

	return nil
}
