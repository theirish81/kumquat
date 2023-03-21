package internal

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestNixShellOp_Run(t *testing.T) {
	_ = os.Setenv(EnvValuesPath, "../etc/values.yaml")
	scope := NewScope()
	op, _ := NewNixShellIOp(map[string]any{"command": "ls ../"}, &scope)
	_ = op.Run(context.Background())
	res := op.GetResult()
	if !strings.Contains(res.(string), "main.go") {
		t.Error("error while running NixShellOp")
	}
}
