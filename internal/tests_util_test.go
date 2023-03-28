package internal

import "testing"

func TestPrototypeCheck(t *testing.T) {
	if err := PrototypeCheck(map[string]any{"aString": "bar", "anInt": 22, "anArray": []string{"foo"}},
		Proto{"aString": TYPE_STRING, "anInt": TYPE_INT, "anArray": TYPE_ARRAY}); err != nil {
		t.Error("prototype check returned a failure", err)
	}

	if err := PrototypeCheck(map[string]any{"aString": "bar", "anInt": 22, "anArray": []string{"foo"}},
		Proto{"aString": TYPE_INT, "anInt": TYPE_INT, "anArray": TYPE_ARRAY}); err == nil {
		t.Error("prototype check did not return a failure", err)
	}
}

func TestIsSequenceAllowed(t *testing.T) {
	if !IsSequenceAllowed("foo") {
		t.Error("correct sequence name was rejected")
	}
	if IsSequenceAllowed("foo.bar") {
		t.Error("invalid sequence not rejected")
	}
	if IsSequenceAllowed("foo/bar") {
		t.Error("invalid sequence not rejected")
	}
	if IsSequenceAllowed("foo\\bar") {
		t.Error("invalid sequence not rejected")
	}
	if IsSequenceAllowed("foo\"bar") {
		t.Error("invalid sequence not rejected")
	}
}
