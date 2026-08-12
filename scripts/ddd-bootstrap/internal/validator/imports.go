package validator

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func checkImports(root string, add resultAdder) {
	_ = filepath.WalkDir(filepath.Join(root, "app"), func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() {
			if entry.Name() == "vendor" || entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		imports, parseErr := importsOf(path)
		if parseErr != nil {
			add("WARN", "imports.parse", rel, "could not parse imports", parseErr.Error())
			return nil
		}
		checkImportRules(rel, imports, add)
		return nil
	})
}

func importsOf(path string) ([]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	imports := make([]string, 0, len(file.Imports))
	for _, spec := range file.Imports {
		imports = append(imports, strings.Trim(spec.Path.Value, `"`))
	}
	return imports, nil
}

func checkImportRules(rel string, imports []string, add resultAdder) {
	hasResponseWrapper := false
	for _, imported := range imports {
		if imported == "git.yugeeker.com/SHARED/go-lazy/app/entities" {
			hasResponseWrapper = true
		}
		switch {
		case strings.HasPrefix(rel, "app/domain/") && imported == "github.com/gofiber/fiber/v3":
			add("FAIL", "imports.domain-fiber", rel, "domain layer imports Fiber", "move protocol handling to controller or OHS PL adapter")
		case strings.HasPrefix(rel, "app/domain/") && strings.Contains(imported, "/app/acl/adapter/repository/postgres/model"):
			add("FAIL", "imports.domain-postgres-model", rel, "domain layer imports postgres model", "map persistence models in ACL PL")
		case strings.HasPrefix(rel, "app/domain/") && strings.Contains(imported, "/app/ohs/"):
			add("FAIL", "imports.domain-ohs", rel, "domain layer imports OHS", "move protocol types and use-case orchestration to app/ohs")
		case strings.HasPrefix(rel, "app/domain/") && strings.Contains(imported, "/app/acl/adapter"):
			add("FAIL", "imports.domain-acl-adapter", rel, "domain layer imports an ACL implementation", "depend only on app/acl/port interfaces")
		case strings.HasPrefix(rel, "app/ohs/remote/controller/") && strings.Contains(imported, "/app/acl/adapter"):
			add("FAIL", "imports.controller-acl-adapter", rel, "controller imports ACL adapter", "depend on appservice and OHS PL adapter only")
		case strings.HasPrefix(rel, "app/ohs/remote/controller/") && strings.Contains(imported, "/app/domain/") && strings.Contains(imported, "/service"):
			add("FAIL", "imports.controller-domain-service", rel, "controller imports domain service directly", "call appservice instead")
		case strings.HasPrefix(rel, "app/acl/pl/") && strings.Contains(imported, "/app/ohs/remote/controller"):
			add("FAIL", "imports.acl-pl-controller", rel, "ACL PL imports controller", "keep protocol types out of ACL PL")
		case strings.HasPrefix(rel, "app/ohs/local/appservice/") && strings.Contains(imported, "/app/acl/adapter"):
			add("FAIL", "imports.appservice-acl-adapter", rel, "appservice imports an ACL implementation", "depend on app/acl/port interfaces")
		}
	}
	if strings.HasPrefix(rel, "app/ohs/remote/controller/") && !hasResponseWrapper {
		add("FAIL", "imports.controller-response-wrapper", rel, "controller does not import shared response wrapper", "use git.yugeeker.com/SHARED/go-lazy/app/entities")
	}
}
