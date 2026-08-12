package generator

import (
	"io"

	"ddd-bootstrap/internal/config"
	rendertemplates "ddd-bootstrap/internal/templates"
)

type Generator struct {
	renderer *rendertemplates.Renderer
	out      io.Writer
}

func New(out io.Writer) (*Generator, error) {
	renderer, err := rendertemplates.NewRenderer("")
	if err != nil {
		return nil, err
	}
	return &Generator{renderer: renderer, out: out}, nil
}

func (g *Generator) Init(opts config.InitOptions) error {
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		return err
	}

	dbInfo, err := prepareInitDBInfo(&opts)
	if err != nil {
		return err
	}

	data := opts.TemplateData()
	specs := projectFiles()
	specs = append(specs, cmdServerFiles()...)
	specs = append(specs, gormGenFiles(data)...)
	if opts.WithHealth {
		specs = append(specs, domainFiles(data)...)
	}
	specs = append(specs, agentDocFiles()...)
	specs = append(specs, ruleFiles()...)

	files, err := g.renderSpecs(specs, data)
	if err != nil {
		return err
	}
	if dbInfo != nil {
		infoFile, err := renderInitDBInfo(dbInfo)
		if err != nil {
			return err
		}
		files = append(files, infoFile)
	}
	if err := applyPlan(g.out, opts, projectDirs(data, opts.WithHealth), files); err != nil {
		return err
	}
	if opts.DryRun {
		return nil
	}
	if opts.WithHealth {
		if err := updateDomainRegistrations(g.out, config.AddDomainOptions{
			ProjectRoot: opts.Output,
			Domain:      data.Domain,
			Aggregate:   data.DomainPascal,
			Schema:      data.Schema,
			Table:       data.Table,
			Force:       opts.Force,
		}, data); err != nil {
			return err
		}
	}
	if !opts.SkipGoModTidy {
		if err := runGoModDownload(g.out, opts.Output); err != nil {
			return err
		}
	}
	if dbInfo == nil {
		return nil
	}
	if err := activateInitDB(g.out, opts.Output, dbInfo); err != nil {
		return err
	}
	if err := runGormGen(g.out, opts.Output); err != nil {
		return err
	}
	if !opts.SkipGoModTidy {
		if err := runGoModTidy(g.out, opts.Output); err != nil {
			return err
		}
	}
	return nil
}

func projectFiles() []fileSpec {
	return []fileSpec{
		{template: "project/go.mod.tpl", target: "go.mod"},
		{template: "project/Makefile.tpl", target: "Makefile"},
		{template: "project/gitignore.tpl", target: ".gitignore"},
		{template: "project/AGENTS.md.tpl", target: "AGENTS.md"},
		{template: "project/README.md.tpl", target: "README.md"},
		{template: "project/config.go.tpl", target: "configs/config.go"},
		{template: "project/config.yml.tpl", target: "configs/config.yml"},
		{template: "project/gorm-gen-main.go.tpl", target: "cmd/gorm-gen/main.go"},
		{template: "project/auth.go.tpl", target: "app/auth.go"},
	}
}

func cmdServerFiles() []fileSpec {
	return []fileSpec{
		{template: "cmd-server/main.go.tpl", target: "cmd/server/main.go"},
		{template: "cmd-server/di/infrastructure.go.tpl", target: "cmd/server/di/infrastructure.go"},
		{template: "cmd-server/di/acl.go.tpl", target: "cmd/server/di/acl.go"},
		{template: "cmd-server/di/domain.go.tpl", target: "cmd/server/di/domain.go"},
		{template: "cmd-server/di/ohs.go.tpl", target: "cmd/server/di/ohs.go"},
		{template: "cmd-server/di/routers.go.tpl", target: "cmd/server/di/routers.go"},
		{template: "cmd-server/di/invoke.go.tpl", target: "cmd/server/di/invoke.go"},
	}
}
