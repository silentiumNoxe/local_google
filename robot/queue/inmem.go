package queue

import (
	"encoding/hex"
	"log/slog"
	"sync"
)

type InMemStorage struct {
	cfg *Config
	s   Storage

	first *inmemEntry
	last  *inmemEntry
	m     *sync.Mutex
}

func NewInMemStorage(cfg *Config, delegate Storage) *InMemStorage {
	return &InMemStorage{cfg: cfg, s: delegate, m: &sync.Mutex{}}
}

func (i *InMemStorage) Put(entry *Entry) error {
	i.m.Lock()
	defer i.m.Unlock()

	if len(entry.ID) == 0 {
		entry.ID = idgen([]byte(entry.Addr))
	}

	slog.Debug("(queue in-mem) Store entry", slog.String("addr", entry.Addr), slog.String("entry", hex.EncodeToString(entry.ID)))

	e := &inmemEntry{data: entry}

	if i.first == nil {
		i.first = e
		i.last = i.first
		return nil
	}

	last := i.last
	i.last = e
	last.next = e

	return nil
}

func (i *InMemStorage) Pop(amount int) ([]*Entry, error) {
	i.m.Lock()
	defer i.m.Unlock()

	arr := make([]*Entry, amount)

	var x = i.first
	if x == nil {
		return arr, nil
	}

	for i := 0; i < amount; i++ {
		arr[i] = x.data
		x = x.next
	}

	return arr, nil
}

func (i *InMemStorage) Close() error {
	return i.s.Close()
}

type inmemEntry struct {
	next *inmemEntry
	data *Entry
}
