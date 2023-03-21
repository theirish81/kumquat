package internal

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
)

// Step is one step in a sequence
type Step struct {
	Name    string         `yaml:"name"`
	Type    string         `yaml:"type"`
	Config  map[string]any `yaml:"config"`
	scope   *Scope
	Filters []Filter `yaml:"filters"`
	Hide    bool     `yaml:"hide"`
}

// GetImplementation will return the proper IOperation implementation for the current step
func (s *Step) GetImplementation() (IOperation, error) {
	switch s.Type {
	case OP_NIX_SHELL:
		return NewNixShellIOp(s.Config, s.scope)
	case OP_TEMPLATE:
		return NewTemplateOp(s.Config, s.scope)
	case OP_WAIT:
		return NewWaitOp(s.Config)
	case OP_SQL:
		return NewSqlOp(s.Config, s.scope)
	case OP_MONGO:
		return NewMongoOp(s.Config, s.scope)
	default:
		return nil, errors.New("could not find a proper implementation for " + s.Type)
	}
}

// Run will run the step
func (s *Step) Run(ctx context.Context) (any, error) {
	op, err := s.GetImplementation()
	if err != nil {
		return nil, err
	}
	err = op.Run(ctx)
	if err != nil {
		return nil, err
	}
	res := op.GetResult()
	for _, f := range s.Filters {
		rx, err := f.Run(ctx, res)
		if err == nil {
			res = rx
		} else {
			log.Err(err).Str("filter", f.Type).Msg("could not run filter")
		}
	}
	return res, nil
}
