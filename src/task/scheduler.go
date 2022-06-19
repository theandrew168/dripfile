package task

import (
	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/src/config"
	"github.com/theandrew168/dripfile/src/jsonlog"
)

type Scheduler struct {
	cfg    config.Config
	logger *jsonlog.Logger
}

func NewScheduler(cfg config.Config, logger *jsonlog.Logger) *Scheduler {
	s := Scheduler{
		cfg:    cfg,
		logger: logger,
	}
	return &s
}

func (s *Scheduler) Run() error {
	redis, err := asynq.ParseRedisURI(s.cfg.RedisURI)
	if err != nil {
		return err
	}

	sched := asynq.NewScheduler(redis, nil)

	sessionPruneTask, err := NewSessionPruneTask()
	if err != nil {
		return err
	}

	_, err = sched.Register("0 * * * *", sessionPruneTask)
	if err != nil {
		return err
	}

	err = sched.Run()
	if err != nil {
		return err
	}

	return nil
}
