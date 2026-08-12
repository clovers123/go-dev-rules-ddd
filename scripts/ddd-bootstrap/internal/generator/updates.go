package generator

import (
	"bufio"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type goFileUpdate struct {
	path           string
	imports        []string
	importMarker   string
	entries        []string
	entryMarker    string
	force          bool
	dryRun         bool
	actionLabel    string
	alreadyMessage string
}

func readModulePath(root string) (string, error) {
	file, err := os.Open(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[0] == "module" {
			return fields[1], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan go.mod: %w", err)
	}
	return "", fmt.Errorf("module path not found in go.mod")
}

func filepathJoin(root, rel string) string {
	return filepath.Join(root, filepath.FromSlash(rel))
}

func applyGoFileUpdate(out io.Writer, update goFileUpdate) error {
	contentBytes, err := os.ReadFile(update.path)
	if err != nil {
		return fmt.Errorf("read %s: %w", filepath.ToSlash(update.path), err)
	}
	content := string(contentBytes)

	missingImports := missingImportSpecs(content, update.imports)
	missingEntries := missingEntries(content, update.entries)
	if len(missingImports) == 0 && len(missingEntries) == 0 {
		fmt.Fprintf(out, "%s %s\n", update.alreadyMessage, filepath.ToSlash(update.path))
		return nil
	}

	if update.dryRun {
		fmt.Fprintf(out, "%s %s\n", update.actionLabel, filepath.ToSlash(update.path))
		return nil
	}

	var changed bool
	if len(missingImports) > 0 {
		var inserted bool
		content, inserted, err = insertBeforeMarker(content, update.importMarker, missingImports)
		if err != nil {
			return fmt.Errorf("update imports in %s: %w", filepath.ToSlash(update.path), err)
		}
		changed = changed || inserted
	}
	if len(missingEntries) > 0 {
		var inserted bool
		content, inserted, err = insertBeforeMarker(content, update.entryMarker, missingEntries)
		if err != nil {
			return fmt.Errorf("update entries in %s: %w", filepath.ToSlash(update.path), err)
		}
		changed = changed || inserted
	}
	if !changed {
		fmt.Fprintf(out, "%s %s\n", update.alreadyMessage, filepath.ToSlash(update.path))
		return nil
	}

	formatted, err := format.Source([]byte(content))
	if err != nil {
		return fmt.Errorf("format updated Go file %s: %w", filepath.ToSlash(update.path), err)
	}
	if err := os.WriteFile(update.path, formatted, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filepath.ToSlash(update.path), err)
	}
	fmt.Fprintf(out, "%s %s\n", update.actionLabel, filepath.ToSlash(update.path))
	return nil
}

func missingImportSpecs(content string, specs []string) []string {
	var missing []string
	for _, spec := range specs {
		path := quotedImportPath(spec)
		if path == "" {
			continue
		}
		if strings.Contains(content, "\""+path+"\"") {
			continue
		}
		missing = append(missing, spec)
	}
	return missing
}

func quotedImportPath(spec string) string {
	first := strings.IndexByte(spec, '"')
	last := strings.LastIndexByte(spec, '"')
	if first < 0 || last <= first {
		return ""
	}
	return spec[first+1 : last]
}

func missingEntries(content string, entries []string) []string {
	var missing []string
	for _, entry := range entries {
		key := entryKey(entry)
		if key == "" || strings.Contains(content, key) {
			continue
		}
		missing = append(missing, entry)
	}
	return missing
}

func entryKey(entry string) string {
	for _, line := range strings.Split(entry, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "GenerateModelAs(") {
			return strings.TrimSuffix(line, ",")
		}
	}
	for _, line := range strings.Split(entry, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		return strings.TrimSuffix(line, ",")
	}
	return strings.TrimSpace(entry)
}

func insertBeforeMarker(content, marker string, entries []string) (string, bool, error) {
	if marker == "" {
		return content, false, fmt.Errorf("marker is empty")
	}
	idx := strings.Index(content, marker)
	if idx < 0 {
		return content, false, fmt.Errorf("marker %q not found", marker)
	}
	lineStart := strings.LastIndex(content[:idx], "\n") + 1
	lineEnd := idx + strings.Index(content[idx:], "\n")
	if lineEnd < idx {
		lineEnd = len(content)
	}
	markerLine := content[lineStart:lineEnd]
	indent := leadingWhitespace(markerLine)

	var b strings.Builder
	for _, entry := range entries {
		for _, line := range strings.Split(entry, "\n") {
			if strings.TrimSpace(line) == "" {
				b.WriteByte('\n')
				continue
			}
			b.WriteString(indent)
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}

	return content[:lineStart] + b.String() + content[lineStart:], true, nil
}

func leadingWhitespace(line string) string {
	var b strings.Builder
	for _, r := range line {
		if r != ' ' && r != '\t' {
			break
		}
		b.WriteRune(r)
	}
	return b.String()
}
