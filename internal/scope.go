package internal

import (
	"context"
	"encoding/json"
	"github.com/cbroglie/mustache"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

// Scope is the variable Scope of the sequence
type Scope struct {
	Scope   map[string]any
	Env     []string
	Results []StepResult
	Errors  map[string]string
}

const EnvValuesPath = "VALUES_PATH"
const EnvValuesPathDefault = "etc/values.yaml"
const VarParams = "params"
const VarValues = "values"

// NewScope is the constructor for Scope
func NewScope() Scope {
	scope := Scope{Scope: map[string]any{}, Env: make([]string, 0),
		Results: make([]StepResult, 0), Errors: make(map[string]string)}
	valuesPath := os.Getenv(EnvValuesPath)
	if len(valuesPath) == 0 {
		valuesPath = EnvValuesPathDefault
	}
	data, err := os.ReadFile(valuesPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load values file")
	}
	obj := make(map[string]any)
	if err = yaml.Unmarshal(data, &obj); err != nil {
		log.Fatal().Err(err).Msg("could not parse values file")
	}
	scope.Scope[VarValues] = obj
	return scope
}

// LoadEnvs loads the listed environment variables in the Scope
func (s *Scope) LoadEnvs(envs []string) {
	for _, env := range envs {
		s.Scope[env] = os.Getenv(env)
		s.Env = append(s.Env, env+"="+os.Getenv(env))
	}
}

// InsertParams inserts the params expressed as JSON into the Scope
func (s *Scope) InsertParams(data map[string]any) {
	s.Scope[VarParams] = data
}

// PushResult will store the result into the Scope and in the "Results" structure
func (s *Scope) PushResult(name string, value any, hide bool) {
	s.Scope[name] = value
	if !hide {
		s.Results = append(s.Results, StepResult{Name: name, Value: value})
	}
}

// Render renders a string template, against the Scope
func (s *Scope) Render(ctx context.Context, data string) (string, error) {
	templ, err := mustache.ParseString(data)
	if err != nil {
		return "", err
	}
	return templ.Render(s.Scope)
}

// RenderMap will recursively traverse a map and try to render all strings it finds
func (s *Scope) RenderMap(ctx context.Context, config map[string]any) (map[string]any, error) {
	var err error
	for k, v := range config {
		switch typed := v.(type) {
		case string:
			config[k], err = s.Render(ctx, typed)
		case map[string]any:
			config[k], err = s.RenderMap(ctx, typed)
		case []string:
			res := make([]string, 0)
			for _, item := range typed {
				rendered := ""
				rendered, err = s.Render(ctx, item)
				if err != nil {
					break
				}
				res = append(res, rendered)
			}
			config[k] = res
		}
		if err != nil {
			break
		}
	}
	return config, err
}

// ToJSON marshals the Scope to JSON
func (s *Scope) ToJSON() string {
	data, _ := json.Marshal(s.Scope)
	return string(data)
}

// ToProcessResult converts the internal representation of results, to a publishable version
func (s *Scope) ToProcessResult() ProcessResult {
	return ProcessResult{Results: s.Results, Errors: s.Errors}
}
