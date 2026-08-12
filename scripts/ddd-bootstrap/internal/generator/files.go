package generator

import (
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"

	"ddd-bootstrap/internal/config"
)

type fileSpec struct {
	template string
	target   string
}

type plannedFile struct {
	target  string
	content string
}

type planOptions struct {
	Root   string
	Force  bool
	DryRun bool
}

func (g *Generator) renderSpecs(specs []fileSpec, data config.TemplateData) ([]plannedFile, error) {
	files := make([]plannedFile, 0, len(specs))
	for _, spec := range specs {
		content, err := g.renderer.Render(spec.template, data)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(spec.target, ".go") {
			formatted, err := format.Source([]byte(content))
			if err != nil {
				return nil, fmt.Errorf("format generated Go file %s: %w", spec.target, err)
			}
			content = string(formatted)
		}
		files = append(files, plannedFile{target: spec.target, content: content})
	}
	return files, nil
}

func applyPlan(out io.Writer, opts config.InitOptions, dirs []string, files []plannedFile) error {
	return applyFilePlan(out, planOptions{
		Root:   opts.Output,
		Force:  opts.Force,
		DryRun: opts.DryRun,
	}, dirs, files)
}

func applyFilePlan(out io.Writer, opts planOptions, dirs []string, files []plannedFile) error {
	if opts.Root == "" {
		return fmt.Errorf("output root is required")
	}

	if opts.DryRun {
		for _, dir := range dirs {
			target := filepath.Join(opts.Root, filepath.FromSlash(dir))
			action := "CREATE-DIR"
			if exists, _ := pathExists(target); exists {
				action = "SKIP-DIR"
			}
			fmt.Fprintf(out, "%s %s\n", action, filepath.ToSlash(target))
		}
		for _, file := range files {
			action := "CREATE"
			target := filepath.Join(opts.Root, filepath.FromSlash(file.target))
			if exists, _ := pathExists(target); exists {
				if opts.Force {
					action = "OVERWRITE"
				} else {
					action = "SKIP"
				}
			}
			fmt.Fprintf(out, "%s %s\n", action, filepath.ToSlash(target))
		}
		return nil
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(opts.Root, filepath.FromSlash(dir)), 0o755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	for _, file := range files {
		target := filepath.Join(opts.Root, filepath.FromSlash(file.target))
		action := "CREATE"
		if exists, err := pathExists(target); err != nil {
			return err
		} else if exists {
			if !opts.Force {
				fmt.Fprintf(out, "SKIP %s\n", filepath.ToSlash(target))
				continue
			}
			action = "OVERWRITE"
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create parent directory for %s: %w", target, err)
		}
		if err := os.WriteFile(target, []byte(file.content), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", target, err)
		}
		fmt.Fprintf(out, "%s %s\n", action, filepath.ToSlash(target))
	}
	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
