package utils

import (
	"os"
	"strings"
	"github.com/mdgspace/sysreplicate/internal/domain"
)

// DetectDistro returns the distro and base distro.
func DetectDistro() (string, string) {
	data, err := os.ReadFile(domain.OsReleasePath)
	if err != nil {
		return "unknown", "unknown"
	}
	var distro, baseDistro string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			distro = strings.Trim(strings.SplitN(line, "=", 2)[1], `"`)
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			baseDistro = strings.Trim(strings.SplitN(line, "=", 2)[1], `"`)
		}
	}
	return distro, baseDistro
}