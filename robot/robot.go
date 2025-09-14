package robot

import (
	"context"
	"fmt"
	"io"
	"local_google/html"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	r.stopped = false
	slog.Info("Robot started", slog.String("rid", r.id))
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
	slog.Debug("Robot step", slog.String("rid", r.id))

	task := r.queue.Pop()
	if task == nil {
		slog.Info("Queue is empty", slog.String("rid", r.id))
		return nil
	}

	if task.Target == "" {
		return fmt.Errorf("empty target")
	}

	target, err := url.Parse(task.Target)
	if err != nil {
		return err
	}

	reader, err := r.request(target)
	if err != nil {
		return err
	}

	slog.Debug("Parse page", slog.String("rid", r.id))
	node, err := html.Parse(reader)
	if err != nil {
		return err
	}

	slog.Debug("Analyze page", slog.String("rid", r.id))
	result, err := analyzePage(node, fmt.Sprintf("%s://%s", target.Scheme, target.Host))
	if err != nil {
		return err
	}

	r.queue.Mutex.Lock()
	for _, l := range result.Links {
		r.queue.Push(l)
	}
	r.queue.Mutex.Unlock()

	return nil
}

func (r *Robot) request(u *url.URL) (io.Reader, error) {
	slog.Debug("Request page", slog.String("rid", r.id), slog.String("url", u.String()))
	resp, err := http.Get(u.String())
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

func analyzePage(root *html.Node, domain string) (*AnalyzeResult, error) {
	var result = AnalyzeResult{
		Index: make(map[string]int),
	}

	var extractContent = func(s string) {
		s = strings.ToLower(s)
		entries := strings.Split(s, " ")
		for _, x := range entries {
			x = strings.TrimSpace(x)
			if x == "" {
				continue
			}
			result.Index[x]++
		}
	}

	var walker = html.NewWalker(root)

	for {
		n := walker.Next()
		if n == nil {
			break
		}

		if n.Tag == "a" {
			uri := n.Attr["href"]
			if uri != "" {
				if !strings.HasPrefix(uri, "http") {
					uri = domain + uri
				}

				result.Links = append(result.Links, uri)
			}
		}

		extractContent(n.Content)
	}

	return &result, nil
}

type AnalyzeResult struct {
	Index map[string]int
	Links []string
}
