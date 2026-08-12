package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type InitOptions struct {
	ProjectName   string
	ModulePath    string
	Output        string
	DB            string
	WithHealth    bool
	RemoveHealth  bool
	AutoGormGen   bool
	UseExistingDB bool
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	SkipGoModTidy bool
	Force         bool
	DryRun        bool
}

type ValidateOptions struct {
	ProjectRoot string
	Strict      bool
	Format      string
}

type AddDomainOptions struct {
	ProjectRoot string
	Domain      string
	Aggregate   string
	Schema      string
	Table       string
	WithCRUD    bool
	Force       bool
	DryRun      bool
}

type AddEntityOptions struct {
	ProjectRoot string
	Domain      string
	Entity      string
	Schema      string
	Table       string
	Force       bool
	DryRun      bool
}

type RemoveHealthOptions struct {
	ProjectRoot string
	DryRun      bool
}

type TemplateData struct {
	ProjectName        string
	ModulePath         string
	DB                 string
	DBHost             string
	DBPort             int
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	Domain             string
	Entity             string
	DomainPascal       string
	DomainPlural       string
	DomainPascalPlural string
	Schema             string
	Table              string
	TablePascal        string
	SQLFileName        string
	GormGenFileName    string
	WithHealth         bool
	ProjectNameYAML    string
	DBHostYAML         string
	DBUserYAML         string
	DBPasswordYAML     string
	DBNameYAML         string
	DBSSLModeYAML      string
}

func DefaultInitOptions() InitOptions {
	return InitOptions{
		DB:         "postgres",
		WithHealth: true,
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBSSLMode:  "disable",
	}
}

func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{ProjectRoot: ".", Format: "text"}
}

func DefaultAddDomainOptions() AddDomainOptions {
	return AddDomainOptions{
		ProjectRoot: ".",
		Schema:      "public",
		Table:       "",
	}
}

func DefaultAddEntityOptions() AddEntityOptions {
	return AddEntityOptions{
		ProjectRoot: ".",
		Schema:      "public",
		Table:       "",
	}
}

func DefaultRemoveHealthOptions() RemoveHealthOptions {
	return RemoveHealthOptions{ProjectRoot: "."}
}

func (o *InitOptions) Normalize() {
	o.ProjectName = strings.TrimSpace(o.ProjectName)
	o.ModulePath = strings.TrimSpace(o.ModulePath)
	o.Output = strings.TrimSpace(o.Output)
	o.DB = strings.TrimSpace(o.DB)
	o.DBHost = strings.TrimSpace(o.DBHost)
	o.DBUser = strings.TrimSpace(o.DBUser)
	o.DBPassword = strings.TrimSpace(o.DBPassword)
	o.DBName = strings.TrimSpace(o.DBName)
	o.DBSSLMode = strings.TrimSpace(o.DBSSLMode)
	if o.DB == "" {
		o.DB = "postgres"
	}
	if o.DBHost == "" {
		o.DBHost = "localhost"
	}
	if o.DBPort == 0 {
		o.DBPort = 5432
	}
	if o.DBUser == "" {
		o.DBUser = "postgres"
	}
	if o.DBPassword == "" {
		o.DBPassword = "postgres"
	}
	if o.DBName == "" {
		o.DBName = o.ProjectName
	}
	if o.DBSSLMode == "" {
		o.DBSSLMode = "disable"
	}
	if o.RemoveHealth {
		o.WithHealth = false
	} else {
		o.WithHealth = true
	}
	if o.Output != "" {
		o.Output = filepath.Clean(o.Output)
	}
}

func (o InitOptions) Validate() error {
	var missing []string
	if o.ProjectName == "" {
		missing = append(missing, "--project-name")
	}
	if o.ModulePath == "" {
		missing = append(missing, "--module")
	}
	if o.Output == "" {
		missing = append(missing, "--output")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
	}
	if o.DB != "postgres" {
		return fmt.Errorf("--db currently supports only postgres, got %q", o.DB)
	}
	if o.AutoGormGen && o.UseExistingDB {
		return errors.New("--auto-gormgen and --use-existing-db are mutually exclusive")
	}
	if o.DBPort < 1 || o.DBPort > 65535 {
		return fmt.Errorf("--db-port must be between 1 and 65535, got %d", o.DBPort)
	}
	return nil
}

