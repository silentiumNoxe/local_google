package queue

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	cfg   *Config
	files map[string]*os.File
}

func NewLocalStorage(cfg *Config) *LocalStorage {
	if cfg.BaseDir == "" {
		cfg.BaseDir = "data"
	}

	return &LocalStorage{cfg: cfg, files: make(map[string]*os.File)}
}

func (l *LocalStorage) Put(entry *Entry) error {
	return l.write("00001", entry)
}

func (l *LocalStorage) write(id string, entry *Entry) error {
	file, err := l.open(id)
	if err != nil {
		return err
	}

	if len(entry.ID) == 0 {
		entry.ID = idgen([]byte(entry.Addr))
	}

	slog.Debug("(queue local) Store entry", slog.String("addr", entry.Addr), slog.String("entry", hex.EncodeToString(entry.ID)))

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	_, err = file.Write(entry.ID)
	if err != nil {
		return err
	}

	score := make([]byte, 4)
	binary.BigEndian.PutUint32(score, uint32(entry.Score))
	_, err = file.Write(score)
	if err != nil {
		return err
	}

	addr := make([]byte, 512)
	if len(entry.Addr) > 512 {
		slog.Warn("Address too long. Max 512", slog.String("addr", entry.Addr))
	}
	copy(addr, entry.Addr)

	_, err = file.Write(addr)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte("\n"))

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

func (l *LocalStorage) Pop(amount int) ([]*Entry, error) {
	//TODO implement me
	panic("implement me")
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
