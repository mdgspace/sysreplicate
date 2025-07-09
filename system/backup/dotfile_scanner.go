package backup

import (
	"os"
	"path/filepath"
	"strings"
)

var DotfilePaths = []string{
	"~/.bashrc",
	"~/.zshrc",
	"~/.vimrc",
	"~/.config",
	"~/.bash_history",
	"~/.zsh_history",
	"~/.gitconfig",
	"~/.profile",
	"~/.npmrc",
}

type Dotfile struct {
	Path     string
	RelPath  string
	IsDir    bool
	IsBinary bool
	Mode     os.FileMode
	Content  string // ignore for the binary files
}

// expand ~ to home dir
func expandHome(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
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
func ScanDotfiles() ([]Dotfile, error) {
	var results []Dotfile
	home, _ := os.UserHomeDir()

	for _, raw := range DotfilePaths {
		full := expandHome(raw)

		info, err := os.Stat(full)
		if err != nil {
			continue
		}

		relPath, _ := filepath.Rel(home, full)
		entry := Dotfile{
			Path:    full,
			RelPath: relPath,
			IsDir:   info.IsDir(),
			Mode:    info.Mode(),
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
