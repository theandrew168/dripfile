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

	projects, err := store.Project.ReadAll()
	if err != nil {
		errorLog.Println(err)
		return 1
	}

	// start a scheduler for each project
	var schedulers []*gocron.Scheduler
	for _, project := range projects {

		// read all transfers linked to this project
		transfers, err := store.Transfer.ReadAllByProject(project)
		if err != nil {
			errorLog.Println(err)
			return 1
		}

		// add each transfer for the scheduler
		scheduler := gocron.NewScheduler(time.UTC)
		for _, transfer := range transfers {
			scheduler.Cron(transfer.Schedule.Expr).Do(func() {
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

		// start the scheduler in a background goro
		scheduler.StartAsync()
		schedulers = append(schedulers, scheduler)
	}

	// prune sessions hourly
	primary := gocron.NewScheduler(time.UTC)
	primary.Cron("0 * * * *").Do(func() {
		t, err := task.PruneSessions()
		if err != nil {
			errorLog.Println(err)
		}

		err = queue.Push(t)
		if err != nil {
			errorLog.Println(err)
		}
	})

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// run the primary scheduler forever
	primary.StartBlocking()

	return 0
}
