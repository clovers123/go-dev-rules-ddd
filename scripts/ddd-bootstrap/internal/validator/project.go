package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"ddd-bootstrap/internal/config"
)

type Result struct {
	Level  string `json:"level"`
	Rule   string `json:"rule"`
	File   string `json:"file"`
	Reason string `json:"reason"`
	Fix    string `json:"fix,omitempty"`
}

func (r Result) String() string {
	line := fmt.Sprintf("%s %s %s %s", r.Level, r.Rule, filepath.ToSlash(r.File), r.Reason)
	if r.Fix != "" && r.Level != "PASS" {
		line += " | fix: " + r.Fix
	}
	return line
}

func Validate(opts config.ValidateOptions) ([]Result, bool, error) {
	opts.Normalize()
	if err := opts.Validate(); err != nil {
		return nil, false, err
	}

	var results []Result
	add := func(level, rule, file, reason, fix string) {
		results = append(results, Result{
			Level:  level,
			Rule:   rule,
			File:   file,
			Reason: reason,
			Fix:    fix,
		})
	}

	checkCoreProject(opts.ProjectRoot, add)
	checkNoApplicationLayer(opts.ProjectRoot, add)
	checkLegacyGuidance(opts.ProjectRoot, add)
	checkAgentDispatch(opts.ProjectRoot, add)
	checkGitignore(opts.ProjectRoot, add)
	flows := detectGeneratedFlows(opts.ProjectRoot, add)
	checkGeneratedFlows(opts.ProjectRoot, flows, add)
	checkImports(opts.ProjectRoot, add)
	checkGormGen(opts.ProjectRoot, flows, add)
	checkRepositoryRules(opts.ProjectRoot, add)

	failed := false
	for _, r := range results {
		if r.Level == "FAIL" || opts.Strict && r.Level == "WARN" {
			failed = true
			break
		}
	}
	return results, failed, nil
}

func checkCoreProject(root string, add resultAdder) {
	requiredDirs := []string{
		"app/ohs",
		"app/ohs/local/appservice",
		"app/ohs/remote/controller",
		"app/ohs/remote/routers",
		"app/domain",
		"app/acl",
		"app/infra",
		"app/ohs/pl",
		"app/domain",
		"app/acl",
		"cmd/server/di",
		"cmd/gorm-gen",
		"sql",
		"configs",
		".agent/rules",
	}
	for _, dir := range uniqueStrings(requiredDirs) {
		checkDir(root, dir, "project.core-dir", add)
	}

	diFiles := []string{"acl.go", "domain.go", "infrastructure.go", "invoke.go", "ohs.go", "routers.go"}
	for _, name := range diFiles {
		checkFile(root, filepath.Join("cmd/server/di", name), "project.di-file", add)
	}
}

func checkNoApplicationLayer(root string, add resultAdder) {
	rel := "app/application"
	path := filepath.Join(root, filepath.FromSlash(rel))
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		add("FAIL", "project.no-application-layer", rel, "standalone application layer exists", "move application services to app/ohs/local/appservice and remove app/application")
		return
	}
	add("PASS", "project.no-application-layer", rel, "standalone application layer is absent", "")
}

func checkLegacyGuidance(root string, add resultAdder) {
	legacyDirs := []string{
		"spec-rules",
		".agent/spec-rules",
		".agent/dir-rules",
	}
	for _, dir := range legacyDirs {
		path := filepath.Join(root, filepath.FromSlash(dir))
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			add("FAIL", "agent.no-legacy-rules", dir, "legacy rule directory exists", "remove legacy guidance directory; keep active guidance in .agent/rules, AGENTS.md, and directory .agent.md files")
		} else {
			add("PASS", "agent.no-legacy-rules", dir, "not present", "")
		}
	}
}

