package validator

import (
	"os"
	"path/filepath"
	"strings"
)

func checkRepositoryRules(root string, add resultAdder) {
	base := filepath.Join(root, "app/acl/adapter/repository/postgres")
	rawChains := []string{".Where(", ".Order(", ".Limit(", ".Offset(", ".Like(", ".Or(", ".Joins(", ".Raw(", ".Exec("}
	_ = filepath.WalkDir(base, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() {
			switch entry.Name() {
			case "model", "repository":
				return filepath.SkipDir
			default:
				return nil
			}
		}
		if !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		contentBytes, readErr := os.ReadFile(path)
		if readErr != nil {
			add("WARN", "repository.read", rel, "could not read repository impl", readErr.Error())
			return nil
		}
		content := string(contentBytes)
		for _, chain := range rawChains {
			if strings.Contains(content, chain) {
				add("FAIL", "repository.no-raw-gorm", rel, "repository impl uses "+chain, "move query logic to cmd/gorm-gen SQL comments")
			}
		}
		if strings.Contains(content, "repository.Q.") && !strings.Contains(content, "ddd-bootstrap: allow-global-dao") {
			add("FAIL", "repository.no-global-dao", rel, "repository impl calls repository.Q directly", "use repository.GetXDao(ctx, tx)")
		}
		if strings.Contains(content, "&model.") || strings.Contains(content, "model.") {
			add("FAIL", "repository.no-model-mapping", rel, "repository impl maps postgres model directly", "delegate model/entity mapping to ACL adapter")
		} else {
			add("PASS", "repository.no-model-mapping", rel, "repository impl avoids direct model mapping", "")
		}
		if strings.Contains(content, "BO") || strings.Contains(content, "DTO") {
			add("FAIL", "repository.no-dto-bo", rel, "repository impl mentions DTO/BO types", "repository ports should return domain entities")
		}
		return nil
	})
}
