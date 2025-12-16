package platform

import (
	"log"
	"os/exec"
	"strings"
	"github.com/mdgspace/sysreplicate/internal/domain"
)

// FetchPackages returns a list of installed packages for the given base distro.
func FetchPackages(baseDistro string) map[string][]string {
    cmds := make(map[string]*exec.Cmd)
    
    switch baseDistro {
    case "debian":
        cmds["official_packages"] = exec.Command("sh", "-c", domain.DebianFetchCmd)

    case "arch":
        cmds["official_packages"] = exec.Command("sh", "-c", domain.ArchOfficialFetchCmd)
        cmds["yay_packages"] = exec.Command("sh", "-c", domain.ArchYayFetchCmd)

    case "rhel", "fedora":
        cmds["official_packages"] = exec.Command("sh", "-c", domain.RhelFetchCmd)

    case "void":
        cmds["official_packages"] = exec.Command("sh", "-c", domain.VoidFetchCmd)
    default:
        log.Println("Unsupported distro for native package manager detection")
    }

	addOptionalPackageManagers(cmds)

    packageMap := make(map[string][]string)

    for key, value := range cmds {
        output, err := value.CombinedOutput()
        if err != nil {
            log.Println("Error retrieving ", key, ": ", err)
            continue
        }
        packageMap[key] = strings.Split(strings.TrimSpace(string(output)), "\n")
    }
    return packageMap
}

func addOptionalPackageManagers(cmds map[string]*exec.Cmd) {
    _, err := exec.LookPath("flatpak")
    if err == nil {
        cmds["flatpak_packages"] = exec.Command("sh", "-c", domain.FlatpakFetchCmd)
    } else {
        cmds["flatpak_packages"] = exec.Command("true")
    }

    _, err = exec.LookPath("snap")
    if err == nil {
        cmds["snap_packages"] = exec.Command("sh", "-c", domain.SnapFetchCmd)
    } else {
        cmds["snap_packages"] = exec.Command("true")
    }
}