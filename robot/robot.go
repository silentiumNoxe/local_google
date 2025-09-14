package robot

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

type Robot struct {
	ctx context.Context
	id  string

	ticker  *time.Ticker
	queue   *TaskQueue
	stopped bool
}

func New(ctx context.Context, cfg Config, id string) *Robot {
	return &Robot{
		ctx:    ctx,
		id:     id,
		ticker: time.NewTicker(cfg.Delay),
		queue:  cfg.Queue,
	}
}

// Start running endless loop. Execute this function as a goroutine.
func (r *Robot) Start() error {
	slog.Info("Robot %s started", r.id)
	for {
		if r.stopped {
			time.Sleep(time.Second * 5)
			continue
		}

		select {
		case <-r.ctx.Done():
			return nil
		case <-r.ticker.C:
			if err := r.step(); err != nil {
				return err
			}
		}
	}
}

func (r *Robot) step() error {
	slog.Info("Robot doing step", slog.String("rid", r.id))

	task := r.queue.Pop()
	if task == nil {
		return nil
	}

	url := task.Target
	if url == "" {
		return nil
	}

	reader, err := r.request(url)
	if err != nil {
		return err
	}

	node, err := html.Parse(reader)
	if err != nil {
		return err
	}

	return nil
}

func (r *Robot) request(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (r *Robot) Stop() {
	r.stopped = true
}

func (r *Robot) Stopped() bool {
	return r.stopped
}
