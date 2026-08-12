package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	{{if .WithHealth}}
	"{{.ModulePath}}/cmd/gorm-gen/{{.Domain}}"
	{{end}}
	// ddd-bootstrap:import gorm-gen

	"github.com/gookit/goutil/strutil"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gopkg.in/yaml.v3"
)

func main() {
	dsn, err := loadPostgresDSN("configs/config.yml", "default")
	if err != nil {
		log.Fatalf("failed to build postgres dsn: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	executeSQLFiles(db)

	g := gen.NewGenerator(gen.Config{
		OutPath:           "./app/acl/adapter/repository/postgres/repository",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		WithUnitTest:      false,
	})

	g.UseDB(db)

	// JSON tag strategy: camelCase
	g.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
		return strutil.CamelCase(columnName)
	})

	{{if .WithHealth}}
	// ==================== {{.DomainPascal}} Domain ====================
	g.WithTableNameStrategy(func(tableName string) (tableContent string) {
		return fmt.Sprintf("{{.Schema}}.%s", tableName)
	})

	g.ApplyInterface(func({{.Domain}}.{{.TablePascal}}) {},
		g.GenerateModelAs("{{.Table}}", "{{.TablePascal}}",
			gen.FieldType("id", "string"),
			gen.FieldType("owned_by", "string"),
			gen.FieldType("created_by", "*string"),
			gen.FieldType("updated_by", "*string"),
			gen.FieldType("deleted_by", "*string"),
			gen.FieldType("create_time", "*time.Time"),
			gen.FieldType("update_time", "*time.Time"),
			gen.FieldType("feature", "datatypes.JSON"),
			gen.FieldType("is_deleted", "bool")),
	)
	{{end}}
	// ddd-bootstrap:apply-interface

	g.Execute()
	fmt.Println("GORM Gen completed successfully")
}

type gormGenConfig struct {
	Postgres map[string]postgresConfig `yaml:"postgres"`
}

type postgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslMode"`
	DB       string `yaml:"db"`
}

func loadPostgresDSN(path string, name string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var cfg gormGenConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	if name == "" {
		name = "default"
	}
	pg, ok := cfg.Postgres[name]
	if !ok {
		return "", fmt.Errorf("postgres config %q not found", name)
	}
	if pg.SSLMode == "" {
		pg.SSLMode = "disable"
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pg.Host,
		pg.Port,
		pg.User,
		pg.Password,
		pg.DB,
		pg.SSLMode,
	), nil
}

func executeSQLFiles(db *gorm.DB) {
	files, err := filepath.Glob("sql/**/*.sql")
	if err != nil {
		log.Fatalf("find sql files: %v", err)
	}
	sort.Strings(files)
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("read sql file %s: %v", file, err)
		}
		sql := strings.TrimSpace(string(content))
		if sql == "" {
			continue
		}
		for _, statement := range strings.Split(sql, ";") {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}
			if err := db.Exec(statement).Error; err != nil {
				log.Fatalf("execute sql file %s: %v", file, err)
			}
		}
	}
}
