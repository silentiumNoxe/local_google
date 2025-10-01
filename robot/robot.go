package robot

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"local_google/html"
	"local_google/index"
	"local_google/robot/queue"
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
	queue   queue.Storage
	idx     index.Storage
	stopped bool
}

func New(ctx context.Context, cfg Config, id string) *Robot {
	return &Robot{
		ctx:    ctx,
		id:     id,
		ticker: time.NewTicker(cfg.Delay),
		queue:  cfg.Queue,
		idx:    cfg.Idx,
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

	entries, err := r.queue.Pop(1)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e == nil {
			continue
		}

		if e.Addr == "" {
			slog.Info("Empty target", slog.String("rid", r.id), slog.String("entry", hex.EncodeToString(e.ID)))
			continue
		}

		target, err := url.Parse(e.Addr)
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

		for _, l := range result.Links {
			err := r.queue.Put(&queue.Entry{Addr: l, Score: 0})
			if err != nil {
				slog.Warn("Unable to put entry to storage", slog.String("rid", r.id), slog.String("err", err.Error()))
			}
		}

		for word := range result.Index {
			id := index.IDGen([]byte(word))
			entry, err := r.idx.FindById(id)
			if err != nil {
				slog.Warn("Unable to find entry by id", slog.String("rid", r.id), slog.String("err", err.Error()))
				continue
			}

			if entry == nil {
				entry = &index.Entry{ID: id}
			}

			var exists = false
			for _, x := range entry.Addr {
				if x == e.Addr {
					exists = true
					break
				}
			}

			if exists {
				continue
			}

			entry.Addr = append(entry.Addr, e.Addr)
			err = r.idx.Save(entry)
			if err != nil {
				slog.Warn("Unable to save entry", slog.String("rid", r.id), slog.String("err", err.Error()))
			}
		}
	}

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
