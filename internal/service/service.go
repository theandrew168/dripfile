package service

import (
	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
)

type Service struct {
	logger *jsonlog.Logger
	store  *storage.Storage
	queue  *task.Queue
	box    *secret.Box
	mailer mail.Mailer
}

func New(
	logger *jsonlog.Logger,
	store *storage.Storage,
	queue *task.Queue,
	box *secret.Box,
	mailer mail.Mailer,
) *Service {
	s := Service{
		logger: logger,
		store:  store,
		queue:  queue,
		box:    box,
		mailer: mailer,
	}
	return &s
}
