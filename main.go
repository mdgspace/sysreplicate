package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	osType := runtime.GOOS
	fmt.Println("Detected OS Type:", osType)

	switch osType {
	case "darwin":
		fmt.Println("MacOS is not supported")
	case "windows":
		fmt.Println("Windows is not supported")
	case "linux":
		distro, base_distro := fetchLinuxDistro()
		if distro == "unknown" && base_distro == "unknown" {
			fmt.Println("Failed to fetch the details of your distro")
		}
		fmt.Println("Distribution:", distro)
		fmt.Println("Built On:", base_distro)
		packages := fetchPackages(base_distro)
		jsonObj, err := buildSystemJSON(osType, distro, base_distro, packages)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
		} else {
			os.MkdirAll("outputs/sys", 0744)
			os.WriteFile("outputs/sys/package_info.json", jsonObj, 0744)
			os.MkdirAll("outputs/scripts", 0744)
			generateInstallScript(base_distro, packages, "outputs/scripts/setup.sh")
		}
	default:
		fmt.Println("OS not supported")
	}
}

func fetchLinuxDistro() (string, string) {
	data, err := os.ReadFile("/etc/os-release")

	if err != nil {
		return "unknown", "unknown"
	}
	var distro, base_distro string

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			distro = strings.Trim(strings.SplitN(line, "=", 2)[1], `"`)
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			base_distro = strings.Trim(strings.SplitN(line, "=", 2)[1], `"`)
		}
	}

	return distro, base_distro
}

func fetchPackages(base_distro string) []string {
	var cmd *exec.Cmd
	var cmdYay *exec.Cmd
	switch base_distro {
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
		fmt.Println("Your distro is unsupported, cannot identify package manager !")
		return []string{"unknown"}
	}

	if base_distro != "arch" {

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error in retrieving packages: ", err)
		}
		fmt.Println("Installed Packages: ")
		packages := strings.Split(string(output), "\n")
		return packages
	}

	outputPacman, err := cmd.CombinedOutput()
	outPutYay, errYay := cmdYay.CombinedOutput()

	if err != nil {
		fmt.Println("Error in retrieving Pacman packages: ", err)
	}
	if errYay != nil {
		fmt.Println("Error in retrieving Yay packages: ", errYay)
	}

	pacmanPackages := strings.Split(string(outputPacman), "\n")
	yayPackages := strings.Split(string(outPutYay), "\n")

	yayPackages = append([]string{"YayPackages"}, yayPackages...)

	return append(pacmanPackages, yayPackages...)
}

func buildSystemJSON(osType, distro, base_distro string, packages []string) ([]byte, error) {
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

	if base_distro == "arch" {
		official := []string{}
		aur := []string{}
		isAUR := false
		for _, pkg := range packages {
			if pkg == "YayPackages" {
				isAUR = true
				continue
			}
			if isAUR {
				aur = append(aur, pkg)
			} else {
				official = append(official, pkg)
			}
		}
		archPkgs := ArchPackages{Official: official, AUR: aur}
		info := SystemInfo{
			OS:         osType,
			Distro:     distro,
			BaseDistro: base_distro,
			Packages:   archPkgs,
		}
		return json.MarshalIndent(info, "", "  ")
	}

	info := SystemInfo{
		OS:         osType,
		Distro:     distro,
		BaseDistro: base_distro,
		Packages:   packages,
	}
	return json.MarshalIndent(info, "", "  ")
}

func generateInstallScript(base_distro string, packages []string, scriptPath string) {
	f, err := os.Create(scriptPath)
	if err != nil {
		fmt.Println("Error creating script:", err)
		return
	}
	defer f.Close()

	f.WriteString("#!/bin/bash\n")
	f.WriteString("set -e\n")
	f.WriteString("echo 'Starting package installation...'\n")

	var installCmd string
	switch base_distro {
	case "debian":
		installCmd = "sudo apt-get install -y"
	case "arch":
		installCmd = "sudo pacman -S --noconfirm"
	case "rhel", "fedora":
		installCmd = "sudo dnf install -y"
	case "void":
		installCmd = "sudo xbps-install -y"
	default:
		f.WriteString("echo 'Unsupported distro for script generation.'\n")
		return
	}

	if base_distro == "arch" {
		official := []string{}
		aur := []string{}
		isAUR := false
		for _, pkg := range packages {
			if pkg == "YayPackages" {
				isAUR = true
				continue
			}
			if isAUR {
				aur = append(aur, pkg)
			} else {
				official = append(official, pkg)
			}
		}
		f.WriteString("echo 'Installing official packages with pacman...'\n")
		for _, pkg := range official {
			if pkg == "" {
				continue
			}
			f.WriteString(fmt.Sprintf("%s %s || true\n", installCmd, pkg))
		}
		f.WriteString("if ! command -v yay >/dev/null; then\n  echo 'yay not found, installing yay...'\n  sudo pacman -S --noconfirm yay\nfi\n")
		f.WriteString("echo 'Installing AUR packages with yay...'\n")
		for _, pkg := range aur {
			if pkg == "" {
				continue
			}
			f.WriteString(fmt.Sprintf("yay -S --noconfirm %s || true\n", pkg))
		}
		fmt.Println("Script generated successfully")
		return
	}

	f.WriteString(fmt.Sprintf("echo 'Installing packages with %s...'\n", installCmd))
	for _, pkg := range packages {
		if pkg == "" {
			continue
		}
		f.WriteString(fmt.Sprintf("%s %s || true\n", installCmd, pkg))
	}
	fmt.Println("Script generated successfully")
}
