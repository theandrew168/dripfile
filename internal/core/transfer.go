package core

import (
	"context"
)

type Transfer struct {
	Pattern string
	Src     Location
	Dst     Location
	Project Project

	// readonly (from database, after creation)
	ID int
}

func NewTransfer(pattern string, src, dst Location, project Project) Transfer {
	transfer := Transfer{
		Pattern: pattern,
		Src:     src,
		Dst:     dst,
		Project: project,
	}
	return transfer
}

type TransferStorage interface {
	CreateTransfer(ctx context.Context, transfer *Transfer) error
	ReadTransfer(ctx context.Context, id int) (Transfer, error)
	UpdateTransfer(ctx context.Context, transfer Transfer) error
	DeleteTransfer(ctx context.Context, transfer Transfer) error
}
