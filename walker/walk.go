package walker

import (
	"context"
	"log/slog"
	"time"
)

func (w *Walker) walk(ctx context.Context) {
	defer w.wg.Done()

	ticker := time.NewTicker(w.dur)

	for {
		if w.paused {
			slog.Debug("Walker %s is paused. Sleep %v", w.ID, w.dur)
			time.Sleep(w.dur)
		}

		select {
		case <-ticker.C:
			w.step()
		case <-ctx.Done():
			slog.Info("Walker %s is stopped", w.ID)
			return
		}
	}
}

func (w *Walker) step() {
	slog.Debug("Walker %s doing step", w.ID)

}
