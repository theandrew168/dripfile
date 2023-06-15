package location

type Storage interface {
	Create(l *Location) error
	Read(id string) (*Location, error)
	List() ([]*Location, error)
	Delete(id string) error
}
