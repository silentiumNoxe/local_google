package main

import (
	"context"
	"fmt"
	"local_google/html"
	"local_google/robot"
	"log/slog"
	"strings"
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

func AnalyzePage(root *html.Node) (*AnalyzeResult, error) {
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
