package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkGormGen(root string, flows []generatedFlow, add resultAdder) {
	checkGormGenMain(root, add)
	checkGormGenOutputPaths(root, add)
	checkGeneratedFileHeaders(root, add)
	for _, flow := range flows {
		tablePascal := flow.TypeName
		checkGeneratedDAOFiles(root, flow, add)

		exportPath := fmt.Sprintf("app/acl/adapter/repository/postgres/repository/export.%s.%s.go", flow.Schema, flow.Table)
		exportContent := readOptional(root, exportPath)
		if exportContent != "" {
			expectations := []string{
				"func Get" + tablePascal + "Dao(ctx context.Context, tx *Query)",
				"if tx != nil",
				"return tx." + tablePascal + ".WithContext(ctx)",
				"return Q." + tablePascal + ".WithContext(ctx)",
			}
			for _, expected := range expectations {
				if strings.Contains(exportContent, expected) {
					add("PASS", "gormgen.wrapper", exportPath, expected+" exists", "")
				} else {
					add("FAIL", "gormgen.wrapper", exportPath, expected+" is missing", "regenerate DAO wrapper")
				}
			}
		}

		interfacePath := flow.GormFile
		interfaceContent := readOptional(root, interfacePath)
		if interfaceContent != "" {
			if strings.Contains(interfaceContent, "Querier interface") {
				add("FAIL", "gormgen.interface-name", interfacePath, "interface name uses Querier suffix", "name the interface after the generated model type, e.g. Health or AppInfo")
			} else {
				add("PASS", "gormgen.interface-name", interfacePath, "interface name has no Querier suffix", "")
			}
			if strings.Contains(interfaceContent, "owned_by") && strings.Contains(interfaceContent, "is_deleted=false") {
				add("PASS", "gormgen.sql-tenant-soft-delete", interfacePath, "owned_by and is_deleted=false constraints exist", "")
			} else {
				add("FAIL", "gormgen.sql-tenant-soft-delete", interfacePath, "tenant or soft delete constraint is missing", "add owned_by and is_deleted=false to query comments")
			}
		}

		implPath := fmt.Sprintf("app/acl/adapter/repository/postgres/%s/%s.go", flow.Domain, flow.Entity)
		implContent := readOptional(root, implPath)
		if implContent != "" {
			if strings.Contains(implContent, "repository.Get"+tablePascal+"Dao(ctx, tx)") || strings.Contains(implContent, "repository.Get"+tablePascal+"Dao(ctx, nil)") {
				add("PASS", "gormgen.tx-wrapper", implPath, "repository impl uses the DAO wrapper", "")
			} else {
				add("FAIL", "gormgen.tx-wrapper", implPath, "repository impl does not use the DAO wrapper", "call repository.Get"+tablePascal+"Dao(ctx, tx) or repository.Get"+tablePascal+"Dao(ctx, nil)")
			}
			if strings.Contains(implContent, "repository.Q.") && !strings.Contains(implContent, "ddd-bootstrap: allow-global-dao") {
				add("FAIL", "gormgen.no-global-dao", implPath, "repository impl calls repository.Q directly", "use repository.Get"+tablePascal+"Dao(ctx, tx)")
			} else {
				add("PASS", "gormgen.no-global-dao", implPath, "repository impl avoids repository.Q direct calls", "")
			}
			if strings.Contains(implContent, "owned_by") && strings.Contains(implContent, "is_deleted") {
				add("PASS", "repository.tenant-soft-delete", implPath, "repository keeps tenant and soft delete constraints visible", "")
			} else {
				add("WARN", "repository.tenant-soft-delete", implPath, "repository relies on GORM Gen SQL for tenant and soft delete constraints", "keep owned_by and is_deleted filters in query comments")
			}
		}
	}
}

func checkGormGenMain(root string, add resultAdder) {
	const rel = "cmd/gorm-gen/main.go"
	content := readOptional(root, rel)
	if content == "" {
		add("FAIL", "gormgen.main", rel, "GORM Gen entrypoint is missing", "regenerate project with ddd-bootstrap init")
		return
	}
	if strings.Contains(content, "config.Load(") ||
		strings.Contains(content, "config.Postgres.DSN") ||
		strings.Contains(content, "/configs\"") {
		add("FAIL", "gormgen.config-loader", rel, "GORM Gen entrypoint depends on project configs API", "regenerate cmd/gorm-gen/main.go from the bootstrap template; it must load configs/config.yml itself")
		return
	}
	if strings.Contains(content, "ModelPkgPath") {
		add("FAIL", "gormgen.model-path", rel, "GORM Gen entrypoint sets ModelPkgPath and can nest model output under the repository path", "remove ModelPkgPath; let GORM Gen create the sibling postgres/model directory")
		return
	}
	if strings.Contains(content, "loadPostgresDSN(") && strings.Contains(content, "gorm.Open(postgres.Open(dsn)") {
		add("PASS", "gormgen.config-loader", rel, "GORM Gen entrypoint builds a DSN before opening postgres", "")
	} else {
		add("WARN", "gormgen.config-loader", rel, "GORM Gen entrypoint does not match the bootstrap DSN loader shape", "keep cmd/gorm-gen/main.go self-contained and avoid project configs API coupling")
	}
}

func checkGormGenOutputPaths(root string, add resultAdder) {
	const rel = "app/acl/adapter/repository/postgres/app"
	path := filepath.Join(root, filepath.FromSlash(rel))
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		add("FAIL", "gormgen.model-path", rel, "generated model output is nested under postgres/app", "delete the nested generated output, remove ModelPkgPath, and rerun go run ./cmd/gorm-gen from the project root")
		return
	}
	add("PASS", "gormgen.model-path", rel, "no nested postgres/app generated output", "")
}

func checkGeneratedDAOFiles(root string, flow generatedFlow, add resultAdder) {
	files := []string{
		fmt.Sprintf("app/acl/adapter/repository/postgres/repository/%s.%s.gen.go", flow.Schema, flow.Table),
		fmt.Sprintf("app/acl/adapter/repository/postgres/model/%s.%s.gen.go", flow.Schema, flow.Table),
	}
	for _, file := range files {
		content := readOptional(root, file)
		if content == "" {
			add("WARN", "gormgen.generated-file", file, "generated DAO/model file is missing", "run go run ./cmd/gorm-gen")
			continue
		}
		if strings.Contains(content, "Code generated") && strings.Contains(content, "DO NOT EDIT") {
			add("PASS", "gormgen.generated-file", file, "generated file exists and has generated header", "")
		} else {
			add("FAIL", "gormgen.generated-file", file, "file looks hand-written or generated header is missing", "regenerate with go run ./cmd/gorm-gen; do not hand-write .gen.go files")
		}
	}
}

func checkGeneratedFileHeaders(root string, add resultAdder) {
	base := filepath.Join(root, "app/acl/adapter/repository/postgres")
	_ = filepath.WalkDir(base, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(entry.Name(), ".gen.go") {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return nil
		}
		file := filepath.ToSlash(rel)
		content := readOptional(root, file)
		if strings.Contains(content, "Code generated") && strings.Contains(content, "DO NOT EDIT") {
			add("PASS", "gormgen.no-handwritten-gen", file, "generated header exists", "")
		} else {
			add("FAIL", "gormgen.no-handwritten-gen", file, ".gen.go file is missing generated header", "delete the hand-written file and rerun GORM Gen")
		}
		return nil
	})
}

func readOptional(root, rel string) string {
	content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return ""
	}
	return string(content)
}
