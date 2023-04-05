package internal

import (
	"os"
	"testing"
)

func TestScope_LoadEnvs(t *testing.T) {
	scope := NewScope()
	_ = os.Setenv("FOO", "BAR")
	scope.LoadEnvs([]string{"FOO"})
	if scope.Scope["FOO"] != "BAR" {
		t.Error("env variable not loaded correctly")
	}
	if len(scope.Scope) != 2 {
		t.Error("unexpected variables in Scope")
	}
}

func TestScope_InsertParams(t *testing.T) {
	scope := NewScope()
	data := map[string]any{"foo": "bar"}
	scope.InsertParams(data)
	if k, _ := scope.Scope["Params"].(map[string]any)["foo"]; k != "bar" {
		t.Error("params not set correctly")
	}
}
