package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	texttemplate "text/template"
)

type Renderer struct {
	root string
}

func NewRenderer(root string) (*Renderer, error) {
	if root == "" {
		found, err := findTemplateRoot()
		if err != nil {
			return nil, err
		}
		root = found
	}
	return &Renderer{root: root}, nil
}

func (r *Renderer) Render(relativePath string, data any) (string, error) {
	path := filepath.Join(r.root, filepath.FromSlash(relativePath))
	tpl, err := texttemplate.New(filepath.Base(relativePath)).
		Option("missingkey=error").
		ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("parse template %s: %w", relativePath, err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render template %s: %w", relativePath, err)
	}
	return buf.String(), nil
}

func findTemplateRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(wd, "templates")
		marker := filepath.Join(candidate, "project", "go.mod.tpl")
		if _, err := os.Stat(marker); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return "", fmt.Errorf("template root not found; run ddd-bootstrap from the skill root or provide templates directory")
}
