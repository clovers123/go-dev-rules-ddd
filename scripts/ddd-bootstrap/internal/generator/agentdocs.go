package generator

import "ddd-bootstrap/internal/config"

func agentDocFiles() []fileSpec {
	return []fileSpec{
		{template: "agent-docs/ohs.agent.md.tpl", target: "app/ohs/.agent.md"},
		{template: "agent-docs/ohs-remote.agent.md.tpl", target: "app/ohs/remote/.agent.md"},
		{template: "agent-docs/controller.agent.md.tpl", target: "app/ohs/remote/controller/.agent.md"},
		{template: "agent-docs/routers.agent.md.tpl", target: "app/ohs/remote/routers/.agent.md"},
		{template: "agent-docs/appservice.agent.md.tpl", target: "app/ohs/local/appservice/.agent.md"},
		{template: "agent-docs/ohs-pl.agent.md.tpl", target: "app/ohs/pl/.agent.md"},
		{template: "agent-docs/domain.agent.md.tpl", target: "app/domain/.agent.md"},
		{template: "agent-docs/acl.agent.md.tpl", target: "app/acl/.agent.md"},
		{template: "agent-docs/acl-port.agent.md.tpl", target: "app/acl/port/.agent.md"},
		{template: "agent-docs/acl-pl.agent.md.tpl", target: "app/acl/pl/.agent.md"},
		{template: "agent-docs/acl-adapter.agent.md.tpl", target: "app/acl/adapter/.agent.md"},
		{template: "agent-docs/postgres-repository.agent.md.tpl", target: "app/acl/adapter/repository/postgres/.agent.md"},
		{template: "agent-docs/infra.agent.md.tpl", target: "app/infra/.agent.md"},
		{template: "agent-docs/di.agent.md.tpl", target: "cmd/server/di/.agent.md"},
		{template: "agent-docs/gorm-gen.agent.md.tpl", target: "cmd/gorm-gen/.agent.md"},
		{template: "agent-docs/sql.agent.md.tpl", target: "sql/.agent.md"},
	}
}

func ruleFiles() []fileSpec {
	return []fileSpec{
		{template: "rules/database-system-fields.md.tpl", target: ".agent/rules/database-system-fields.md"},
		{template: "rules/ddd-adapter.md.tpl", target: ".agent/rules/ddd-adapter.md"},
		{template: "rules/ddd-layer-mapping.md.tpl", target: ".agent/rules/ddd-layer-mapping.md"},
		{template: "rules/di-registration.md.tpl", target: ".agent/rules/di-registration.md"},
		{template: "rules/feature-field-usage.md.tpl", target: ".agent/rules/feature-field-usage.md"},
		{template: "rules/git-commit-convention.md.tpl", target: ".agent/rules/git-commit-convention.md"},
		{template: "rules/ohs-routing.md.tpl", target: ".agent/rules/ohs-routing.md"},
		{template: "rules/repository-dao.md.tpl", target: ".agent/rules/repository-dao.md"},
		{template: "rules/sql-file-naming.md.tpl", target: ".agent/rules/sql-file-naming.md"},
	}
}

func projectDirs(data config.TemplateData, includeDomain bool) []string {
	dirs := []string{
		"app/ohs/local/appservice",
		"app/ohs/remote/controller",
		"app/ohs/remote/routers",
		"app/ohs/remote/middleware",
		"app/ohs/pl",
		"app/domain",
		"app/acl/port/repository",
		"app/acl/port/client",
		"app/acl/port/publisher",
		"app/acl/port/subscriber",
		"app/acl/port/security",
		"app/acl/pl",
		"app/acl/adapter/repository/postgres/model",
		"app/acl/adapter/repository/postgres/repository",
		"app/infra/constants",
		"app/infra/enums",
		"app/infra/cron",
		"app/infra/syserrors",
		"app/infra/validation",
		"app/infra/abstract",
		"cmd/server/di",
		"cmd/gorm-gen",
		"configs",
		"sql",
		"tools",
		"test",
		".agent",
		".agent/rules",
	}
	if !includeDomain {
		return dirs
	}

	d := data.Domain
	return append(dirs,
		"app/domain/"+d+"/entity",
		"app/domain/"+d+"/service",
		"app/ohs/pl/"+d+"/adapter",
		"app/ohs/local/appservice/"+d,
		"app/ohs/remote/controller/"+d,
		"app/acl/port/repository/"+d,
		"app/acl/pl/"+d,
		"app/acl/adapter/repository/postgres/"+d,
		"cmd/gorm-gen/"+d,
		"sql/"+data.Schema,
	)
}
