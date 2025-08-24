package output

import (
	"fmt"
	"os"
)

// generateInstallScript creates a shell script to install all packages for the given distro.
// Returns an error if the script cannot be created or written.
func GenerateInstallScript(baseDistro string, packages map[string][]string, scriptPath string) error {
	f, err := os.Create(scriptPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, "#!/bin/bash\nset -e\necho 'Starting package installation...'")
	if err != nil {
		return err
	}

	var officialInstallCmd string
	switch baseDistro {
	case "debian":
		officialInstallCmd = "sudo apt-get install -y"
	case "arch":
		officialInstallCmd = "sudo pacman -S --noconfirm"
	case "rhel", "fedora":
		officialInstallCmd = "sudo dnf install -y"
	case "void":
		officialInstallCmd = "sudo xbps-install -y"
	default:
		_, _ = fmt.Fprintln(f, "echo 'Unsupported distro for script generation.'")
		return nil
	}

	for repo, pkgs := range packages {
		switch repo {
		case "official_packages":
			_, err = fmt.Fprintf(f, "echo 'Installing packages with %s...'\n", officialInstallCmd)
			if err != nil {
				return err
			}
			for _, pkg := range pkgs {
				if pkg == "" {
					continue
				}
				_, err = fmt.Fprintf(f, "%s %s || true\n", officialInstallCmd, pkg)
				if err != nil {
					return err
				}
			}
		case "yay_packages":
			_, err = fmt.Fprintln(f, "if ! command -v yay >/dev/null; then\n  echo 'yay not found, installing yay...'\n  sudo pacman -S --noconfirm yay\nfi")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(f, "echo 'Installing AUR packages with yay...'")
			if err != nil {
				return err
			}
			for _, pkg := range pkgs {
				if pkg == "" {
					continue
				}
				_, err = fmt.Fprintf(f, "yay -S --noconfirm %s || true\n", pkg)
				if err != nil {
					return err
				}
			}

		case "flatpak_packages":
			_, err = fmt.Fprintf(f, "if ! command -v flatpak >/dev/null; then\n  echo 'flatpak not found, installing flatpak...'\n  %s flatpak\nfi\n", officialInstallCmd)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(f, "echo 'Installing Flatpak packages...'")
			if err != nil {
				return err
			}
			for _, pkg := range pkgs {
				if pkg == "" {
					continue
				}
				_, err = fmt.Fprintf(f, "sudo flatpak install --noninteractive %s || true\n", pkg)
				if err != nil {
					return err
				}
			}

		case "snap_packages":
			_, err = fmt.Fprintf(f, "if ! command -v snap >/dev/null; then\n  echo 'snap not found, installing snapd...'\n  %s snapd\nsudo systemctl enable --now snapd.socket\nfi\n", officialInstallCmd)
			// this limits it to systemctl, but need to replace this in future to support non systemd systems
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(f, "echo 'Installing Snap packages...'")
			if err != nil {
				return err
			}
			for _, pkg := range pkgs {
				if pkg == "" {
					continue
				}
				_, err = fmt.Fprintf(f, "sudo snap install %s || true\n", pkg)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil

}
