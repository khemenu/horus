package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/google/uuid"
)

type Reporter interface {
	Report(raw any) (string, error)
}

func NewReporter(format string, opt string) (Reporter, error) {
	switch format {
	case "plain":
		return nil, nil

	case "template":
		if opt == "" {
			return nil, fmt.Errorf("no template expression is provided")
		}

		tmpl := template.New("reporter")
		tmpl.Funcs(template.FuncMap{
			"uuid": func(i interface{}) (uuid.UUID, error) {
				switch v := i.(type) {
				case string:
					return uuid.Parse(v)
				case []byte:
					return uuid.FromBytes(v)
				default:
					return uuid.Nil, fmt.Errorf("invalid input type: %T", i)
				}
			},
		})

		tmpl, err := tmpl.Parse(opt)
		if err != nil {
			return nil, fmt.Errorf("parse template: %w", err)
		}

		// tmpl.Funcs()

		return &TemplateReporter{
			Template: tmpl,
		}, nil

	case "json":
		return &JsonReporter{}, nil

	default:
		return nil, fmt.Errorf("unknown formatter")
	}
}

type TemplateReporter struct {
	*template.Template
}

func (r *TemplateReporter) Report(raw any) (string, error) {
	var b bytes.Buffer
	if err := r.Execute(&b, raw); err != nil {
		return "", err
	}

	return b.String(), nil
}

type JsonReporter struct{}

func (r *JsonReporter) Report(raw any) (string, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(raw); err != nil {
		return "", err
	}

	return b.String(), nil
}
