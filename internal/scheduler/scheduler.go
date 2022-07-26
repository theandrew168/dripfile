package scheduler

import (
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/go-co-op/gocron"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
)

type Scheduler struct {
	cfg    config.Config
	logger *jsonlog.Logger
	store  *storage.Storage
	queue  *task.Queue
}

func New(
	cfg config.Config,
	logger *jsonlog.Logger,
	store *storage.Storage,
	queue *task.Queue,
) *Scheduler {
	s := Scheduler{
		cfg:    cfg,
		logger: logger,
		store:  store,
		queue:  queue,
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
		t := task.NewSessionPruneTask()
		err := s.queue.Submit(t)
		if err != nil {
			s.logger.Error(err, nil)
			return
		}
	})

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// run the scheduler in the background
	sched.StartAsync()

	// Load transfers at startup
	// Maintain set of currently scheduled transfers
	// Every minute, read transfers
	// For each transfer_id / tag, add or remove
	c := time.Tick(time.Minute)
	for {
		// read tags of currently scheduled jobs
		have := make(map[string]bool)
		for _, j := range sched.Jobs() {
			// skip untagged jobs
			tags := j.Tags()
			if len(tags) == 0 {
				continue
			}

			have[tags[0]] = true
		}

		// read all transfers from database
		transfers, err := s.store.Transfer.ReadAll()
		if err != nil {
			s.logger.Error(err, nil)
			continue
		}

		// read tags of transfers in database
		want := make(map[string]bool)
		for _, t := range transfers {
			want[t.ID] = true
		}

		// diff scheduled transfers vs transfers in the database
		add, remove := diff(have, want)

		// add missing transfers
		for _, transfer := range transfers {
			if _, ok := add[transfer.ID]; !ok {
				continue
			}

			s.logger.Info("schedule transfer", map[string]string{
				"transfer_id": transfer.ID,
			})

			sched.Cron(transfer.Schedule.Expr).Tag(transfer.ID).Do(func() {
				t := task.NewTransferTryTask(transfer.ID)
				err = s.queue.Submit(t)
				if err != nil {
					s.logger.Error(err, nil)
					return
				}
			})
		}

		// remove old transfers
		for id := range remove {
			s.logger.Info("unschedule transfer", map[string]string{
				"transfer_id": id,
			})

			err = sched.RemoveByTags(id)
			if err != nil {
				s.logger.Error(err, nil)
				continue
			}
		}

		<-c
	}

	return nil
}

func diff(have, want map[string]bool) (map[string]bool, map[string]bool) {
	// add = want but not have
	add := make(map[string]bool)
	for s := range want {
		if _, ok := have[s]; !ok {
			add[s] = true
		}
	}

	// remove = have but not want
	remove := make(map[string]bool)
	for s := range have {
		if _, ok := want[s]; !ok {
			remove[s] = true
		}
	}

	return add, remove
}
