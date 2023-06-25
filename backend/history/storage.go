package history

type Storage interface {
	Create(l *History) error
	Read(id string) (*History, error)
	List() ([]*History, error)
}
