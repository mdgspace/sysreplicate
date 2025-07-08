package output

import (
	"encoding/json"
)

// BuildSystemJSON creates a well-structured JSON object for the system info and packages.
func BuildSystemJSON(osType, distro, baseDistro string, packages map[string][]string) ([]byte, error) {

	type SystemInfo struct {
		OS         string              `json:"os"`
		Distro     string              `json:"distro"`
		BaseDistro string              `json:"base_distro"`
		Packages   map[string][]string `json:"packages"`
	}

	info := SystemInfo{
		OS:         osType,
		Distro:     distro,
		BaseDistro: baseDistro,
		Packages:   packages,
	}
	return json.MarshalIndent(info, "", "  ")
}
