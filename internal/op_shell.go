package internal

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// NixShellOp is an operation that runs a *NIX command
type NixShellOp struct {
	Command string
	Stdin   bool
	Result  string
	Timeout time.Duration
	scope   *Scope
}

const nixShellOpSchema = `{
	"required": [
		"command",
		"stdin",
		"timeout"
	],
	"properties": {
		"command": {
			"type":"string"
		},
		"stdin": {
			"type":"boolean"
		},
		"timeout": {
			"type":"string"
		}
	}
}`

// NewNixShellIOp constructor for NixShellOp
func NewNixShellIOp(config map[string]any, scope *Scope) (*NixShellOp, error) {
	config = SetDefault(config, "stdin", false)
	config = SetDefault(config, "timeout", "10s")
	duration, err := time.ParseDuration(config["timeout"].(string))
	if err != nil {
		return nil, err
	}
	if err := PrototypeCheck(config, nixShellOpSchema); err == nil {
		return &NixShellOp{Command: config["command"].(string),
			Stdin:   config["stdin"].(bool),
			Timeout: duration,
			scope:   scope}, nil
	} else {
		return nil, err
	}
}

// Run runs the operation
func (o *NixShellOp) Run(ctx context.Context) error {
	log.Debug().Str("operation", "NixShellOp").Str("command", o.Command).Msg("running NIX op")
	evaluatedCommand, err := o.scope.Render(ctx, o.Command)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, o.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", evaluatedCommand)
	cmd.Env = o.scope.Env
	// if the Stdin switch is on, then we pipe the Scope into the process standard input
	if o.Stdin {
		_ = o.pipeScope(cmd)
	}
	data, err := cmd.Output()
	if err != nil {
		return err
	}
	if ctx.Err() == context.DeadlineExceeded {
		return ctx.Err()
	}
	o.Result = strings.TrimSpace(string(data))
	return nil
}

// pipeScope pipes the scope into the child process standard input
func (o *NixShellOp) pipeScope(cmd *exec.Cmd) error {
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	_, _ = pipe.Write([]byte(o.scope.ToJSON()))
	return pipe.Close()
}

// GetResult returns the result of the operation
func (o *NixShellOp) GetResult() any {
	return o.Result
}
