package generator

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ddd-bootstrap/internal/config"
)

func TestInitGeneratesHealthOnlyAndSelfContainedGormGen(t *testing.T) {
	root := t.TempDir()
	gen, err := New(&bytes.Buffer{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	opts := config.DefaultInitOptions()
	opts.ProjectName = "rem-sso-backend-ddd"
	opts.ModulePath = "rem-sso-backend-ddd"
	opts.Output = root
	opts.DBPassword = "EVS3%Jthp$Cb7WDkGVzhDgkDhhrNjmfD"
	opts.SkipGoModTidy = true

	if err := gen.Init(opts); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	for _, rel := range []string{
		"app/domain/health/entity/health.go",
		"app/ohs/local/appservice/health/health.go",
		"app/ohs/remote/controller/health/health.go",
		"app/ohs/remote/routers/health.go",
		"cmd/gorm-gen/health/health.health.go",
		"sql/health/health.health.sql",
		"app/acl/adapter/repository/postgres/repository/export.health.health.go",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}

	for _, rel := range []string{
		"app/domain/user",
		"app/application",
		"app/acl/pl/user",
		"cmd/gorm-gen/user",
		"sql/user",
		"cmd/gorm-gen/health/health.go",
		"sql/health/health.sql",
		"app/acl/adapter/repository/postgres/repository/export.health.go",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); !os.IsNotExist(err) {
			t.Fatalf("init generated an unexpected legacy or inferred path: %s", rel)
		}
	}

	mainBytes, err := os.ReadFile(filepath.Join(root, "cmd", "gorm-gen", "main.go"))
	if err != nil {
		t.Fatalf("read gorm-gen main: %v", err)
	}
	main := string(mainBytes)
	for _, forbidden := range []string{"configs.Load(", "config.Postgres.DSN", "/configs\"", "ModelPkgPath", "HealthQuerier"} {
		if strings.Contains(main, forbidden) {
			t.Fatalf("gorm-gen main must not contain %q", forbidden)
		}
	}
	for _, expected := range []string{
		"loadPostgresDSN(\"configs/config.yml\", \"default\")",
		"gorm.Open(postgres.Open(dsn)",
		"gopkg.in/yaml.v3",
		"datatypes.JSON",
		"gen.FieldType(\"create_time\", \"*time.Time\")",
		"gen.FieldType(\"update_time\", \"*time.Time\")",
		"WithUnitTest:      false",
		"g.ApplyInterface(func(health.Health) {}",
	} {
		if !strings.Contains(main, expected) {
			t.Fatalf("gorm-gen main missing %q", expected)
		}
	}

	auth := readGeneratedFile(t, root, "app/auth.go")
	for _, forbidden := range []string{"\"errors\"", "\"fmt\""} {
		if strings.Contains(auth, forbidden) {
			t.Fatalf("auth.go must not contain unused import %q", forbidden)
		}
	}

	acl := readGeneratedFile(t, root, "app/acl/pl/health/health.go")
	for _, expected := range []string{
		"import (\n\t\"time\"",
		"CreateTime:  HealthTimeValue(record.CreateTime)",
		"UpdateTime:  HealthTimeValue(record.UpdateTime)",
		"func HealthTimeValue(value *time.Time) time.Time",
	} {
		if !strings.Contains(acl, expected) {
			t.Fatalf("acl adapter missing %q", expected)
		}
	}

	goMod := readGeneratedFile(t, root, "go.mod")
	for _, expected := range []string{
		"go 1.24.3",
		"github.com/redis/go-redis/v9 v9.12.1",
		"gorm.io/driver/postgres v1.6.0",
		"gorm.io/gorm v1.30.0",
		"gorm.io/plugin/dbresolver v1.6.0",
	} {
		if !strings.Contains(goMod, expected) {
			t.Fatalf("go.mod missing %q", expected)
		}
	}

	configYAML := readGeneratedFile(t, root, "configs/config.yml")
	if !strings.Contains(configYAML, "password: 'EVS3%Jthp$Cb7WDkGVzhDgkDhhrNjmfD'") {
		t.Fatalf("config.yml did not preserve quoted password with shell-sensitive characters:\n%s", configYAML)
	}

	configGo := readGeneratedFile(t, root, "configs/config.go")
	for _, expected := range []string{
		"Basic    *config.BasicConfiguration",
		"var basic config.BasicConfiguration",
		"config.SetGlobalConfiguration(cfg.Basic)",
		"pg.Db",
		"pg.SslMode",
	} {
		if !strings.Contains(configGo, expected) {
			t.Fatalf("config.go missing %q", expected)
		}
	}

	infrastructure := readGeneratedFile(t, root, "cmd/server/di/infrastructure.go")
	for _, forbidden := range []string{
		"instance.PostgresClient",
		"instance.RedisClient",
	} {
		if strings.Contains(infrastructure, forbidden) {
			t.Fatalf("infrastructure.go must not contain %q", forbidden)
		}
	}
	for _, expected := range []string{
		"func(cfg *configs.Config) (redis.UniversalClient, error)",
		"cfg.PostgresDSN(\"default\")",
		"gorm.Open(postgres.Open(dsn), &gorm.Config{})",
	} {
		if !strings.Contains(infrastructure, expected) {
			t.Fatalf("infrastructure.go missing %q", expected)
		}
	}

	for _, rel := range []string{
		"cmd/server/di/acl.go",
		"cmd/server/di/domain.go",
		"cmd/server/di/ohs.go",
		"cmd/server/di/routers.go",
		"cmd/server/di/infrastructure.go",
		"cmd/server/di/invoke.go",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("expected layered DI file %s: %v", rel, err)
		}
	}
}

