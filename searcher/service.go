package searcher

import (
	"errors"
	"log/slog"
	"net/http"
	"sync"
)

type Searcher struct {
	wg *sync.WaitGroup
}

func New(wg *sync.WaitGroup) *Searcher {
	return &Searcher{wg: wg}
}

func (s *Searcher) Search(query string) []string {
	return []string{}
}

func (s *Searcher) StartServer() error {
	mux := http.NewServeMux()

	s.wg.Add(1)
	go func(addr string, mux *http.ServeMux) {
		defer s.wg.Done()

		slog.Info("Starting server on %s", addr)
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
			slog.Info("Server closed")
		}
	}(":8080", mux)
	return nil
}
