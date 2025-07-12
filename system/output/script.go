package output

import (
	"fmt"
	"os"
)

// splitArchPackages splits the combined package list into official and AUR packages for Arch-based distros.
func SplitArchPackages(packages []string) (official, aur []string) {
	isAUR := false
	for _, pkg := range packages {
		if pkg == "YayPackages" {
			isAUR = true
			continue
		}
		if isAUR {
			if pkg != "" {
				aur = append(aur, pkg)
			}
		} else {
			if pkg != "" {
				official = append(official, pkg)
			}
		}
	}
	return
}

// generateInstallScript creates a shell script to install all packages for the given distro.
// Returns an error if the script cannot be created or written.
func GenerateInstallScript(baseDistro string, packages []string, scriptPath string) error {
	f, err := os.Create(scriptPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("#!/bin/bash\nset -e\necho 'Starting package installation...'\n")
	if err != nil {
		return err
	}

	var installCmd string
	switch baseDistro {
	case "debian":
		installCmd = "sudo apt-get install -y"
	case "arch":
		installCmd = "sudo pacman -S --noconfirm"
	case "rhel", "fedora":
		installCmd = "sudo dnf install -y"
	case "void":
		installCmd = "sudo xbps-install -y"
	default:
		_, _ = f.WriteString("echo 'Unsupported distro for script generation.'\n")
		return nil
	}

	if baseDistro == "arch" {
		official, aur := SplitArchPackages(packages)
		_, err = f.WriteString("echo 'Installing official packages with pacman...'\n")
		if err != nil {
			return err
		}
		for _, pkg := range official {
			if pkg == "" {
				continue
			}
			_, err = f.WriteString(fmt.Sprintf("%s %s || true\n", installCmd, pkg))
			if err != nil {
				return err
			}
		}
		_, err = f.WriteString("if ! command -v yay >/dev/null; then\n  echo 'yay not found, installing yay...'\n  sudo pacman -S --noconfirm yay\nfi\n")
		if err != nil {
			return err
		}
		_, err = f.WriteString("echo 'Installing AUR packages with yay...'\n")
		if err != nil {
			return err
		}
		for _, pkg := range aur {
			if pkg == "" {
				continue
			}
			_, err = f.WriteString(fmt.Sprintf("yay -S --noconfirm %s || true\n", pkg))
			if err != nil {
				return err
			}
		}
		return nil
	}

	_, err = f.WriteString(fmt.Sprintf("echo 'Installing packages with %s...'\n", installCmd))
	if err != nil {
		return err
	}
	for _, pkg := range packages {
		if pkg == "" {
			continue
		}
		_, err = f.WriteString(fmt.Sprintf("%s %s || true\n", installCmd, pkg))
		if err != nil {
			return err
		}
	}
	return nil
}