func TestInitDownloadsModules(t *testing.T) {
	root := t.TempDir()
	gen, err := New(&bytes.Buffer{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	oldRunGoModDownload := runGoModDownload
	t.Cleanup(func() { runGoModDownload = oldRunGoModDownload })
	called := false
	runGoModDownload = func(out io.Writer, gotRoot string) error {
		called = true
		if gotRoot != root {
			t.Fatalf("go mod download root = %q, want %q", gotRoot, root)
		}
		return nil
	}

	opts := config.DefaultInitOptions()
	opts.ProjectName = "tidy-check"
	opts.ModulePath = "tidy-check"
	opts.Output = root

	if err := gen.Init(opts); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if !called {
		t.Fatal("expected init to run go mod download")
	}
}

func TestInitRunsGoModTidyAfterGormGen(t *testing.T) {
	root := t.TempDir()
	gen, err := New(&bytes.Buffer{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	oldRunGoModDownload := runGoModDownload
	oldRunGormGen := runGormGen
	oldRunGoModTidy := runGoModTidy
	t.Cleanup(func() {
		runGoModDownload = oldRunGoModDownload
		runGormGen = oldRunGormGen
		runGoModTidy = oldRunGoModTidy
	})

	var calls []string
	runGoModDownload = func(out io.Writer, gotRoot string) error {
		calls = append(calls, "download")
		return nil
	}
	runGormGen = func(out io.Writer, gotRoot string) error {
		calls = append(calls, "gorm-gen")
		return nil
	}
	runGoModTidy = func(out io.Writer, gotRoot string) error {
		calls = append(calls, "tidy")
		return nil
	}

	opts := config.DefaultInitOptions()
	opts.ProjectName = "tidy-check"
	opts.ModulePath = "tidy-check"
	opts.Output = root
	opts.UseExistingDB = true

	if err := gen.Init(opts); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if strings.Join(calls, ",") != "download,gorm-gen,tidy" {
		t.Fatalf("calls = %v, want download, gorm-gen, tidy", calls)
	}
}

func TestInitSkipsGoModCommandsWhenConfigured(t *testing.T) {
	root := t.TempDir()
	gen, err := New(&bytes.Buffer{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	oldRunGoModDownload := runGoModDownload
	oldRunGoModTidy := runGoModTidy
	oldRunGormGen := runGormGen
	t.Cleanup(func() {
		runGoModDownload = oldRunGoModDownload
		runGoModTidy = oldRunGoModTidy
		runGormGen = oldRunGormGen
	})

	var calls []string
	runGoModDownload = func(out io.Writer, gotRoot string) error {
		calls = append(calls, "download")
		return nil
	}
	runGoModTidy = func(out io.Writer, gotRoot string) error {
		calls = append(calls, "tidy")
		return nil
	}
	runGormGen = func(out io.Writer, gotRoot string) error {
		calls = append(calls, "gorm-gen")
		return nil
	}

	opts := config.DefaultInitOptions()
	opts.ProjectName = "skip-check"
	opts.ModulePath = "skip-check"
	opts.Output = root
	opts.UseExistingDB = true
	opts.SkipGoModTidy = true

	if err := gen.Init(opts); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if strings.Join(calls, ",") != "gorm-gen" {
		t.Fatalf("calls = %v, want only gorm-gen", calls)
	}
}

func TestInitDBInfoDoesNotPersistPassword(t *testing.T) {
	root := t.TempDir()
	gen, err := New(&bytes.Buffer{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	oldRunGoModDownload := runGoModDownload
	oldRunGormGen := runGormGen
	oldRunGoModTidy := runGoModTidy
	t.Cleanup(func() {
		runGoModDownload = oldRunGoModDownload
		runGormGen = oldRunGormGen
		runGoModTidy = oldRunGoModTidy
	})
	runGoModDownload = func(out io.Writer, gotRoot string) error { return nil }
	runGormGen = func(out io.Writer, gotRoot string) error { return nil }
	runGoModTidy = func(out io.Writer, gotRoot string) error { return nil }

	opts := config.DefaultInitOptions()
	opts.ProjectName = "db-info-check"
	opts.ModulePath = "db-info-check"
	opts.Output = root
	opts.UseExistingDB = true
	opts.DBPassword = "do-not-save-me"

	if err := gen.Init(opts); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	info := readGeneratedFile(t, root, ".agent/.init-db-info.json")
	if strings.Contains(info, "do-not-save-me") || strings.Contains(info, "\"password\"") {
		t.Fatalf("init db info persisted password:\n%s", info)
	}
	if !strings.Contains(info, "\"passwordStoredInConfig\": true") {
		t.Fatalf("init db info missing passwordStoredInConfig:\n%s", info)
	}
}

func readGeneratedFile(t *testing.T, root, rel string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(content)
}
