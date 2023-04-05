package internal

import (
	"reflect"
	"testing"
)

func TestStep_GetImplementation(t *testing.T) {
	step := Step{Type: OpNixShell, Config: map[string]any{}}
	impl, _ := step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.NixShellOp" {
		t.Error("nix shell op not implemented correctly")
	}
	step = Step{Type: OpSql, Config: map[string]any{}}
	impl, _ = step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.SqlOp" {
		t.Error("sql op not implemented correctly")
	}
	step = Step{Type: OpMongo, Config: map[string]any{}}
	impl, _ = step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.MongoOp" {
		t.Error("mongo op not implemented correctly")
	}
	step = Step{Type: OpTemplate, Config: map[string]any{}}
	impl, _ = step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.TemplateOp" {
		t.Error("template op not implemented correctly")
	}
	step = Step{Type: OpWait, Config: map[string]any{}}
	impl, _ = step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.WaitOp" {
		t.Error("wait op not implemented correctly")
	}
}
