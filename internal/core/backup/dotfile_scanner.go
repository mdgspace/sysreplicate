package backup

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mdgspace/sysreplicate/internal/domain"
)

// expand ~ to home dir
func expandHome(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			if len(path) == 1 {
				return home
			}
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

// check for binary files
func containsNullByte(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return true
		}
	}
	return false
}

// ScanDotfiles scans all dotfiles and returns their metadata + content
func ScanDotfiles() ([]domain.Dotfile, error) {
	var results []domain.Dotfile
	home, _ := os.UserHomeDir()

	for _, raw := range domain.DotfilePaths {
		full := expandHome(raw)

		info, err := os.Stat(full)
		if err != nil {
			continue
		}

		realPath, _ := filepath.Rel(home, full)
		entry := domain.Dotfile{
			Path:     full,
			RealPath: realPath,
			IsDir:    info.IsDir(),
			Mode:     info.Mode(),
		}

		if !info.IsDir() {
			data, err := os.ReadFile(full)
			if err != nil {
				continue
			}
			if containsNullByte(data) {
				entry.IsBinary = true
			} else {
				entry.Content = string(data)
			}
		}

		results = append(results, entry)
	}

	return results, nil
}
