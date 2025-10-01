package index

type Storage interface {
	Save(*Entry) error
	FindById(id [32]byte) (*Entry, error)
	Close() error
}
