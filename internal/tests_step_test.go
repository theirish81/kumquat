package internal

import (
	"reflect"
	"testing"
)

func TestStep_GetImplementation(t *testing.T) {
	step := Step{Type: OP_NIX_SHELL, Config: map[string]any{}}
	impl, _ := step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.NixShellOp" {
		t.Error("nix shell op not implemented correctly")
	}
	step = Step{Type: OP_SQL, Config: map[string]any{}}
	impl, _ = step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.SqlOp" {
		t.Error("sql op not implemented correctly")
	}
	step = Step{Type: OP_MONGO, Config: map[string]any{}}
	impl, _ = step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.MongoOp" {
		t.Error("mongo op not implemented correctly")
	}
	step = Step{Type: OP_TEMPLATE, Config: map[string]any{}}
	impl, _ = step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.TemplateOp" {
		t.Error("template op not implemented correctly")
	}
	step = Step{Type: OP_WAIT, Config: map[string]any{}}
	impl, _ = step.GetImplementation()
	if reflect.TypeOf(impl).String() != "*internal.WaitOp" {
		t.Error("wait op not implemented correctly")
	}
}