func checkAgentDispatch(root string, add resultAdder) {
	content := readOptional(root, "AGENTS.md")
	if strings.TrimSpace(content) == "" {
		add("FAIL", "agents.exists", "AGENTS.md", "AGENTS.md is missing or empty", "generate AGENTS.md with ddd-bootstrap init")
	} else {
		add("PASS", "agents.exists", "AGENTS.md", "AGENTS.md is populated", "")
		for _, expected := range []string{"add-domain", "add-entity", ".agent/rules", ".agent.md", "validate", "health", ".gen.go", "owned_by", "is_deleted"} {
			if strings.Contains(content, expected) {
				add("PASS", "agents.dispatch", "AGENTS.md", expected+" is referenced", "")
			} else {
				add("FAIL", "agents.dispatch", "AGENTS.md", expected+" is missing", "regenerate AGENTS.md")
			}
		}
	}

	for _, file := range requiredRuleFiles() {
		checkFile(root, filepath.Join(".agent/rules", file), "agent.rules-file", add)
	}
	for _, file := range requiredAgentDocs() {
		checkFile(root, file, "agent.dir-doc", add)
		content := readOptional(root, file)
		if content == "" {
			continue
		}
		if strings.Contains(content, ".agent/rules/") && strings.Contains(content, "Allowed") && strings.Contains(content, "Forbidden") {
			add("PASS", "agent.dir-doc-content", file, "role guidance references rules and allowed/forbidden actions", "")
		} else {
			add("WARN", "agent.dir-doc-content", file, "directory agent doc is missing rules or allowed/forbidden guidance", "regenerate directory .agent.md")
		}
	}
}

func requiredRuleFiles() []string {
	return []string{
		"database-system-fields.md",
		"ddd-adapter.md",
		"ddd-layer-mapping.md",
		"di-registration.md",
		"feature-field-usage.md",
		"git-commit-convention.md",
		"ohs-routing.md",
		"repository-dao.md",
		"sql-file-naming.md",
	}
}

func requiredAgentDocs() []string {
	return []string{
		"app/ohs/.agent.md",
		"app/ohs/remote/.agent.md",
		"app/ohs/remote/controller/.agent.md",
		"app/ohs/remote/routers/.agent.md",
		"app/ohs/local/appservice/.agent.md",
		"app/ohs/pl/.agent.md",
		"app/domain/.agent.md",
		"app/acl/.agent.md",
		"app/acl/port/.agent.md",
		"app/acl/pl/.agent.md",
		"app/acl/adapter/.agent.md",
		"app/acl/adapter/repository/postgres/.agent.md",
		"app/infra/.agent.md",
		"cmd/server/di/.agent.md",
		"cmd/gorm-gen/.agent.md",
		"sql/.agent.md",
	}
}

func checkGitignore(root string, add resultAdder) {
	path := filepath.Join(root, ".gitignore")
	content, err := os.ReadFile(path)
	if err != nil {
		add("FAIL", "gitignore.exists", ".gitignore", ".gitignore is missing", "generate .gitignore with ddd-bootstrap init")
		return
	}
	text := string(content)
	for _, entry := range []string{"**/.agent.md", ".agent/.init-db-info.json"} {
		if strings.Contains(text, entry) {
			add("PASS", "gitignore.agent-entry", ".gitignore", entry+" is ignored", "")
		} else {
			add("FAIL", "gitignore.agent-entry", ".gitignore", entry+" is missing", "add "+entry+" to .gitignore")
		}
	}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == ".agent/" || line == ".agent/*" {
			add("WARN", "gitignore.agent-rules", ".gitignore", line+" may hide shared .agent/rules", "ignore only generated/sensitive .agent files, not .agent/rules/")
			break
		}
	}
}

type generatedFlow struct {
	Domain   string
	Entity   string
	Table    string
	Schema   string
	TypeName string
	SQLFile  string
	GormFile string
}

