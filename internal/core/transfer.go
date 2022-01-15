package core

type Transfer struct {
	// readonly (from database, after creation)
	ID int64

	Pattern string
	Src     Location
	Dst     Location
	Project Project
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
	Create(transfer *Transfer) error
	Read(id int64) (Transfer, error)
	Update(transfer Transfer) error
	Delete(transfer Transfer) error
}
