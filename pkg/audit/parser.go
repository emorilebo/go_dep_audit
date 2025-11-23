package audit

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Module represents a Go module in the dependency graph
type Module struct {
	Path      string       `json:"Path"`
	Version   string       `json:"Version"`
	Time      string       `json:"Time"`
	Indirect  bool         `json:"Indirect"`
	Dir       string       `json:"Dir"`
	GoMod     string       `json:"GoMod"`
	GoVersion string       `json:"GoVersion"`
	Main      bool         `json:"Main"`
	Replace   *Module      `json:"Replace,omitempty"`
}

// GetModuleGraph returns the full dependency graph using 'go list -m -json all'
func GetModuleGraph(ctx context.Context, projectPath string) ([]Module, error) {
	// Check if go.mod exists
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); os.IsNotExist(err) {
		return nil, fmt.Errorf("go.mod not found in %s", projectPath)
	}

	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-json", "all")
	cmd.Dir = projectPath
	
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run go list: %w, stderr: %s", err, stderr.String())
	}

	var modules []Module
	decoder := json.NewDecoder(&stdout)
	for decoder.More() {
		var mod Module
		if err := decoder.Decode(&mod); err != nil {
			return nil, fmt.Errorf("failed to decode module json: %w", err)
		}
		modules = append(modules, mod)
	}

	return modules, nil
}

// ParseGoMod parses the go.mod file directly (fallback or for direct deps only)
// Note: This is a simple parser and doesn't handle complex replace directives or transitive deps
func ParseGoMod(path string) ([]Module, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var modules []Module
	scanner := bufio.NewScanner(file)
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if line == "require (" {
			inRequire = true
			continue
		}
		if line == ")" && inRequire {
			inRequire = false
			continue
		}

		if strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				modules = append(modules, parseModuleLine(parts[1:]))
			}
		} else if inRequire {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				modules = append(modules, parseModuleLine(parts))
			}
		}
	}

	return modules, scanner.Err()
}

func parseModuleLine(parts []string) Module {
	mod := Module{
		Path:    parts[0],
		Version: parts[1],
	}
	// Check for // indirect comment
	for _, p := range parts {
		if p == "//" {
			continue
		}
		if p == "indirect" {
			mod.Indirect = true
		}
	}
	return mod
}
