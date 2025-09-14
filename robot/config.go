package robot

import "time"

type Config struct {
	Delay time.Duration
	Queue *TaskQueue
}

func DefaultConfig() Config {
	return Config{
		Delay: 20 * time.Second,
	}
}
