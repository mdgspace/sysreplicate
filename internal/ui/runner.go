package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/mdgspace/sysreplicate/system/output"
	"github.com/mdgspace/sysreplicate/internal/platform"
	"github.com/mdgspace/sysreplicate/internal/domain"
)

// Run is the entry point for the system orchestrator.
func Run() {
	osType := runtime.GOOS
	fmt.Println("Detected OS Type:", osType)

	switch osType {
	case "darwin":
		fmt.Println("MacOS is not supported")
		return
	case "windows":
		fmt.Println("Windows is not supported")
		return
	case "linux":
		showMenu() ////main menu component
	default:
		fmt.Println("OS not supported")
	}
}

// showMenu displays the main menu for Linux users
func showMenu() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=== SysReplicate - Distro Hopping Tool ===")
		fmt.Println("1. Create Complete System Backup (Recommended)")
		fmt.Println("2. Restore System from Backup")
		fmt.Println("3. Generate package replication files only")
		fmt.Println("4. Backup SSH/GPG keys only")
		fmt.Println("5. Backup dotfiles only")
		fmt.Println("6. Exit")
		fmt.Print("Choose an option (1-6): ")

		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			RunUnifiedBackup()
		case "2":
			RunRestore()
		case "3":
			runPackageReplication()
		case "4":
			RunBackup()
		case "5":
			RunDotfileBackup()
		case "6":
			fmt.Println("Goodbye Captain!")
			return
		default:
			fmt.Println("Invalid choice. Please select 1-6.")
		}
	}
}

// this handles the original package replication functionality
func runPackageReplication() {
	distro, baseDistro := platform.DetectDistro()
	if distro == "unknown" && baseDistro == "unknown" {
		log.Println("Failed to fetch the details of your distro")
		return
	}

	fmt.Println("Distribution:", distro)
	fmt.Println("Built On:", baseDistro)

	packages := platform.FetchPackages(baseDistro)
	jsonObj, err := output.BuildSystemJSON("linux", distro, baseDistro, packages)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}

	if err := os.MkdirAll(domain.OutputSysDirPath, 0744); err != nil {
		log.Println("Error creating sys output directory:", err)
		return
	}

	if err := os.WriteFile(domain.JsonOutputPath, jsonObj, 0644); err != nil {
		log.Println("Error writing JSON output:", err)
		return
	}

	if err := os.MkdirAll(domain.OutputScriptsDirPath, 0744); err != nil {
		log.Println("Error creating scripts output directory:", err)
		return
	}

	if err := output.GenerateInstallScript(baseDistro, packages, nil, domain.ScriptOutputPath); err != nil {
		log.Println("Error generating install script:", err)
	} else {
		fmt.Println("Script generated successfully at:", domain.ScriptOutputPath)
	}
}
