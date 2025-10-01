package robot

import (
	"local_google/index"
	"local_google/robot/queue"
	"time"
)

type Config struct {
	Delay time.Duration
	Queue queue.Storage
	Idx   index.Storage
}

func DefaultConfig() Config {
	return Config{
		Delay: 20 * time.Second,
	}
}
