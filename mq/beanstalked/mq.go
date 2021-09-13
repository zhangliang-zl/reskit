package beanstalked

import "time"

const (
	MinPriority    uint32 = 1024
	MaxPriority    uint32 = 2048
	MaxWorkingTTL         = time.Second * 120
	FailRetryDelay        = time.Second * 10
)