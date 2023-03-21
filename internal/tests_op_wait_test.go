package internal

import (
	"context"
	"testing"
	"time"
)

func TestWaitOp_Run(t *testing.T) {
	op, _ := NewWaitOp(map[string]any{"duration": "10ms"})
	start := time.Now()
	_ = op.Run(context.Background())
	end := time.Now()
	if end.Sub(start).Milliseconds() < 10 {
		t.Error("wait did not... wait")
	}
}
