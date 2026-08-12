package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"ddd-bootstrap/internal/config"
	"ddd-bootstrap/internal/generator"
	"ddd-bootstrap/internal/validator"
)

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "add-domain":
		return runAddDomain(args[1:], stdout, stderr)
	case "add-entity":
		return runAddEntity(args[1:], stdout, stderr)
	case "remove-health":
		return runRemoveHealth(args[1:], stdout, stderr)
	case "-h", "--help", "help":
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func runInit(args []string, stdout, stderr io.Writer) int {
	opts := config.DefaultInitOptions()
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.ProjectName, "project-name", opts.ProjectName, "project name")
	fs.StringVar(&opts.ModulePath, "module", opts.ModulePath, "Go module path")
	fs.StringVar(&opts.Output, "output", opts.Output, "output directory")
	fs.StringVar(&opts.DB, "db", opts.DB, "database type")
	fs.BoolVar(&opts.AutoGormGen, "auto-gormgen", false, "start a temporary PostgreSQL container and run GORM Gen")
	fs.BoolVar(&opts.UseExistingDB, "use-existing-db", false, "use the configured PostgreSQL database and run GORM Gen")
	fs.BoolVar(&opts.RemoveHealth, "remove-health", false, "skip the bootstrap health module after verification")
	fs.BoolVar(&opts.SkipGoModTidy, "skip-go-mod-tidy", false, "skip go mod download and go mod tidy")
	fs.StringVar(&opts.DBHost, "db-host", opts.DBHost, "PostgreSQL host")
	fs.IntVar(&opts.DBPort, "db-port", opts.DBPort, "PostgreSQL port")
	fs.StringVar(&opts.DBUser, "db-user", opts.DBUser, "PostgreSQL user")
	fs.StringVar(&opts.DBPassword, "db-password", opts.DBPassword, "PostgreSQL password")
	fs.StringVar(&opts.DBName, "db-name", opts.DBName, "PostgreSQL database name")
	fs.StringVar(&opts.DBSSLMode, "db-sslmode", opts.DBSSLMode, "PostgreSQL sslmode")
	fs.BoolVar(&opts.Force, "force", false, "overwrite managed files")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "print planned files without writing")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	gen, err := generator.New(stdout)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := gen.Init(opts); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	opts := config.DefaultValidateOptions()
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.ProjectRoot, "project-root", opts.ProjectRoot, "project root")
	fs.BoolVar(&opts.Strict, "strict", opts.Strict, "treat WARN as FAIL")
	fs.StringVar(&opts.Format, "format", opts.Format, "output format: text or json")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	opts.Normalize()

	results, failed, err := validator.Validate(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if opts.Format == "json" {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(results); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	} else {
		for _, result := range results {
			fmt.Fprintln(stdout, result.String())
		}
	}
	if failed {
		return 1
	}
	return 0
}

func runAddDomain(args []string, stdout, stderr io.Writer) int {
	opts := config.DefaultAddDomainOptions()
	fs := flag.NewFlagSet("add-domain", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.ProjectRoot, "project-root", opts.ProjectRoot, "project root")
	fs.StringVar(&opts.Domain, "domain", opts.Domain, "domain package name")
	fs.StringVar(&opts.Aggregate, "aggregate", opts.Aggregate, "aggregate name")
	fs.StringVar(&opts.Schema, "schema", opts.Schema, "database schema")
	fs.StringVar(&opts.Table, "table", opts.Table, "database table")
	fs.BoolVar(&opts.WithCRUD, "with-crud", false, "generate CRUD chain")
	fs.BoolVar(&opts.Force, "force", false, "overwrite managed files")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "print planned files without writing")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	gen, err := generator.New(stdout)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := gen.AddDomain(opts); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func runAddEntity(args []string, stdout, stderr io.Writer) int {
	opts := config.DefaultAddEntityOptions()
	fs := flag.NewFlagSet("add-entity", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.ProjectRoot, "project-root", opts.ProjectRoot, "project root")
	fs.StringVar(&opts.Domain, "domain", opts.Domain, "domain package name")
	fs.StringVar(&opts.Entity, "entity", opts.Entity, "entity or value object package/file name")
	fs.StringVar(&opts.Schema, "schema", opts.Schema, "database schema")
	fs.StringVar(&opts.Table, "table", opts.Table, "database table")
	fs.BoolVar(&opts.Force, "force", false, "overwrite managed files")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "print planned files without writing")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	gen, err := generator.New(stdout)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := gen.AddEntity(opts); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func runRemoveHealth(args []string, stdout, stderr io.Writer) int {
	opts := config.DefaultRemoveHealthOptions()
	fs := flag.NewFlagSet("remove-health", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.ProjectRoot, "project-root", opts.ProjectRoot, "project root")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "print planned removals without deleting")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	gen, err := generator.New(stdout)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := gen.RemoveHealth(opts); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: ddd-bootstrap <init|validate|add-domain|add-entity|remove-health> [flags]")
}
