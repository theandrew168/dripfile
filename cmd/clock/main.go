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

	queue := task.NewQueue(pool)

	s := gocron.NewScheduler(time.UTC)

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

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	// run the scheduler forever
	s.StartBlocking()

	return 0
}
