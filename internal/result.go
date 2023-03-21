package internal

import "encoding/json"

// ProcessResult is the final result of a sequence execution
type ProcessResult struct {
	Description string            `json:"description"`
	Results     []StepResult      `json:"results"`
	Errors      map[string]string `json:"errors"`
}

// ToJSON will turn the Result data structure into a JSON string
func (r *ProcessResult) ToJSON() []byte {
	data, _ := json.Marshal(r)
	return data
}

// StepResult is the outcome of one step
type StepResult struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}
