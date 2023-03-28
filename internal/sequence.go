package internal

import (
	"context"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Sequence is a sequence of operations
type Sequence struct {
	Env          []string `yaml:"env"`
	Description  string   `yaml:"description"`
	Steps        []*Step  `yaml:"steps"`
	AcceptParams bool     `yaml:"accept_params"`
	Requires     []string `yaml:"requires"`
	Scope        *Scope
}

// LoadSequence loads a sequence from a file path
func LoadSequence(filePath string) (Sequence, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Sequence{}, err
	}
	scope := NewScope()
	seq := Sequence{Scope: &scope}
	err = yaml.Unmarshal(data, &seq)
	if err != nil {
		return seq, err
	}
	scope.LoadEnvs(seq.Env)
	for _, step := range seq.Steps {
		step.scope = &scope
	}
	return seq, nil
}

// Run will run the sequence
func (s *Sequence) Run(ctx context.Context) {
	for _, step := range s.Steps {
		if val, err := step.Run(ctx); err == nil {
			s.Scope.PushResult(step.Name, val, step.Hide)
		} else {
			log.Error().Err(err).Str("step", step.Name).Msg("could not run step")
			s.Scope.Errors[step.Name] = err.Error()
		}

	}
}

// CheckRequires will check that the provided map has at least the fields described in the "Requires" list
func (s *Sequence) CheckRequires(data map[string]any) error {
	if s.Requires != nil {
		for _, r := range s.Requires {
			if _, ok := data[r]; !ok {
				return errors.New("required parameter missing: " + r)
			}
		}
	}
	return nil
}

// Result will return publishable results for the sequence
func (s *Sequence) Result() ProcessResult {
	res := s.Scope.ToProcessResult()
	res.Description = s.Description
	return res
}
