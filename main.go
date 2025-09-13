package main

import (
	"fmt"
	"io"
	"local_google/html"
	"os"
	"sync"
)

func main() {
	var wg = &sync.WaitGroup{}
	//var ctx = context.Background()

	file, err := os.OpenFile("sample_3.html", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	root, err := html.Parse(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", root)

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

type SimpleWriter struct {
	arr []byte
	i   int
}

func NewSimpleWriter(arr []byte) *SimpleWriter {
	return &SimpleWriter{arr: arr}
}

func (w *SimpleWriter) WriteByte(b byte) error {
	if w.i >= len(w.arr) {
		return io.EOF
	}

	w.arr[w.i] = b
	w.i++

	return nil
}

func (w *SimpleWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if err := w.WriteByte(b); err != nil {
			return n, err
		}
	}

	return len(p), nil
}

func (w *SimpleWriter) Reset() {
	w.i = 0
}

type Node struct {
	Tag   string
	Attr  map[string]string
	Child []*Node
}

func NewNode(tag string) *Node {
	return &Node{
		Tag:   tag,
		Attr:  map[string]string{},
		Child: []*Node{},
	}
}

func clearBuf(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