func (o *ValidateOptions) Normalize() {
	o.ProjectRoot = strings.TrimSpace(o.ProjectRoot)
	if o.ProjectRoot == "" {
		o.ProjectRoot = "."
	}
	o.ProjectRoot = filepath.Clean(o.ProjectRoot)
	o.Format = strings.TrimSpace(strings.ToLower(o.Format))
	if o.Format == "" {
		o.Format = "text"
	}
}

func (o ValidateOptions) Validate() error {
	if o.ProjectRoot == "" {
		return errors.New("missing --project-root")
	}
	switch o.Format {
	case "text", "json":
		return nil
	default:
		return fmt.Errorf("--format must be text or json, got %q", o.Format)
	}
}

func (o *AddDomainOptions) Normalize() {
	o.ProjectRoot = strings.TrimSpace(o.ProjectRoot)
	if o.ProjectRoot == "" {
		o.ProjectRoot = "."
	}
	o.ProjectRoot = filepath.Clean(o.ProjectRoot)
	o.Domain = NormalizeIdentifier(o.Domain, "")
	o.Schema = NormalizeIdentifier(o.Schema, o.Domain)
	o.Table = NormalizeIdentifier(o.Table, o.Domain)
	o.Aggregate = strings.TrimSpace(o.Aggregate)
	if o.Aggregate == "" && o.Domain != "" {
		o.Aggregate = PascalName(o.Domain)
	}
}

func (o AddDomainOptions) Validate() error {
	if o.Domain == "" {
		return errors.New("missing --domain")
	}
	return nil
}

func (o *AddEntityOptions) Normalize() {
	o.ProjectRoot = strings.TrimSpace(o.ProjectRoot)
	if o.ProjectRoot == "" {
		o.ProjectRoot = "."
	}
	o.ProjectRoot = filepath.Clean(o.ProjectRoot)
	o.Domain = NormalizeIdentifier(o.Domain, "")
	o.Entity = NormalizeIdentifier(o.Entity, o.Domain)
	o.Schema = NormalizeIdentifier(o.Schema, o.Domain)
	o.Table = NormalizeIdentifier(o.Table, o.Entity)
}

func (o AddEntityOptions) Validate() error {
	var missing []string
	if o.Domain == "" {
		missing = append(missing, "--domain")
	}
	if o.Entity == "" {
		missing = append(missing, "--entity")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
	}
	return nil
}

func (o *RemoveHealthOptions) Normalize() {
	o.ProjectRoot = strings.TrimSpace(o.ProjectRoot)
	if o.ProjectRoot == "" {
		o.ProjectRoot = "."
	}
	o.ProjectRoot = filepath.Clean(o.ProjectRoot)
}

func (o RemoveHealthOptions) Validate() error {
	if o.ProjectRoot == "" {
		return errors.New("missing --project-root")
	}
	return nil
}

func (o InitOptions) TemplateData() TemplateData {
	domain := "health"
	schema := domain
	table := domain
	return TemplateData{
		ProjectName:        o.ProjectName,
		ModulePath:         strings.TrimSuffix(o.ModulePath, "/"),
		DB:                 o.DB,
		DBHost:             o.DBHost,
		DBPort:             o.DBPort,
		DBUser:             o.DBUser,
		DBPassword:         o.DBPassword,
		DBName:             o.DBName,
		DBSSLMode:          o.DBSSLMode,
		Domain:             domain,
		Entity:             domain,
		DomainPascal:       PascalName(domain),
		DomainPlural:       PluralName(domain),
		DomainPascalPlural: PascalName(PluralName(domain)),
		Schema:             schema,
		Table:              table,
		TablePascal:        PascalName(table),
		SQLFileName:        SQLFileName(schema, table),
		GormGenFileName:    GormGenFileName(domain, table),
		WithHealth:         o.WithHealth,
		ProjectNameYAML:    YAMLString(o.ProjectName),
		DBHostYAML:         YAMLString(o.DBHost),
		DBUserYAML:         YAMLString(o.DBUser),
		DBPasswordYAML:     YAMLString(o.DBPassword),
		DBNameYAML:         YAMLString(o.DBName),
		DBSSLModeYAML:      YAMLString(o.DBSSLMode),
	}
}

