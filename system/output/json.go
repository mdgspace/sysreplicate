package output

import (
	"encoding/json"
)

// BuildSystemJSON creates a well-structured JSON object for the system info and packages.
func BuildSystemJSON(osType, distro, baseDistro string, packages []string) ([]byte, error) {
	type ArchPackages struct {
		Official []string `json:"official_packages"`
		AUR      []string `json:"aur_packages"`
	}
	type SystemInfo struct {
		OS         string      `json:"os"`
		Distro     string      `json:"distro"`
		BaseDistro string      `json:"base_distro"`
		Packages   interface{} `json:"packages"`
	}
	if baseDistro == "arch" {
		official, aur := SplitArchPackages(packages)
		archPkgs := ArchPackages{Official: official, AUR: aur}
		info := SystemInfo{
			OS:         osType,
			Distro:     distro,
			BaseDistro: baseDistro,
			Packages:   archPkgs,
		}
		return json.MarshalIndent(info, "", "  ")
	}
	info := SystemInfo{
		OS:         osType,
		Distro:     distro,
		BaseDistro: baseDistro,
		Packages:   packages,
	}
	return json.MarshalIndent(info, "", "  ")
}