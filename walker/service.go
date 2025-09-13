package walker

import (
	"context"
	"sync"
	"time"
)

type Walker struct {
	ID     string
	paused bool
	inited bool

	wg  *sync.WaitGroup
	dur time.Duration
}

func New(ID string, wg *sync.WaitGroup, dur time.Duration) *Walker {
	return &Walker{
		ID:  ID,
		wg:  wg,
		dur: dur,
	}
}

func (w *Walker) Start(ctx context.Context) {
	w.paused = false
	if w.inited {
		return
	}

	w.inited = true
	w.wg.Add(1)
	go w.walk(ctx)
}

func (w *Walker) Stop() {
	w.paused = true
}
