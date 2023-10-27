package service

import (
	"log/slog"

	"github.com/theandrew168/dripfile/backend/repository"
)

type Service struct {
	Transfer *Transfer
}

func New(logger *slog.Logger, repo *repository.Repository) *Service {
	srvc := Service{
		Transfer: NewTransfer(logger, repo),
	}
	return &srvc
}
