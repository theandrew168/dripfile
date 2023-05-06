package transfer

type Transfer struct {
	// readonly (from database, after creation)
	ID string

	Pattern        string
	FromLocationID string
	ToLocationID   string
}

func New(pattern, fromLocationID, toLocationID string) Transfer {
	transfer := Transfer{
		Pattern:        pattern,
		FromLocationID: fromLocationID,
		ToLocationID:   toLocationID,
	}
	return transfer
}

type Repository interface {
	Create(transfer *Transfer) error
	Read(id string) (Transfer, error)
	List() ([]Transfer, error)
	ListByLocationID(locationID string) ([]Transfer, error)
	Update(transfer Transfer) error
	Delete(id string) error
}
