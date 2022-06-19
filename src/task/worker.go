package task

import (
	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/src/config"
	"github.com/theandrew168/dripfile/src/jsonlog"
	"github.com/theandrew168/dripfile/src/mail"
)

type Worker struct {
	cfg    config.Config
	logger *jsonlog.Logger
	mailer mail.Mailer
}

func NewWorker(cfg config.Config, logger *jsonlog.Logger, mailer mail.Mailer) *Worker {
	w := Worker{
		cfg:    cfg,
		logger: logger,
		mailer: mailer,
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
	mux.HandleFunc(TypeEmailSend, w.HandleEmailSend)

	err = srv.Run(mux)
	if err != nil {
		return err
	}

	return nil
}
