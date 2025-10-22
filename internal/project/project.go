package project

import (
	"os"
	"path/filepath"
	"strings"
)

func FindProjectRoot(startPath string) (string, string, error) {
	dir, err := filepath.Abs(startPath)
	if err != nil {
		return "", "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			content, err := os.ReadFile(filepath.Join(dir, "go.mod"))
			if err != nil {
				return dir, "", err
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "module ") {
					moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module"))
					return dir, moduleName, nil
				}
			}
			return dir, "", nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return startPath, "", nil
		}
		dir = parent
	}
}
