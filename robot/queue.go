package robot

import (
	"fmt"
	"log/slog"
	"sync"
)

type TaskQueue struct {
	First *Task
	Last  *Task
	Mutex *sync.Mutex
}

func (q *TaskQueue) Push(target string) {
	slog.Debug(fmt.Sprintf("Added to queue target: %s", target))

	task := &Task{Target: target}

	if q.First == nil {
		q.First = task
		q.Last = q.First
		return
	}

	last := q.Last
	q.Last = task
	last.Next = task
}

func (q *TaskQueue) Pop() *Task {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	var x = q.First
	if x == nil {
		return nil
	}

	q.First = x.Next
	return x
}

func (q *TaskQueue) Empty() bool {
	return q.First == nil
}

type Task struct {
	Next   *Task
	Target string
}
