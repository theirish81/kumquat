package internal

import (
	"testing"
)

func TestProcessResult_ToJSON(t *testing.T) {
	res := ProcessResult{Errors: map[string]string{"foo": "bar"}}
	str := string(res.ToJSON())
	if str != "{\"description\":\"\",\"results\":null,\"errors\":{\"foo\":\"bar\"}}" {
		t.Error("json conversion of ProcessResult did not work")
	}
}
