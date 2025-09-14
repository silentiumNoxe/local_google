package robot

import (
	"fmt"
	"log/slog"
	"sync"
)

type TaskQueue struct {
	First *Task
	Last  *Task
	lock  *sync.Mutex
}

func (q *TaskQueue) Push(target string) {
	q.lock.Lock()
	defer q.lock.Unlock()

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
	q.lock.Lock()
	defer q.lock.Unlock()

	var x = q.First
	if x == nil {
		return nil
	}

	q.First = x.Next
	return x
}

type Task struct {
	Next   *Task
	Target string
}
