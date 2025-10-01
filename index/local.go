package index

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

type LocalStorage struct {
	cfg   *Config
	files map[string]*os.File
}

func NewLocalStorage(cfg *Config) *LocalStorage {
	if cfg.BaseDir == "" {
		cfg.BaseDir = "data"
	}

	if cfg.FileEntries == 0 {
		cfg.FileEntries = 100_000
	}

	return &LocalStorage{cfg: cfg, files: make(map[string]*os.File)}
}

func (l *LocalStorage) Save(entry *Entry) error {
	var entryId = l.getEntryID(entry.ID)
	var fileId = l.getFileID(entryId)

	file, err := l.open(fileId)
	if err != nil {
		return err
	}

	return nil
}

func (l *LocalStorage) getEntryID(b [32]byte) uint64 {
	var id uint64
	binary.BigEndian.Uint64(b[:7])
	return id
}

func (l *LocalStorage) getFileID(id uint64) string {
	var block = l.cfg.FileEntries
	var target = ((int(id)-1)/block + 1) * block

	return strconv.Itoa(target)
}

func (l *LocalStorage) readHead(f *os.File) error {
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}

	var endHeader = "</head>"
	var endHeaderLen = len(endHeader)

	var header []byte
	var buf = make([]byte, 1)
	for {
		_, err = f.Read(buf)
		if err != nil {
			return err
		}

		header = append(header, buf[0])
		end := header[(len(header)-endHeaderLen)-1:]

		if bytes.Equal([]byte(endHeader), end) {
			break
		}
	}

}

func (l *LocalStorage) FindById(id [32]byte) (*Entry, error) {
	var fileId = l.getFileID(id)

	return nil, nil
}

func (l *LocalStorage) write(id string, entry *Entry) error {
	//file, err := l.open(id)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (l *LocalStorage) open(id string) (*os.File, error) {
	if err := l.mkdir(); err != nil {
		return nil, err
	}

	fd := l.files[id]
	if fd != nil {
		return fd, nil
	}

	base := l.cfg.BaseDir
	path := filepath.Join(base, "queue")

	file, err := os.OpenFile(filepath.Join(path, id), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	l.files[id] = file
	return file, nil
}

func (l *LocalStorage) mkdir() error {
	base := l.cfg.BaseDir
	path := filepath.Join(base, "queue")
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *LocalStorage) Close() error {
	for name, fd := range l.files {
		err := fd.Close()
		if err != nil {
			slog.Warn("Unable to close file", slog.String("name", name), slog.String("err", err.Error()))
		}
	}
	return nil
}
