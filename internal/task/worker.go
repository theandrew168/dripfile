package task

import (
	"github.com/coreos/go-systemd/daemon"
	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/stripe"
)

type Worker struct {
	cfg     config.Config
	logger  *jsonlog.Logger
	store   *storage.Storage
	box     *secret.Box
	mailer  mail.Mailer
	billing stripe.Billing
}

func NewWorker(
	cfg config.Config,
	logger *jsonlog.Logger,
	store *storage.Storage,
	box *secret.Box,
	mailer mail.Mailer,
	billing stripe.Billing,
) *Worker {
	w := Worker{
		cfg:     cfg,
		logger:  logger,
		store:   store,
		box:     box,
		mailer:  mailer,
		billing: billing,
	}
	return &w
}

// TODO: signals and stuff?
func (w *Worker) Run() error {
	redis, err := asynq.ParseRedisURI(w.cfg.RedisURI)
	if err != nil {
		return err
	}

	srv := asynq.NewServer(redis, asynq.Config{Concurrency: 10})

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeSessionPrune, w.HandleSessionPrune)
	mux.HandleFunc(TypeEmailSend, w.HandleEmailSend)
	mux.HandleFunc(TypeTransferTry, w.HandleTransferTry)

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	err = srv.Run(mux)
	if err != nil {
		return err
	}

	return nil
}
