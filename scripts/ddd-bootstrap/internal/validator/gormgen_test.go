package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckGormGenMainRejectsProjectConfigAPICoupling(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "cmd", "gorm-gen")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := `package main

func main() {
	config, _ := config.Load("configs/config.yml")
	_ = config.Postgres.DSN
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	var results []Result
	checkGormGenMain(root, func(level, rule, file, reason, fix string) {
		results = append(results, Result{Level: level, Rule: rule, File: file, Reason: reason, Fix: fix})
	})

	for _, result := range results {
		if result.Level == "FAIL" && result.Rule == "gormgen.config-loader" {
			return
		}
	}
	t.Fatalf("expected gormgen.config-loader failure, got %#v", results)
}

func TestCheckGormGenMainRejectsModelPkgPath(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "cmd", "gorm-gen")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := `package main

func main() {
	_ = gen.Config{ModelPkgPath: "./app/acl/adapter/repository/postgres/model"}
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	var results []Result
	checkGormGenMain(root, func(level, rule, file, reason, fix string) {
		results = append(results, Result{Level: level, Rule: rule, File: file, Reason: reason, Fix: fix})
	})

	for _, result := range results {
		if result.Level == "FAIL" && result.Rule == "gormgen.model-path" {
			return
		}
	}
	t.Fatalf("expected gormgen.model-path failure, got %#v", results)
}

func TestCheckGormGenOutputPathsRejectsNestedModelOutput(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "app", "acl", "adapter", "repository", "postgres", "app")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested output: %v", err)
	}

	var results []Result
	checkGormGenOutputPaths(root, func(level, rule, file, reason, fix string) {
		results = append(results, Result{Level: level, Rule: rule, File: file, Reason: reason, Fix: fix})
	})

	for _, result := range results {
		if result.Level == "FAIL" && result.Rule == "gormgen.model-path" {
			return
		}
	}
	t.Fatalf("expected gormgen.model-path failure, got %#v", results)
}

func TestCheckGormGenRejectsQuerierSuffix(t *testing.T) {
	root := t.TempDir()
	gormDir := filepath.Join(root, "cmd", "gorm-gen")
	if err := os.MkdirAll(filepath.Join(gormDir, "health"), 0o755); err != nil {
		t.Fatalf("mkdir gorm dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(gormDir, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}
	content := `package health

type HealthQuerier interface {
	// SELECT * FROM @@table WHERE owned_by=@ownedBy AND is_deleted=false
	FindPage(ownedBy string) error
}
`
	if err := os.WriteFile(filepath.Join(gormDir, "health", "health.health.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("write health interface: %v", err)
	}

	var results []Result
	checkGormGen(root, []generatedFlow{{
		Domain:   "health",
		Entity:   "health",
		Table:    "health",
		Schema:   "health",
		TypeName: "Health",
		GormFile: "cmd/gorm-gen/health/health.health.go",
	}}, func(level, rule, file, reason, fix string) {
		results = append(results, Result{Level: level, Rule: rule, File: file, Reason: reason, Fix: fix})
	})

	for _, result := range results {
		if result.Level == "FAIL" && result.Rule == "gormgen.interface-name" {
			return
		}
	}
	t.Fatalf("expected gormgen.interface-name failure, got %#v", results)
}
