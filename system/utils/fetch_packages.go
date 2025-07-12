package utils

import (
	"log"
	"os/exec"
	"strings"
)

// FetchPackages returns a list of installed packages for the given base distro.
func FetchPackages(baseDistro string) map[string][]string {
	cmds := make(map[string]*exec.Cmd)

	switch baseDistro {
	case "debian":
		// cmds["official_packages"] = exec.Command("dpkg", "--get-selections")
		cmds["official_packages"] = exec.Command("sh", "-c", `dpkg-query -W -f='${Package}\n' | sort > /tmp/all.txt
apt-mark showmanual | sort > /tmp/manual.txt
comm -12 /tmp/all.txt /tmp/manual.txt | xargs -r dpkg-query -W -f='${Package}=${Version}\n'
rm /tmp/all.txt /tmp/manual.txt
`)
		cmds["flatpak_packages"] = exec.Command("flatpak", "list", "--app", "--columns=origin,application")
		cmds["snap_packages"] = exec.Command("sh", "-c", "snap list | awk 'NR>1 {print $1}'")

	case "arch":
		// cmds["official_packages"] = exec.Command("pacman", "-Qen")
		// cmds["yay_packages"] = exec.Command("pacman", "-Qem")
		cmds["official_packages"] = exec.Command("sh", "-c", `pacman -Qen | cut -d' ' -f1`)
		cmds["yay_packages"] = exec.Command("sh", "-c", `pacman -Qem | cut -d' ' -f1`)
		cmds["flatpak_packages"] = exec.Command("flatpak", "list", "--app", "--columns=origin,application")
		cmds["snap_packages"] = exec.Command("sh", "-c", "snap list | awk 'NR>1 {print $1}'")

	case "rhel", "fedora":
		cmds["official_packages"] = exec.Command("rpm", "-qa") // need to change this later
		cmds["flatpak_packages"] = exec.Command("flatpak", "list", "--app", "--columns=origin,application")
		cmds["snap_packages"] = exec.Command("sh", "-c", "snap list | awk 'NR>1 {print $1}'")

	case "void":
		cmds["official_packages"] = exec.Command("xbps-query", "-l") // need to change this later
		cmds["flatpak_packages"] = exec.Command("flatpak", "list", "--app", "--columns=origin,application")
		cmds["snap_packages"] = exec.Command("sh", "-c", "snap list | awk 'NR>1 {print $1}'")

	default:
		log.Println("Your distro is unsupported, cannot identify package manager!")
		return map[string][]string{
			"error": {"unsupported distro"},
		}
	}

	packageMap := make(map[string][]string)

	for key, value := range cmds {
		output, err := value.CombinedOutput()
		if err != nil {
			log.Println("Error in retrieving ", key, ": ", err)
			continue
		}
		packageMap[key] = strings.Split(strings.TrimSpace((string(output))), "\n")
	}
	return packageMap

}