func detectGeneratedFlows(root string, add resultAdder) []generatedFlow {
	base := filepath.Join(root, "cmd/gorm-gen")
	var flows []generatedFlow
	_ = filepath.WalkDir(base, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || entry.Name() == "main.go" {
			return nil
		}
		rel, relErr := filepath.Rel(base, path)
		if relErr != nil {
			return nil
		}
		parts := strings.Split(filepath.ToSlash(rel), "/")
		if len(parts) != 2 {
			return nil
		}
		domain := parts[0]
		stem := strings.TrimSuffix(parts[1], ".go")
		nameParts := strings.SplitN(stem, ".", 2)
		entity := nameParts[0]
		table := entity
		if len(nameParts) == 2 {
			table = nameParts[1]
		}
		content := readOptional(root, filepath.ToSlash(filepath.Join("cmd/gorm-gen", rel)))
		typeName := detectInterfaceType(content, entity)
		entity = detectEntityFile(root, domain, typeName, entity)
		schema, sqlFile := detectSQLSchemaAndFile(root, table, domain)
		flows = append(flows, generatedFlow{
			Domain:   domain,
			Entity:   entity,
			Table:    table,
			Schema:   schema,
			TypeName: typeName,
			SQLFile:  sqlFile,
			GormFile: filepath.ToSlash(filepath.Join("cmd/gorm-gen", rel)),
		})
		return nil
	})
	sort.Slice(flows, func(i, j int) bool {
		if flows[i].Domain == flows[j].Domain {
			return flows[i].Table < flows[j].Table
		}
		return flows[i].Domain < flows[j].Domain
	})
	if len(flows) == 0 {
		add("WARN", "generated-flow.detect", "cmd/gorm-gen", "no generated entity/domain flow detected", "run init or add-domain/add-entity")
	}
	return flows
}

func detectEntityFile(root string, domain string, typeName string, fallback string) string {
	base := filepath.Join(root, "app", "domain", domain, "entity")
	result := fallback
	_ = filepath.WalkDir(base, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil || !regexp.MustCompile(`type\s+`+regexp.QuoteMeta(typeName)+`\s+struct`).Match(content) {
			return nil
		}
		result = strings.TrimSuffix(entry.Name(), ".go")
		return filepath.SkipAll
	})
	return result
}

func detectInterfaceType(content, fallback string) string {
	re := regexp.MustCompile(`type\s+([A-Za-z][A-Za-z0-9]*)\s+interface`)
	matches := re.FindStringSubmatch(content)
	if len(matches) == 2 {
		return strings.TrimSuffix(matches[1], "Querier")
	}
	return config.PascalName(fallback)
}

func detectSQLSchemaAndFile(root, table, fallback string) (string, string) {
	sqlRoot := filepath.Join(root, "sql")
	var schema string
	var sqlFile string
	_ = filepath.WalkDir(sqlRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			return nil
		}
		if entry.Name() != table+".sql" && !strings.HasSuffix(entry.Name(), "."+table+".sql") {
			return nil
		}
		parent := filepath.Base(filepath.Dir(path))
		stem := strings.TrimSuffix(entry.Name(), ".sql")
		prefix := strings.TrimSuffix(stem, "."+table)
		if prefix != "" {
			schema = prefix
		} else {
			schema = parent
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr == nil {
			sqlFile = filepath.ToSlash(rel)
		}
		return filepath.SkipAll
	})
	if schema == "" {
		schema = fallback
	}
	if sqlFile == "" {
		sqlFile = fmt.Sprintf("sql/%s/%s.%s.sql", schema, schema, table)
	}
	return schema, sqlFile
}

func checkGeneratedFlows(root string, flows []generatedFlow, add resultAdder) {
	for _, flow := range flows {
		persistentFiles := []string{
			fmt.Sprintf("app/domain/%s/entity/%s.go", flow.Domain, flow.Entity),
			fmt.Sprintf("app/acl/port/repository/%s/%s.go", flow.Domain, flow.Entity),
			fmt.Sprintf("app/acl/pl/%s/%s.go", flow.Domain, flow.Entity),
			fmt.Sprintf("app/acl/adapter/repository/postgres/%s/%s.go", flow.Domain, flow.Entity),
			fmt.Sprintf("app/acl/adapter/repository/postgres/repository/export.%s.%s.go", flow.Schema, flow.Table),
			flow.GormFile,
			flow.SQLFile,
		}
		for _, file := range persistentFiles {
			checkFile(root, file, "generated-flow.persistent-file", add)
		}

		checkSQLSystemFields(root, flow, add)
		checkAdapterShape(root, flow, add)
		checkDIRegistration(root, flow, add)

		if shouldHaveFullDomain(root, flow) {
			checkFullDomainFiles(root, flow, add)
			checkFullDomainDI(root, flow, add)
		}
	}
}

