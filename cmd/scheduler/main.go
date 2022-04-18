package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/go-co-op/gocron"

	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/postgres"
	"github.com/theandrew168/dripfile/pkg/storage"
	"github.com/theandrew168/dripfile/pkg/task"
)

func main() {
	os.Exit(run())
}

func run() int {
	infoLog := log.New(os.Stdout, "", 0)
	errorLog := log.New(os.Stderr, "error: ", 0)

	// check for config file flag
	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		errorLog.Println(err)
		return 1
	}

	// open a database connection pool
	pool, err := postgres.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		errorLog.Println(err)
		return 1
	}
	defer pool.Close()

	store := storage.New(pool)
	queue := task.NewQueue(pool)

	// main scheduler (handles sessions, transfers, etc)
	s := gocron.NewScheduler(time.UTC)
	s.WaitForScheduleAll()
	s.TagsUnique()

	// prune sessions hourly
	s.Cron("0 * * * *").Do(func() {
		t, err := task.PruneSessions()
		if err != nil {
			errorLog.Println(err)
		}

		err = queue.Push(t)
		if err != nil {
			errorLog.Println(err)
		}
	})

	// run the scheduler in the background
	s.StartAsync()

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// Load transfers at startup
	// Maintain set of currently scheduled transfers
	// Every minute, read transfers
	// For each transfer_id / tag, add or remove
	c := time.Tick(time.Minute)
	for {
		// read tags of currently scheduled jobs
		have := make(map[string]bool)
		for _, j := range s.Jobs() {
			// skip untagged jobs
			tags := j.Tags()
			if len(tags) == 0 {
				continue
			}

			have[tags[0]] = true
		}

		// read all transfers from database
		transfers, err := store.Transfer.ReadAll()
		if err != nil {
			errorLog.Println(err)
			return 1
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

			infoLog.Printf("transfer %v add\n", transfer.ID)

			s.Cron(transfer.Schedule.Expr).Tag(transfer.ID).Do(func() {
				t, err := task.DoTransfer(transfer.ID)
				if err != nil {
					errorLog.Println(err)
					return
				}

				err = queue.Push(t)
				if err != nil {
					errorLog.Println(err)
					return
				}
			})
		}

		// remove old transfers
		for id, _ := range remove {
			infoLog.Printf("transfer %v remove\n", id)

			err = s.RemoveByTags(id)
			if err != nil {
				errorLog.Println(err)
				return 1
			}
		}

		<-c
	}

	return 0
}

func diff(have, want map[string]bool) (map[string]bool, map[string]bool) {
	// add = want but not have
	add := make(map[string]bool)
	for s, _ := range want {
		_, ok := have[s]
		if !ok {
			add[s] = true
		}
	}

	// remove = have but not want
	remove := make(map[string]bool)
	for s, _ := range have {
		_, ok := want[s]
		if !ok {
			remove[s] = true
		}
	}

	return add, remove
}
