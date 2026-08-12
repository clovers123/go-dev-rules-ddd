package generator

import (
	"fmt"
	"io"

	"ddd-bootstrap/internal/config"
)

func (g *Generator) AddDomain(opts config.AddDomainOptions) error {
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		return err
	}

	modulePath, err := readModulePath(opts.ProjectRoot)
	if err != nil {
		return err
	}
	data := opts.TemplateData(modulePath)
	specs := domainFiles(data)
	files, err := g.renderSpecs(specs, data)
	if err != nil {
		return err
	}
	if err := applyFilePlan(g.out, planOptions{
		Root:   opts.ProjectRoot,
		Force:  opts.Force,
		DryRun: opts.DryRun,
	}, domainDirs(data), files); err != nil {
		return err
	}
	if err := updateDomainRegistrations(g.out, opts, data); err != nil {
		return err
	}
	if opts.DryRun {
		return nil
	}
	return runGormGenIfConfigured(g.out, opts.ProjectRoot)
}

func (g *Generator) AddEntity(opts config.AddEntityOptions) error {
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		return err
	}

	modulePath, err := readModulePath(opts.ProjectRoot)
	if err != nil {
		return err
	}
	data := opts.TemplateData(modulePath)
	specs := entityFiles(data)
	files, err := g.renderSpecs(specs, data)
	if err != nil {
		return err
	}
	if err := applyFilePlan(g.out, planOptions{
		Root:   opts.ProjectRoot,
		Force:  opts.Force,
		DryRun: opts.DryRun,
	}, entityDirs(data), files); err != nil {
		return err
	}
	if err := updateEntityRegistrations(g.out, opts, data); err != nil {
		return err
	}
	if opts.DryRun {
		return nil
	}
	return runGormGenIfConfigured(g.out, opts.ProjectRoot)
}

func domainFiles(data config.TemplateData) []fileSpec {
	d := data.Domain
	table := data.Table
	entity := data.Entity
	schema := data.Schema
	return []fileSpec{
		{template: "domain/entity.go.tpl", target: fmt.Sprintf("app/domain/%s/entity/%s.go", d, entity)},
		{template: "domain/service.go.tpl", target: fmt.Sprintf("app/domain/%s/service/%s.go", d, d)},
		{template: "domain/ohs_create.go.tpl", target: fmt.Sprintf("app/ohs/pl/%s/create_%s.go", d, d)},
		{template: "domain/ohs_list.go.tpl", target: fmt.Sprintf("app/ohs/pl/%s/list_%s.go", d, data.DomainPlural)},
		{template: "domain/ohs_adapter.go.tpl", target: fmt.Sprintf("app/ohs/pl/%s/adapter/%s.go", d, d)},
		{template: "domain/appservice.go.tpl", target: fmt.Sprintf("app/ohs/local/appservice/%s/%s.go", d, entity)},
		{template: "domain/controller.go.tpl", target: fmt.Sprintf("app/ohs/remote/controller/%s/%s.go", d, entity)},
		{template: "domain/router.go.tpl", target: fmt.Sprintf("app/ohs/remote/routers/%s.go", d)},
		{template: "domain/repository_port.go.tpl", target: fmt.Sprintf("app/acl/port/repository/%s/%s.go", d, d)},
		{template: "domain/acl_adapter.go.tpl", target: fmt.Sprintf("app/acl/pl/%s/%s.go", d, d)},
		{template: "domain/repository_impl.go.tpl", target: fmt.Sprintf("app/acl/adapter/repository/postgres/%s/%s.go", d, d)},
		{template: "domain/gormgen_interface.go.tpl", target: fmt.Sprintf("cmd/gorm-gen/%s/%s", d, data.GormGenFileName)},
		{template: "domain/repository_export.go.tpl", target: fmt.Sprintf("app/acl/adapter/repository/postgres/repository/export.%s.%s.go", schema, table)},
		{template: "domain/sql.sql.tpl", target: fmt.Sprintf("sql/%s/%s", schema, data.SQLFileName)},
	}
}

