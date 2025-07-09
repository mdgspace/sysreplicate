package backup

import (
	"os"
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
	Path  string
	Found bool
}

// Expand ~ to home directory
func expandHome(path string) string {
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			return strings.Replace(path, "~", home, 1)
		}
	}
	return path
}

// Check which dotfiles exist
func CheckDotfiles() ([]Dotfile, error) {
	var result []Dotfile

	for _, rawPath := range DotfilePaths {
		fullPath := expandHome(rawPath)

		_, err := os.Stat(fullPath)
		result = append(result, Dotfile{
			Path:  fullPath,
			Found: err == nil,
		})
	}
	return result, nil
}
