package internal

import (
	"context"
	"time"
)

// WaitOp is an operation that pauses the sequence for a certain amount of time
type WaitOp struct {
	duration time.Duration
}

// NewWaitOp is the constructor for WaitOp
func NewWaitOp(config map[string]any) (*WaitOp, error) {
	if err := PrototypeCheck(config, Proto{"duration": TYPE_STRING}); err == nil {
		d, err := time.ParseDuration(config["duration"].(string))
		return &WaitOp{duration: d}, err
	} else {
		return nil, err
	}

}

// Run will run the wait command
func (o *WaitOp) Run(_ context.Context) error {
	time.Sleep(o.duration)
	return nil
}

// GetResult will return the result, that's nothing more than an informative string
func (o *WaitOp) GetResult() any {
	return "waited for " + o.duration.String()
}