func entityFiles(data config.TemplateData) []fileSpec {
	d := data.Domain
	entity := data.Entity
	schema := data.Schema
	return []fileSpec{
		{template: "domain/entity_sub.go.tpl", target: fmt.Sprintf("app/domain/%s/entity/%s.go", d, entity)},
		{template: "domain/repository_port.go.tpl", target: fmt.Sprintf("app/acl/port/repository/%s/%s.go", d, entity)},
		{template: "domain/acl_adapter.go.tpl", target: fmt.Sprintf("app/acl/pl/%s/%s.go", d, entity)},
		{template: "domain/repository_impl.go.tpl", target: fmt.Sprintf("app/acl/adapter/repository/postgres/%s/%s.go", d, entity)},
		{template: "domain/gormgen_interface.go.tpl", target: fmt.Sprintf("cmd/gorm-gen/%s/%s", d, data.GormGenFileName)},
		{template: "domain/repository_export.go.tpl", target: fmt.Sprintf("app/acl/adapter/repository/postgres/repository/export.%s.%s.go", schema, data.Table)},
		{template: "domain/sql.sql.tpl", target: fmt.Sprintf("sql/%s/%s", schema, data.SQLFileName)},
	}
}

func gormGenFiles(data config.TemplateData) []fileSpec {
	return nil
}

func domainDirs(data config.TemplateData) []string {
	return []string{
		"app/domain/" + data.Domain + "/entity",
		"app/domain/" + data.Domain + "/valueobject",
		"app/domain/" + data.Domain + "/service",
		"app/domain/" + data.Domain + "/factory",
		"app/domain/" + data.Domain + "/event",
		"app/ohs/local/appservice/" + data.Domain,
		"app/ohs/remote/controller/" + data.Domain,
		"app/ohs/remote/routers",
		"app/ohs/pl/" + data.Domain + "/adapter",
		"app/acl/port/repository/" + data.Domain,
		"app/acl/pl/" + data.Domain,
		"app/acl/adapter/repository/postgres/" + data.Domain,
		"cmd/gorm-gen/" + data.Domain,
		"sql/" + data.Schema,
	}
}

func entityDirs(data config.TemplateData) []string {
	return []string{
		"app/domain/" + data.Domain + "/entity",
		"app/acl/port/repository/" + data.Domain,
		"app/acl/pl/" + data.Domain,
		"app/acl/adapter/repository/postgres/" + data.Domain,
		"cmd/gorm-gen/" + data.Domain,
		"sql/" + data.Schema,
	}
}

func updateDomainRegistrations(out io.Writer, opts config.AddDomainOptions, data config.TemplateData) error {
	if err := updateEntityRegistrations(out, config.AddEntityOptions{
		ProjectRoot: opts.ProjectRoot,
		Domain:      data.Domain,
		Entity:      data.Entity,
		Schema:      data.Schema,
		Table:       data.Table,
		Force:       opts.Force,
		DryRun:      opts.DryRun,
	}, data); err != nil {
		return err
	}
	updates := []goFileUpdate{
		{
			path: filepathJoin(opts.ProjectRoot, "cmd/server/di/domain.go"),
			imports: []string{
				fmt.Sprintf("%sservice \"%s/app/domain/%s/service\"", data.Domain, data.ModulePath, data.Domain),
			},
			importMarker:   "ddd-bootstrap:import domain",
			entries:        []string{fmt.Sprintf("%sservice.New%sDomainService,", data.Domain, data.DomainPascal)},
			entryMarker:    "ddd-bootstrap:provide domain",
			force:          opts.Force,
			dryRun:         opts.DryRun,
			actionLabel:    "UPDATE",
			alreadyMessage: "SKIP",
		},
		{
			path: filepathJoin(opts.ProjectRoot, "cmd/server/di/ohs.go"),
			imports: []string{
				fmt.Sprintf("%sappservice \"%s/app/ohs/local/appservice/%s\"", data.Domain, data.ModulePath, data.Domain),
				fmt.Sprintf("%sadapter \"%s/app/ohs/pl/%s/adapter\"", data.Domain, data.ModulePath, data.Domain),
				fmt.Sprintf("%scontroller \"%s/app/ohs/remote/controller/%s\"", data.Domain, data.ModulePath, data.Domain),
			},
			importMarker: "ddd-bootstrap:import ohs",
			entries: []string{
				fmt.Sprintf("%sadapter.New%sAdapter,", data.Domain, data.DomainPascal),
				fmt.Sprintf("%sappservice.New%sAppService,", data.Domain, data.DomainPascal),
				fmt.Sprintf("%scontroller.New%sController,", data.Domain, data.DomainPascal),
			},
			entryMarker:    "ddd-bootstrap:provide ohs",
			force:          opts.Force,
			dryRun:         opts.DryRun,
			actionLabel:    "UPDATE",
			alreadyMessage: "SKIP",
		},
		{
			path: filepathJoin(opts.ProjectRoot, "cmd/server/di/routers.go"),
			imports: []string{
				fmt.Sprintf("\"%s/app/ohs/remote/routers\"", data.ModulePath),
			},
			importMarker:   "ddd-bootstrap:import routers",
			entries:        []string{fmt.Sprintf("routers.Mount%sRouteGroup,", data.DomainPascal)},
			entryMarker:    "ddd-bootstrap:invoke routers",
			force:          opts.Force,
			dryRun:         opts.DryRun,
			actionLabel:    "UPDATE",
			alreadyMessage: "SKIP",
		},
	}
	for _, update := range updates {
		if err := applyGoFileUpdate(out, update); err != nil {
			return err
		}
	}
	return nil
}

