package msgstream

import (
	"context"
	"sync"
)

type BaseMsg struct {
	mu             sync.Mutex
	Ctx            context.Context
	BeginTimestamp Timestamp
	EndTimestamp   Timestamp
	HashValues     []uint32
}
