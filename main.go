package main

import (
	"context"
	"fmt"
	"local_google/html"
	"local_google/robot"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg = &sync.WaitGroup{}
	var ctx = context.Background()

	file, err := os.OpenFile("test.html", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	root, err := html.Parse(file)
	if err != nil {
		panic(err)
	}

	result, err := AnalyzePage(root)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", result)

	var robotCfg = robot.DefaultConfig()
	robotCfg.Queue = &robot.TaskQueue{}
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

	//var reader = bytes.NewReader(body)
	//req, err := http.NewRequest("GET", "https://zn.ua/ukr/UKRAINE/rosijska-ataka-na-kijiv-zahiblikh-vzhe-troje-postrazhdalij-vahitnij-zhintsi-zrobili-terminovu-operatsiju.html", reader)
	//if err != nil {
	//	panic(err)
	//}
	//
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//
	//body, err = io.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("\n%s\n", body)

	//var w = walker.New("1", wg, 10*time.Second)
	//w.Start(ctx)
	//
	//var s = searcher.New(wg)
	//if err := s.StartServer(); err != nil {
	//	panic(err)
	//}

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
