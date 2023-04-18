package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/schedule"
	"github.com/theandrew168/dripfile/internal/transfer"
)

type Scheduler struct {
	logger       *slog.Logger
	scheduleRepo schedule.Repository
	transferRepo transfer.Repository
}

func New(
	logger *slog.Logger,
	scheduleRepo schedule.Repository,
	transferRepo transfer.Repository,
) *Scheduler {
	s := Scheduler{
		logger:       logger,
		scheduleRepo: scheduleRepo,
		transferRepo: transferRepo,
	}
	return &s
}

func (s *Scheduler) Run(ctx context.Context) error {
	// main scheduler (handles sessions, transfers, etc)
	sched := gocron.NewScheduler(time.UTC)
	sched.WaitForScheduleAll()
	sched.TagsUnique()

	// run the scheduler in the background
	sched.StartAsync()

	// Load transfers at startup
	// Maintain set of currently scheduled transfers
	// Every minute, read transfers
	// For each transfer_id / tag, add or remove
	ticker := time.Tick(time.Minute)
	for {
		select {
		case <-ctx.Done():
			goto stop
		case <-ticker:
			err := s.reschedule(sched)
			if err != nil {
				s.logger.Error(err.Error())
				continue
			}
		}
	}

stop:
	s.logger.Info("stopping scheduler")
	sched.Stop()
	s.logger.Info("stopped scheduler")
	return nil
}

func (s *Scheduler) reschedule(sched *gocron.Scheduler) error {
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
	transfers, err := s.transferRepo.List()
	if err != nil {
		return err
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

		schedule, err := s.scheduleRepo.Read(transfer.ScheduleID)
		if err != nil {
			return err
		}

		s.logger.Info("schedule transfer",
			slog.String("transfer_id", transfer.ID),
		)

		sched.Cron(schedule.Expr).Tag(transfer.ID).Do(func() {
			// TODO: do the transfer
		})
	}

	// remove old transfers
	for id := range remove {
		s.logger.Info("unschedule transfer",
			slog.String("transfer_id", id),
		)

		err = sched.RemoveByTags(id)
		if err != nil {
			return err
		}
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
