package internal

import (
	"github.com/bitly/go-simplejson"
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
	dx, _ := simplejson.NewJson([]byte("{\"foo\":\"bar\"}"))
	data, _ := dx.Map()
	scope.InsertParams(data)
	if k, _ := scope.Scope["Params"].(map[string]any)["foo"]; k != "bar" {
		t.Error("params not set correctly")
	}
}
