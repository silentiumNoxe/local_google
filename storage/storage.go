package storage

type Storage interface {
	Write(id [32]byte, data []byte) error
}
