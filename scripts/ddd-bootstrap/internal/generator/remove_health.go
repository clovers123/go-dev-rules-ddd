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

func (g *Generator) RemoveHealth(opts config.RemoveHealthOptions) error {
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		return err
	}

	for _, target := range healthRemovalTargets() {
		if err := removePath(g.out, opts.ProjectRoot, target, opts.DryRun); err != nil {
			return err
		}
	}
	return removeHealthRegistrations(g.out, opts.ProjectRoot, opts.DryRun)
}

func healthRemovalTargets() []string {
	return []string{
		"app/domain/health",
		"app/ohs/pl/health",
		"app/ohs/local/appservice/health",
		"app/ohs/remote/controller/health",
		"app/ohs/remote/routers/health.go",
		"app/acl/port/repository/health",
		"app/acl/pl/health",
		"app/acl/adapter/repository/postgres/health",
		"app/acl/adapter/repository/postgres/repository/export.health.health.go",
		"app/acl/adapter/repository/postgres/repository/health.health.gen.go",
		"app/acl/adapter/repository/postgres/model/health.health.gen.go",
		"app/acl/adapter/repository/postgres/repository/export.health.go",
		"app/acl/adapter/repository/postgres/repository/health.gen.go",
		"app/acl/adapter/repository/postgres/model/health.gen.go",
		"cmd/gorm-gen/health",
		"sql/health",
	}
}

func removePath(out io.Writer, root, rel string, dryRun bool) error {
	path := filepath.Join(root, filepath.FromSlash(rel))
	if exists, err := pathExists(path); err != nil {
		return err
	} else if !exists {
		fmt.Fprintf(out, "SKIP %s\n", filepath.ToSlash(path))
		return nil
	}
	if dryRun {
		fmt.Fprintf(out, "REMOVE %s\n", filepath.ToSlash(path))
		return nil
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove %s: %w", filepath.ToSlash(path), err)
	}
	fmt.Fprintf(out, "REMOVE %s\n", filepath.ToSlash(path))
	return nil
}

func removeHealthRegistrations(out io.Writer, root string, dryRun bool) error {
	updates := []struct {
		path string
		kind string
	}{
		{"cmd/server/di/acl.go", "lines"},
		{"cmd/server/di/domain.go", "lines"},
		{"cmd/server/di/ohs.go", "lines"},
		{"cmd/server/di/routers.go", "routers"},
		{"cmd/gorm-gen/main.go", "gormgen"},
	}
	for _, update := range updates {
		if err := removeHealthFromGoFile(out, filepath.Join(root, filepath.FromSlash(update.path)), update.kind, dryRun); err != nil {
			return err
		}
	}
	return nil
}

func removeHealthFromGoFile(out io.Writer, path, kind string, dryRun bool) error {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(out, "SKIP %s\n", filepath.ToSlash(path))
			return nil
		}
		return fmt.Errorf("read %s: %w", filepath.ToSlash(path), err)
	}
	content := string(contentBytes)
	updated := removeHealthLines(content)
	if kind == "gormgen" {
		updated = removeHealthGormGenBlock(updated)
	}
	if kind == "routers" && !strings.Contains(updated, "routers.Mount") {
		updated = removeLinesContaining(updated, []string{"/app/ohs/remote/routers"})
	}
	if updated == content {
		fmt.Fprintf(out, "SKIP %s\n", filepath.ToSlash(path))
		return nil
	}
	if dryRun {
		fmt.Fprintf(out, "UPDATE %s\n", filepath.ToSlash(path))
		return nil
	}
	formatted, err := format.Source([]byte(updated))
	if err != nil {
		return fmt.Errorf("format %s: %w", filepath.ToSlash(path), err)
	}
	if err := os.WriteFile(path, formatted, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filepath.ToSlash(path), err)
	}
	fmt.Fprintf(out, "UPDATE %s\n", filepath.ToSlash(path))
	return nil
}

func removeHealthLines(content string) string {
	return removeLinesContaining(content, []string{
		"/health\"",
		"healthservice.",
		"healthacl.",
		"healthrepo.",
		"healthadapter.",
		"healthappservice.",
		"healthcontroller.",
		"NewHealth",
		"MountHealthRouteGroup",
	})
}

func removeLinesContaining(content string, needles []string) string {
	var out []string
	for _, line := range strings.Split(content, "\n") {
		remove := false
		for _, needle := range needles {
			if strings.Contains(line, needle) {
				remove = true
				break
			}
		}
		if !remove {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

func removeHealthGormGenBlock(content string) string {
	lines := strings.Split(content, "\n")
	out := make([]string, 0, len(lines))
	skipping := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(line, "==================== Health") {
			skipping = true
			continue
		}
		if skipping {
			if trimmed == ")" {
				skipping = false
			}
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