func updateEntityRegistrations(out io.Writer, opts config.AddEntityOptions, data config.TemplateData) error {
	updates := []goFileUpdate{
		{
			path: filepathJoin(opts.ProjectRoot, "cmd/server/di/acl.go"),
			imports: []string{
				fmt.Sprintf("%sacl \"%s/app/acl/pl/%s\"", data.Domain, data.ModulePath, data.Domain),
				fmt.Sprintf("%srepo \"%s/app/acl/adapter/repository/postgres/%s\"", data.Domain, data.ModulePath, data.Domain),
			},
			importMarker: "ddd-bootstrap:import acl",
			entries: []string{
				fmt.Sprintf("%sacl.New%sAclAdapter,", data.Domain, data.DomainPascal),
				fmt.Sprintf("%srepo.New%sRepository,", data.Domain, data.DomainPascal),
			},
			entryMarker:    "ddd-bootstrap:provide acl",
			force:          opts.Force,
			dryRun:         opts.DryRun,
			actionLabel:    "UPDATE",
			alreadyMessage: "SKIP",
		},
		{
			path: filepathJoin(opts.ProjectRoot, "cmd/gorm-gen/main.go"),
			imports: []string{
				fmt.Sprintf("\"%s/cmd/gorm-gen/%s\"", data.ModulePath, data.Domain),
			},
			importMarker: "ddd-bootstrap:import gorm-gen",
			entries: []string{
				gormGenApplyBlock(data),
			},
			entryMarker:    "ddd-bootstrap:apply-interface",
			force:          opts.Force,
			dryRun:         opts.DryRun,
			actionLabel:    "UPDATE",
			alreadyMessage: "SKIP",
		},
	}
	for _, update := range updates {
		if err := applyGoFileUpdate(out, update); err != nil {
			return err
		}
	}
	return nil
}

func gormGenApplyBlock(data config.TemplateData) string {
	return fmt.Sprintf(`// ==================== %s %s ====================
	g.WithTableNameStrategy(func(tableName string) (tableContent string) {
		return fmt.Sprintf("%s.%%s", tableName)
	})
	g.ApplyInterface(func(%s.%s) {},
		g.GenerateModelAs("%s", "%s",
			gen.FieldType("id", "string"),
			gen.FieldType("owned_by", "string"),
			gen.FieldType("created_by", "*string"),
			gen.FieldType("updated_by", "*string"),
			gen.FieldType("deleted_by", "*string"),
			gen.FieldType("create_time", "*time.Time"),
			gen.FieldType("update_time", "*time.Time"),
			gen.FieldType("feature", "datatypes.JSON"),
			gen.FieldType("is_deleted", "bool")),
	)`, data.DomainPascal, data.Table, data.Schema, data.Domain, data.TablePascal, data.Table, data.TablePascal)
}
