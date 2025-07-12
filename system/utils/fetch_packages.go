package utils

import (
	"log"
	"os/exec"
	"strings"
)

// FetchPackages returns a list of installed packages for the given base distro.
func FetchPackages(baseDistro string) []string {
	var cmd *exec.Cmd
	var cmdYay *exec.Cmd
	switch baseDistro {
	case "debian":
		cmd = exec.Command("dpkg", "--get-selections")
	case "arch":
		cmd = exec.Command("pacman", "-Qn")
		cmdYay = exec.Command("pacman", "-Qm")
	case "rhel", "fedora":
		cmd = exec.Command("rpm", "-qa")
	case "void":
		cmd = exec.Command("xbps-query", "-l")
	default:
		log.Println("Your distro is unsupported, cannot identify package manager!")
		return []string{"unknown"}
	}

	if baseDistro != "arch" {
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("Error in retrieving packages:", err)
		}
		return strings.Split(strings.TrimSpace(string(output)), "\n")
	}

	outputPacman, err := cmd.CombinedOutput()
	outPutYay, errYay := cmdYay.CombinedOutput()
	if err != nil {
		log.Println("Error in retrieving Pacman packages:", err)
	}
	if errYay != nil {
		log.Println("Error in retrieving Yay packages:", errYay)
	}
	pacmanPackages := strings.Split(strings.TrimSpace(string(outputPacman)), "\n")
	yayPackages := strings.Split(strings.TrimSpace(string(outPutYay)), "\n")
	// Mark the split between official and AUR packages
	yayPackages = append([]string{"YayPackages"}, yayPackages...)
	return append(pacmanPackages, yayPackages...)
}