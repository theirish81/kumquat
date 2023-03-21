package internal

import (
	"context"
	"testing"
)

func TestSplitFilter_Run(t *testing.T) {
	filter, _ := NewSplitFilter(map[string]any{"sep": "|"})
	res, _ := filter.Run(context.Background(), "foo|bar")
	resCasted := res.([]string)
	if resCasted[0] != "foo" || resCasted[1] != "bar" {
		t.Error("split did not work as expected for string")
	}

	res, _ = filter.Run(context.Background(), []string{"foo|bar", "foobar"})
	resCasted2 := res.([][]string)
	if resCasted2[0][0] != "foo" || resCasted2[0][1] != "bar" || resCasted2[1][0] != "foobar" {
		t.Error("split did not work as expected for string array")
	}
}

func TestSplitLinesFilter_Run(t *testing.T) {
	filter := NewSplitLinesFilter()
	res, _ := filter.Run(context.Background(), "foo\nbar")
	if res.([]string)[0] != "foo" || res.([]string)[1] != "bar" {
		t.Error("split lines did not work")
	}
}

func TestJsonParseFilter_Run(t *testing.T) {
	filter := NewJsonParseFilter()
	data, _ := filter.Run(context.Background(), `{"foo":"bar"}`)
	out, _ := data.(map[string]any)["foo"]
	if out != "bar" {
		t.Error("json parse filter did not work")
	}
}

func TestRegexpReplaceFilter_Run(t *testing.T) {
	filter, _ := NewRegexpReplaceFilter(map[string]any{"regexp": "\\s\\s+", "replace": " "})
	res, _ := filter.Run(context.Background(), "foo  bar")
	if res != "foo bar" {
		t.Error("regexp filter did not work")
	}
}
