package main

import (
	"context"
	"fmt"
	"local_google/robot"
	"log/slog"
	"sync"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	var wg = &sync.WaitGroup{}
	var ctx = context.Background()

	var robotCfg = robot.DefaultConfig()
	robotCfg.Queue = &robot.TaskQueue{Mutex: &sync.Mutex{}}
	robotCfg.Queue.Push("https://zn.ua/ukr/UKRAINE/rosijska-ataka-na-kijiv-zahiblikh-vzhe-troje-postrazhdalij-vahitnij-zhintsi-zrobili-terminovu-operatsiju.html")

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
}
