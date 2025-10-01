package main

import (
	"context"
	"fmt"
	"local_google/robot"
	"local_google/robot/queue"
	"log/slog"
	"sync"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	var wg = &sync.WaitGroup{}
	var ctx = context.Background()

	var storageCfg = &queue.Config{}
	var storage queue.Storage
	storage = queue.NewLocalStorage(storageCfg)
	storage = queue.NewInMemStorage(storageCfg, storage)

	var robotCfg = robot.DefaultConfig()
	robotCfg.Queue = storage
	err := robotCfg.Queue.Put(&queue.Entry{Addr: "https://zn.ua/ukr/UKRAINE/rosijska-ataka-na-kijiv-zahiblikh-vzhe-troje-postrazhdalij-vahitnij-zhintsi-zrobili-terminovu-operatsiju.html"})
	if err != nil {
		panic(err)
	}

	var robots = make([]*robot.Robot, 3)
	for i := 0; i < len(robots); i++ {
		r := robot.New(ctx, robotCfg, fmt.Sprintf("r-%d", i))
		robots[i] = r

		wg.Add(1)
		go func(r *robot.Robot) {
			defer wg.Done()

			if err := r.Start(); err != nil {
				panic(err)
			}
		}(r)
	}

	wg.Wait()
	if err := storage.Close(); err != nil {
		slog.Warn("Failed close storage", slog.String("err", err.Error()))
	}
}
