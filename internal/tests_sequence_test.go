package internal

import (
	"context"
	"testing"
)

func TestLoadSequence(t *testing.T) {
	seq, _ := LoadSequence("../etc/sequences/http.yaml")
	if len(seq.Steps) == 0 {
		t.Error("sequences not loaded")
	}
}

func TestSequence_Run(t *testing.T) {
	seq, _ := LoadSequence("../etc/sequences/http.yaml")
	seq.Run(context.Background())
	result := seq.Result()
	if len(result.Results) == 0 {
		t.Error("sequence run produced no results")
	}

	seq, err := LoadSequence("bananas.yaml")
	if err == nil {
		t.Error("wrong file path should return an error")
	}

	seq, err = LoadSequence("../README.md")
	if err == nil {
		t.Error("loading a broken YAML syntax should return an error")
	}
}
