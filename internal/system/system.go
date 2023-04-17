package system

import (
	"encoding/hex"

	"golang.org/x/exp/slog"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/service"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
)

type System struct {
	Logger  *slog.Logger
	Config  config.Config
	Box     *secret.Box
	Storage *storage.Storage
	Queue   *task.Queue
	Mailer  mail.Mailer
	Service *service.Service

	pool *pgxpool.Pool
}

func New(logger *slog.Logger, config config.Config) *System {
	s := System{
		Logger: logger,
		Config: config,
	}
	return &s
}

func (s *System) Start() error {
	secretKeyBytes, err := hex.DecodeString(s.Config.SecretKey)
	if err != nil {
		return err
	}

	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)
	s.Box = secret.NewBox(secretKey)

	s.pool, err = database.ConnectPool(s.Config.DatabaseURI)
	if err != nil {
		return err
	}

	s.Storage = storage.New(s.pool)
	s.Queue, err = task.NewQueue(s.pool)
	if err != nil {
		return err
	}

	if s.Config.SMTPURI != "" {
		s.Mailer, err = mail.NewSMTPMailer(s.Config.SMTPURI)
	} else {
		s.Logger.Info("using mock mailer")
		s.Mailer, err = mail.NewMockMailer(s.Logger)
	}
	if err != nil {
		return err
	}

	s.Service = service.New(s.Logger, s.Storage, s.Queue, s.Box, s.Mailer)

	return nil
}

func (s *System) Stop() {
	s.pool.Close()
}
