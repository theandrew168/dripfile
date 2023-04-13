package scheduler

import (
	"context"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/go-co-op/gocron"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
)

type Scheduler struct {
	logger *slog.Logger
	store  *storage.Storage
	queue  *task.Queue
}

func New(
	logger *slog.Logger,
	store *storage.Storage,
	queue *task.Queue,
) *Scheduler {
	s := Scheduler{
		logger: logger,
		store:  store,
		queue:  queue,
	}
	return &s
}

func (s *Scheduler) Run(ctx context.Context) error {
	// main scheduler (handles sessions, transfers, etc)
	sched := gocron.NewScheduler(time.UTC)
	sched.WaitForScheduleAll()
	sched.TagsUnique()

	// prune sessions hourly
	sched.Cron("* * * * *").Do(func() {
		t := task.NewSessionPruneTask()
		err := s.queue.Submit(t)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
	})

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// run the scheduler in the background
	sched.StartAsync()
	<-ctx.Done()

	s.logger.Info("stopping scheduler")
	sched.Stop()
	s.logger.Info("stopped scheduler")
	return nil
}
