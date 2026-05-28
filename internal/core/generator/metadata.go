package generator
	
import (
	"encoding/json"
	"github.com/mdgspace/sysreplicate/internal/domain"
)

// GenerateMetadata creates a well-structured JSON object for the system info and packages.
func GenerateMetadata(info domain.SystemInfo, packages map[string][]string) ([]byte, error) {
	output := struct {
		OS         string              `json:"os"`
		Distro     string              `json:"distro"`
		BaseDistro string              `json:"base_distro"`
		Packages   map[string][]string `json:"packages"`
	}{
		OS:         info.OS,
		Distro:     info.Distro,
		BaseDistro: info.BaseDistro,
		Packages:   packages,
	}
	return json.MarshalIndent(output, "", "  ")
}
