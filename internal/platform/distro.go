package platform

import (
	"os"
	"strings"

	"github.com/mdgspace/sysreplicate/internal/domain"
)

var knownDistros = map[string]string{
	"debian":       "debian",
	"ubuntu":       "debian",
	"linuxmint":    "debian",
	"pop":          "debian",
	"arch":         "arch",
	"manjaro":      "arch",
	"endeavouros":  "arch",
	"rhel":         "rhel",
	"fedora":       "fedora",
	"centos":       "rhel",
	"rocky":        "rhel",
	"alma":         "rhel",
	"void":         "void",
	"opensuse":     "opensuse",
	"opensuse-leap": "opensuse",
	"opensuse-tumbleweed": "opensuse",
	"suse":         "opensuse",
	"alpine":       "alpine",
	"nixos":        "nixos",
	"gentoo":       "gentoo",
}

func normalizeDistro(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func resolveBaseDistro(candidates []string) string {
	for _, c := range candidates {
		normalized := normalizeDistro(c)
		if base, ok := knownDistros[normalized]; ok {
			return base
		}
	}
	return ""
}

func extractField(data, prefix string) string {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return strings.Trim(strings.SplitN(line, "=", 2)[1], `"`)
		}
	}
	return ""
}

func readOSRelease() (string, error) {
	paths := []string{domain.OsReleasePath, "/etc/lsb-release", "/etc/issue"}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		return string(data), nil
	}
	return "", os.ErrNotExist
}

// DetectDistro returns the distro name and the base distro for package management.
func DetectDistro() (string, string) {
	data, err := readOSRelease()
	if err != nil {
		return "unknown", "unknown"
	}

	distro := normalizeDistro(extractField(data, "ID="))
	idLike := extractField(data, "ID_LIKE=")

	candidates := []string{distro}
	candidates = append(candidates, strings.Fields(idLike)...)

	baseDistro := resolveBaseDistro(candidates)
	if baseDistro == "" {
		baseDistro = "unknown"
	}

	return distro, baseDistro
}

func IsImmutableDistro() bool {
	data, err := readOSRelease()
	if err != nil {
		return false
	}

	id := normalizeDistro(extractField(data, "ID="))
	variant := normalizeDistro(extractField(data, "VARIANT_ID="))

	switch {
	case variant == "silverblue" || variant == "kinoite":
		return true
	case strings.Contains(id, "steamos"):
		return true
	case strings.Contains(id, "opensuse-microos"):
		return true
	}

	return false
}
