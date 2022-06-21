package scheduler

import (
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/go-co-op/gocron"
	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
)

type Scheduler struct {
	cfg     config.Config
	logger  *jsonlog.Logger
	storage *storage.Storage
	queue   *asynq.Client
}

func New(
	cfg config.Config,
	logger *jsonlog.Logger,
	storage *storage.Storage,
	queue *asynq.Client,
) *Scheduler {
	s := Scheduler{
		cfg:     cfg,
		logger:  logger,
		storage: storage,
		queue:   queue,
	}
	return &s
}

// TODO: handle signals
// TODO: common code for: systemd notify + run thing in FG + graceful shutdown?
func (s *Scheduler) Run() error {
	// main scheduler (handles sessions, transfers, etc)
	sched := gocron.NewScheduler(time.UTC)
	sched.WaitForScheduleAll()
	sched.TagsUnique()

	// prune sessions hourly
	sched.Cron("* * * * *").Do(func() {
		t, err := task.NewSessionPruneTask()
		if err != nil {
			s.logger.Error(err, nil)
			return
		}

		_, err = s.queue.Enqueue(t)
		if err != nil {
			s.logger.Error(err, nil)
		}
	})

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// run the scheduler in the background
	sched.StartBlocking()

	return nil
}
