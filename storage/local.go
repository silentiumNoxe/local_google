package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	baseDir string
}

func NewLocal(baseDir string) *Local {
	return &Local{baseDir: baseDir}
}

func (l *Local) Write(id [32]byte, tid [2]byte, data []byte) error {
	const fragmentSize = 1024
	const headerSize = 40
	const payloadSize = fragmentSize - headerSize

	if len(data) > payloadSize {
		return errors.New("data too large")
	}

	if err := l.mkdir(l.baseDir); err != nil {
		return err
	}

	var frag = make([]byte, fragmentSize)

	copy(frag, tid[:])
	copy(frag[2:], id[:])
	copy(frag[headerSize-1:], data)

	return nil
}

func (l *Local) open(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return file, nil
		}

		if file != nil {
			_ = file.Close()
		}
		return nil, err
	}

	err = l.writeFileHeader(file)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	return file, nil
}

func (l *Local) writeFileHeader(w io.Writer) error {

}

func (l *Local) mkdir(base string) error {
	p := filepath.Join(base, "db")
	_, err := os.Stat(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return os.MkdirAll(p, 0644)
		}
	}

	return nil
}
