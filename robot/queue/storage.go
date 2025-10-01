package queue

type Storage interface {
	Put(*Entry) error
	Pop(amount int) ([]*Entry, error)
	Close() error
}
