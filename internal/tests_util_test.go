package internal

import "testing"

func TestPrototypeCheck(t *testing.T) {
	if err := PrototypeCheck(map[string]any{"aString": "bar", "anInt": 22, "anArray": []any{"foo"}}, `{
			"properties": {
				"aString": {
					"type": "string"
				},
				"anInt": {
					"type": "integer"
				},
				"anArray": {
					"type": "array",
					"items": {
						"type":"string"
					}
				}
			}
		}`); err != nil {
		t.Error("prototype check returned a failure", err)
	}

	if err := PrototypeCheck(map[string]any{"aString": "bar", "anInt": 22, "anArray": []any{"foo"}},
		`{
			"properties": {
				"aString": {
					"type": "int"
				},
				"anInt": {
					"type": "int"
				},
				"anArray": {
					"type": "array",
					"items": "string"
				}
			}
		}`); err == nil {
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
