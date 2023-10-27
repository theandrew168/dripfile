package repository

import (
	"errors"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/memorydb"
)

// ensure TransferRepository interface is satisfied
var _ TransferRepository = (*MemoryTransferRepository)(nil)

type TransferRepository interface {
	Create(transfer *domain.Transfer) error
	List() ([]*domain.Transfer, error)
	Read(id uuid.UUID) (*domain.Transfer, error)
	Update(transfer *domain.Transfer) error
	Delete(id uuid.UUID) error
}

type MemoryTransferRepository struct {
	db *memorydb.MemoryDB[*domain.Transfer]
}

func NewMemoryTransferRepository() *MemoryTransferRepository {
	repo := MemoryTransferRepository{
		db: memorydb.New[*domain.Transfer](),
	}
	return &repo
}

func (repo *MemoryTransferRepository) Create(transfer *domain.Transfer) error {
	return repo.db.Create(transfer)
}

func (repo *MemoryTransferRepository) List() ([]*domain.Transfer, error) {
	return repo.db.List()
}

func (repo *MemoryTransferRepository) Read(id uuid.UUID) (*domain.Transfer, error) {
	transfer, err := repo.db.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, memorydb.ErrNotFound):
			return nil, ErrNotExist
		default:
			return nil, err
		}
	}

	return transfer, nil
}

func (repo *MemoryTransferRepository) Update(transfer *domain.Transfer) error {
	return repo.db.Update(transfer)
}

func (repo *MemoryTransferRepository) Delete(id uuid.UUID) error {
	err := repo.db.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, memorydb.ErrNotFound):
			return ErrNotExist
		default:
			return err
		}
	}

	return nil
}