func (o AddDomainOptions) TemplateData(modulePath string) TemplateData {
	domain := NormalizeIdentifier(o.Domain, "")
	schema := NormalizeIdentifier(o.Schema, domain)
	table := NormalizeIdentifier(o.Table, domain)
	aggregate := NormalizePascalIdentifier(o.Aggregate, PascalName(domain))
	return TemplateData{
		ModulePath:         strings.TrimSuffix(strings.TrimSpace(modulePath), "/"),
		DB:                 "postgres",
		Domain:             domain,
		Entity:             domain,
		DomainPascal:       aggregate,
		DomainPlural:       PluralName(domain),
		DomainPascalPlural: aggregate + "s",
		Schema:             schema,
		Table:              table,
		TablePascal:        aggregate,
		SQLFileName:        SQLFileName(schema, table),
		GormGenFileName:    GormGenFileName(domain, table),
		WithHealth:         true,
	}
}

func (o AddEntityOptions) TemplateData(modulePath string) TemplateData {
	domain := NormalizeIdentifier(o.Domain, "")
	entity := NormalizeIdentifier(o.Entity, domain)
	schema := NormalizeIdentifier(o.Schema, domain)
	table := NormalizeIdentifier(o.Table, entity)
	entityPascal := PascalName(entity)
	return TemplateData{
		ModulePath:         strings.TrimSuffix(strings.TrimSpace(modulePath), "/"),
		DB:                 "postgres",
		Domain:             domain,
		Entity:             entity,
		DomainPascal:       entityPascal,
		DomainPlural:       PluralName(entity),
		DomainPascalPlural: PascalName(PluralName(entity)),
		Schema:             schema,
		Table:              table,
		TablePascal:        entityPascal,
		SQLFileName:        SQLFileName(schema, table),
		GormGenFileName:    GormGenFileName(domain, table),
		WithHealth:         true,
	}
}

func SQLFileName(schema, table string) string {
	schema = NormalizeIdentifier(schema, "")
	table = NormalizeIdentifier(table, "")
	if schema == "" || table == "" {
		return table + ".sql"
	}
	return schema + "." + table + ".sql"
}

func GormGenFileName(domain, table string) string {
	domain = NormalizeIdentifier(domain, "")
	table = NormalizeIdentifier(table, domain)
	if table == "" {
		return domain + ".go"
	}
	return domain + "." + table + ".go"
}

func NormalizeIdentifier(input, fallback string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		input = fallback
	}
	var b strings.Builder
	lastUnderscore := false
	for _, r := range input {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore && b.Len() > 0 {
			b.WriteByte('_')
			lastUnderscore = true
		}
	}
	out := strings.Trim(b.String(), "_")
	if out == "" {
		out = fallback
	}
	if out == "" {
		return ""
	}
	if out[0] >= '0' && out[0] <= '9' {
		out = "d_" + out
	}
	return out
}

func NormalizePascalIdentifier(input, fallback string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return fallback
	}
	var b strings.Builder
	upperNext := true
	for _, r := range input {
		switch {
		case r >= 'a' && r <= 'z':
			if upperNext {
				b.WriteRune(r - 'a' + 'A')
			} else {
				b.WriteRune(r)
			}
			upperNext = false
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
			upperNext = false
		case r >= '0' && r <= '9':
			if b.Len() == 0 {
				b.WriteByte('D')
			}
			b.WriteRune(r)
			upperNext = false
		default:
			upperNext = true
		}
	}
	out := b.String()
	if out == "" {
		return fallback
	}
	return out
}

func PascalName(input string) string {
	input = NormalizeIdentifier(input, input)
	if input == "" {
		return ""
	}
	parts := strings.Split(input, "_")
	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		b.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			b.WriteString(part[1:])
		}
	}
	return b.String()
}

func PluralName(input string) string {
	input = NormalizeIdentifier(input, input)
	if input == "" {
		return ""
	}
	if strings.HasSuffix(input, "s") {
		return input
	}
	return input + "s"
}

func YAMLString(input string) string {
	return "'" + strings.ReplaceAll(input, "'", "''") + "'"
}