func shouldHaveFullDomain(root string, flow generatedFlow) bool {
	if flow.Entity == flow.Domain {
		return true
	}
	candidates := []string{
		fmt.Sprintf("app/domain/%s/service/%s.go", flow.Domain, flow.Entity),
		fmt.Sprintf("app/ohs/local/appservice/%s/%s.go", flow.Domain, flow.Entity),
		fmt.Sprintf("app/ohs/remote/controller/%s/%s.go", flow.Domain, flow.Entity),
	}
	for _, rel := range candidates {
		if content := readOptional(root, rel); content != "" {
			return true
		}
	}
	return false
}

func checkFullDomainFiles(root string, flow generatedFlow, add resultAdder) {
	name := flow.Entity
	files := []string{
		fmt.Sprintf("app/domain/%s/service/%s.go", flow.Domain, name),
		fmt.Sprintf("app/ohs/pl/%s/create_%s.go", flow.Domain, name),
		fmt.Sprintf("app/ohs/pl/%s/list_%s.go", flow.Domain, config.PluralName(name)),
		fmt.Sprintf("app/ohs/pl/%s/adapter/%s.go", flow.Domain, name),
		fmt.Sprintf("app/ohs/local/appservice/%s/%s.go", flow.Domain, name),
		fmt.Sprintf("app/ohs/remote/controller/%s/%s.go", flow.Domain, name),
		fmt.Sprintf("app/ohs/remote/routers/%s.go", flow.Domain),
	}
	for _, file := range files {
		checkFile(root, file, "generated-flow.full-domain-file", add)
	}
}

func checkSQLSystemFields(root string, flow generatedFlow, add resultAdder) {
	rel := flow.SQLFile
	content := readOptional(root, rel)
	if content == "" {
		return
	}
	required := []string{"owned_by", "is_deleted", "create_time", "update_time", "delete_time", "created_by", "updated_by", "deleted_by", "feature"}
	for _, field := range required {
		if strings.Contains(content, field) {
			add("PASS", "sql.system-field", rel, field+" exists", "")
		} else {
			add("FAIL", "sql.system-field", rel, field+" is missing", "add required system field")
		}
	}
	if strings.Contains(strings.ToUpper(content), "FOREIGN KEY") {
		add("FAIL", "sql.no-foreign-key", rel, "foreign key constraint exists", "remove database FK and enforce consistency in domain code")
	} else {
		add("PASS", "sql.no-foreign-key", rel, "no foreign key constraint", "")
	}
}

func checkAdapterShape(root string, flow generatedFlow, add resultAdder) {
	aclPath := fmt.Sprintf("app/acl/pl/%s/%s.go", flow.Domain, flow.Entity)
	aclContent := readOptional(root, aclPath)
	if aclContent != "" {
		for _, expected := range []string{"I" + flow.TypeName + "AclAdapter", "To" + flow.TypeName + "Aggregations", "From" + flow.TypeName + "AggregationToInfo"} {
			if strings.Contains(aclContent, expected) {
				add("PASS", "adapter.acl-shape", aclPath, expected+" exists", "")
			} else {
				add("FAIL", "adapter.acl-shape", aclPath, expected+" is missing", "regenerate ACL adapter")
			}
		}
	}

	repoPath := fmt.Sprintf("app/acl/adapter/repository/postgres/%s/%s.go", flow.Domain, flow.Entity)
	repoContent := readOptional(root, repoPath)
	if repoContent != "" {
		if strings.Contains(repoContent, "adapter") && strings.Contains(repoContent, "r.adapter.") {
			add("PASS", "adapter.repository-delegates", repoPath, "repository delegates mapping to ACL adapter", "")
		} else {
			add("FAIL", "adapter.repository-delegates", repoPath, "repository mapping does not use ACL adapter", "delegate entity/model conversion to ACL adapter")
		}
	}
}

