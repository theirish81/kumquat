package internal

import "context"

// TemplateOp is an operation that evaluates a template
type TemplateOp struct {
	Template string
	Result   string
	scope    *Scope
}

const templateOpSchema = `{
	"required": [
		"template"
	],
	"properties": {
		"template": {
			"type": "string"
		}
	}
}
`

// NewTemplateOp is the constructor for TemplateOp
func NewTemplateOp(config map[string]any, scope *Scope) (*TemplateOp, error) {
	if err := PrototypeCheck(config, templateOpSchema); err == nil {
		return &TemplateOp{Template: config["template"].(string), scope: scope}, nil
	} else {
		return nil, err
	}
}

// Run runs the template
func (o *TemplateOp) Run(ctx context.Context) error {
	var err error
	o.Result, err = o.scope.Render(ctx, o.Template)
	return err
}

// GetResult returns the result
func (o *TemplateOp) GetResult() any {
	return o.Result
}
