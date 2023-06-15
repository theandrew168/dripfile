package transfer

type Storage interface {
	Create(t *Transfer) error
	Read(id string) (*Transfer, error)
	List() ([]*Transfer, error)
	Delete(id string) error
}