func checkDIRegistration(root string, flow generatedFlow, add resultAdder) {
	aclContent := readOptional(root, "cmd/server/di/acl.go")
	expectations := []string{"New" + flow.TypeName + "AclAdapter", "New" + flow.TypeName + "Repository"}
	for _, expected := range expectations {
		if strings.Contains(aclContent, expected) {
			add("PASS", "di.acl-registration", "cmd/server/di/acl.go", expected+" is registered", "")
		} else {
			add("FAIL", "di.acl-registration", "cmd/server/di/acl.go", expected+" is missing", "register ACL adapter/repository in acl.go")
		}
	}

	gormMain := readOptional(root, "cmd/gorm-gen/main.go")
	if containsApplyInterfaceRegistration(gormMain, flow) && strings.Contains(gormMain, fmt.Sprintf(`GenerateModelAs("%s", "%s"`, flow.Table, flow.TypeName)) {
		add("PASS", "gormgen.registration", "cmd/gorm-gen/main.go", flow.TypeName+" GORM Gen interface is registered", "")
	} else {
		add("FAIL", "gormgen.registration", "cmd/gorm-gen/main.go", flow.TypeName+" GORM Gen interface is missing", "register ApplyInterface in cmd/gorm-gen/main.go")
	}
}

func containsApplyInterfaceRegistration(content string, flow generatedFlow) bool {
	pattern := fmt.Sprintf(`func\s*\(\s*%s\.%s\s*\)`, regexp.QuoteMeta(flow.Domain), regexp.QuoteMeta(flow.TypeName))
	return regexp.MustCompile(pattern).FindStringIndex(content) != nil
}

func checkFullDomainDI(root string, flow generatedFlow, add resultAdder) {
	domainContent := readOptional(root, "cmd/server/di/domain.go")
	ohsContent := readOptional(root, "cmd/server/di/ohs.go")
	routerContent := readOptional(root, "cmd/server/di/routers.go")
	checks := []struct {
		file     string
		content  string
		rule     string
		expected string
		fix      string
	}{
		{"cmd/server/di/domain.go", domainContent, "di.domain-registration", "New" + flow.TypeName + "DomainService", "register domain service in domain.go"},
		{"cmd/server/di/ohs.go", ohsContent, "di.ohs-registration", "New" + flow.TypeName + "Adapter", "register OHS adapter in ohs.go"},
		{"cmd/server/di/ohs.go", ohsContent, "di.ohs-registration", "New" + flow.TypeName + "AppService", "register appservice in ohs.go"},
		{"cmd/server/di/ohs.go", ohsContent, "di.ohs-registration", "New" + flow.TypeName + "Controller", "register controller in ohs.go"},
		{"cmd/server/di/routers.go", routerContent, "di.router-registration", "Mount" + flow.TypeName + "RouteGroup", "invoke router mount in routers.go"},
	}
	for _, check := range checks {
		if strings.Contains(check.content, check.expected) {
			add("PASS", check.rule, check.file, check.expected+" is registered", "")
		} else {
			add("FAIL", check.rule, check.file, check.expected+" is missing", check.fix)
		}
	}
}

func checkDir(root, rel, rule string, add resultAdder) {
	path := filepath.Join(root, filepath.FromSlash(rel))
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		add("PASS", rule, rel, "exists", "")
		return
	}
	add("FAIL", rule, rel, "directory is missing", "create "+rel)
}

func checkFile(root, rel, rule string, add resultAdder) {
	path := filepath.Join(root, filepath.FromSlash(rel))
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		add("PASS", rule, rel, "exists", "")
		return
	}
	add("FAIL", rule, rel, "file is missing", "create "+rel)
}

func uniqueStrings(in []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, value := range in {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

type resultAdder func(level, rule, file, reason, fix string)
