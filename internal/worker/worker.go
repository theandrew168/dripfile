package worker

import (
	"github.com/coreos/go-systemd/daemon"
	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
)

type Worker struct {
	logger   *jsonlog.Logger
	store    *storage.Storage
	box      *secret.Box
	mailer   mail.Mailer
	redisURL string
}

func New(
	logger *jsonlog.Logger,
	store *storage.Storage,
	box *secret.Box,
	mailer mail.Mailer,
	redisURL string,
) *Worker {
	w := Worker{
		logger:   logger,
		store:    store,
		box:      box,
		mailer:   mailer,
		redisURL: redisURL,
	}
	return &w
}

// TODO: signals and stuff?
func (w *Worker) Run() error {
	opts, err := asynq.ParseRedisURI(w.redisURL)
	if err != nil {
		return err
	}

	srv := asynq.NewServer(opts, asynq.Config{Concurrency: 4})

	mux := asynq.NewServeMux()
	mux.HandleFunc(task.KindSessionPrune, w.HandleSessionPrune)
	mux.HandleFunc(task.KindEmailSend, w.HandleEmailSend)
	mux.HandleFunc(task.KindTransferTry, w.HandleTransferTry)

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	err = srv.Run(mux)
	if err != nil {
		return err
	}

	return nil
}
