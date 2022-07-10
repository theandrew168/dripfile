package task

import (
	"github.com/coreos/go-systemd/daemon"
	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
)

type Worker struct {
	logger *jsonlog.Logger
	store  *storage.Storage
	box    *secret.Box
	mailer mail.Mailer

	asynqServer *asynq.Server
}

func NewWorker(
	logger *jsonlog.Logger,
	store *storage.Storage,
	box *secret.Box,
	mailer mail.Mailer,
	asynqServer *asynq.Server,
) *Worker {
	w := Worker{
		logger: logger,
		store:  store,
		box:    box,
		mailer: mailer,

		asynqServer: asynqServer,
	}
	return &w
}

// TODO: signals and stuff?
func (w *Worker) Run() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeSessionPrune, w.HandleSessionPrune)
	mux.HandleFunc(TypeEmailSend, w.HandleEmailSend)
	mux.HandleFunc(TypeTransferTry, w.HandleTransferTry)

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)

	err := w.asynqServer.Run(mux)
	if err != nil {
		return err
	}

	return nil
}
