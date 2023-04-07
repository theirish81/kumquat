package internal

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
)

type FileOp struct {
	Path     string
	BasePath string
	Raw      bool
	scope    *Scope
	Result   any
}

const fileSchema = `{
	"required": [
		"path"
	],
	"properties": {
		"path": {
			"type": "string"
		},
		"raw": {
			"type": "boolean"
		}
	}
}`

func NewFileOp(config map[string]any, scope *Scope) (*FileOp, error) {
	config = SetDefault(config, "basePath", "")
	config = SetDefault(config, "raw", false)
	if err := PrototypeCheck(config, fileSchema); err == nil {
		return &FileOp{Path: config["path"].(string), BasePath: config["basePath"].(string), Raw: config["raw"].(bool),
			scope: scope}, nil
	} else {
		return nil, err
	}
}

func (o *FileOp) Run(ctx context.Context) error {
	log.Debug().Str("operation", "FileOp").Str("path", o.Path).Msg("running File op")
	actualPath, err := o.scope.Render(ctx, o.Path)
	if err != nil {
		return err
	}
	if strings.HasPrefix(actualPath, "/") || strings.Contains(actualPath, "..") {
		return errors.New("invalid path")
	}
	data, err := os.ReadFile(path.Join(o.BasePath, actualPath))
	if err != nil {
		return err
	}
	if o.Raw {
		o.Result = data
	} else {
		o.Result = string(data)
	}

	return nil
}

func (o *FileOp) GetResult() any {
	return o.Result
}
