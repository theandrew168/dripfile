package task

import (
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
	storage *storage.Storage
	box     *secret.Box
	mailer  mail.Mailer
	billing stripe.Billing
}

func NewWorker(
	cfg config.Config,
	logger *jsonlog.Logger,
	storage *storage.Storage,
	box *secret.Box,
	mailer mail.Mailer,
	billing stripe.Billing,
) *Worker {
	w := Worker{
		cfg:     cfg,
		logger:  logger,
		storage: storage,
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

	err = srv.Run(mux)
	if err != nil {
		return err
	}

	return nil
}
