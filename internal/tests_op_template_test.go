package internal

import (
	"context"
	"testing"
)

func TestTemplateOp_Run(t *testing.T) {
	scope := NewScope()
	scope.Scope["foo"] = "bar"
	op, _ := NewTemplateOp(map[string]any{"template": "${foo}"}, &scope)
	_ = op.Run(context.Background())
	res := op.GetResult()
	if res.(string) != "bar" {
		t.Error("")
	}
}
