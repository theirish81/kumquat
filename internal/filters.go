package internal

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

// IFilter is the prototype of all filters
type IFilter interface {
	Run(ctx context.Context, input any) (any, error)
}

const FilterSplitLines = "splitLines"
const FilterSplit = "split"
const FilterJsonParse = "jsonParse"
const FilterReplace = "replace"

// Filter is the wrapping structure of all filters
type Filter struct {
	Type   string         `yaml:"type"`
	Config map[string]any `yaml:"config"`
}

// GetImplementation will return the IFilter proper implementation
func (f *Filter) GetImplementation() (IFilter, error) {
	switch f.Type {
	case FilterSplitLines:
		return NewSplitLinesFilter(), nil
	case FilterSplit:
		return NewSplitFilter(f.Config)
	case FilterJsonParse:
		return NewJsonParseFilter(), nil
	case FilterReplace:
		return NewRegexpReplaceFilter(f.Config)
	default:
		return nil, errors.New("could not find a proper implementation for " + f.Type)
	}
}

// Run runs the filter
func (f *Filter) Run(ctx context.Context, input any) (any, error) {
	op, err := f.GetImplementation()
	if err != nil {
		return nil, err
	}
	return op.Run(ctx, input)
}

// SplitLinesFilter takes the input as string and splits the lines, returning an array
type SplitLinesFilter struct{}

// NewSplitLinesFilter is the constructor for SplitLinesFilter
func NewSplitLinesFilter() *SplitLinesFilter {
	return &SplitLinesFilter{}
}

// Run will run the SplitLinesFilter
func (f *SplitLinesFilter) Run(_ context.Context, input any) (any, error) {
	return strings.Split(input.(string), "\n"), nil
}

// SplitFilter splits a string into an array based on a separation character.
// If the input is an array, the output will be an array of arrays
type SplitFilter struct {
	Sep string
}

const splitFilterSchema = `{
	"required": [
		"sep"
	],
	"properties": {
		"sep": {
			"type": "string"
		}
	}
}`

// NewSplitFilter is the constructor for SplitFilter
func NewSplitFilter(config map[string]any) (*SplitFilter, error) {
	if err := PrototypeCheck(config, splitFilterSchema); err == nil {
		return &SplitFilter{Sep: config["sep"].(string)}, nil
	} else {
		return nil, err
	}

}

// Run will run the filter, transforming the input and returning the transformed data
func (f *SplitFilter) Run(_ context.Context, input any) (any, error) {
	switch input.(type) {
	case string:
		// if the type of the input is a string, then we slit it
		return strings.Split(input.(string), f.Sep), nil
	case []string:
		// if the type of the input is an array of strings, we slit each string
		res := make([][]string, 0)
		for _, str := range input.([]string) {
			res = append(res, strings.Split(str, f.Sep))
		}
		return res, nil
	}
	// in case input type matches nothing, then we return the input as is
	return input, nil
}

// JsonParseFilter is a filter that transforms text to a JSON structure
type JsonParseFilter struct{}

// NewJsonParseFilter is the constructor for JsonParseFilter
func NewJsonParseFilter() *JsonParseFilter {
	return &JsonParseFilter{}
}

// Run will run the filter, transforming the input and returning the transformed data
func (f *JsonParseFilter) Run(_ context.Context, input any) (any, error) {
	var data any
	err := json.Unmarshal([]byte(input.(string)), &data)
	return data, err
}

// RegexpReplaceFilter is a filter that will take an input and replace any matching pattern described by
// Regexp with the value of Replace
type RegexpReplaceFilter struct {
	Regexp  *regexp.Regexp
	Replace string
}

const regexpReplaceFilterSchema = `{
	"required": [
		"regexp",
		"replace"
	],
	"properties": {
		"regexp": {
			"type": "string"
		},
		"replace": {
			"type": "string"
		}
	}
}`

// NewRegexpReplaceFilter is the constructor for RegexpReplaceFilter
func NewRegexpReplaceFilter(config map[string]any) (*RegexpReplaceFilter, error) {
	if err := PrototypeCheck(config, regexpReplaceFilterSchema); err == nil {
		rx, err := regexp.Compile(config["regexp"].(string))
		return &RegexpReplaceFilter{Regexp: rx, Replace: config["replace"].(string)}, err
	} else {
		return nil, err
	}
}

// Run will run the filter, transforming the input and returning the transformed data
func (f *RegexpReplaceFilter) Run(_ context.Context, input any) (any, error) {
	return string(f.Regexp.ReplaceAll([]byte(input.(string)), []byte(f.Replace))), nil
}
