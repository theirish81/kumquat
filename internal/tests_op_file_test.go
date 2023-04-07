package internal

import (
	"context"
	"os"
	"testing"
)

func TestFileOp_Run(t *testing.T) {
	_ = os.Setenv(EnvValuesPath, "../etc/values.yaml")
	scope := NewScope()
	op, _ := NewFileOp(map[string]any{"basePath": "../etc", "path": "test.md"}, &scope)
	_ = op.Run(context.Background())
	if op.GetResult() != "this is a test file" {
		t.Error("could not properly run file op")
	}
}
