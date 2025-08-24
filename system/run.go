package system

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/mdgspace/sysreplicate/system/output"
	"github.com/mdgspace/sysreplicate/system/utils"
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
// MUST BE CHANGED IN THE FUTURE
func showMenu() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=== SysReplicate - Distro Hopping Tool ===")
		fmt.Println("1. Generate package replication files")
		fmt.Println("2. Backup SSH/GPG keys")
		fmt.Println("3. Backup dotfiles")
		fmt.Println("4. Exit")
		fmt.Print("Choose an option (1-4): ")

		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			runPackageReplication()
		case "2":
			RunBackup()
		case "3":
			RunDotfileBackup()
		case "4":
			fmt.Println() //exit
			return
		default:
			fmt.Println("Invalid choice. Please select 1, 2, or 3.")
		}
	}
}

// this handles the original package replication functionality
func runPackageReplication() {
	distro, baseDistro := utils.DetectDistro()
	if distro == "unknown" && baseDistro == "unknown" {
		log.Println("Failed to fetch the details of your distro")
		return
	}

	fmt.Println("Distribution:", distro)
	fmt.Println("Built On:", baseDistro)

	packages := utils.FetchPackages(baseDistro)
	jsonObj, err := output.BuildSystemJSON("linux", distro, baseDistro, packages)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}

	if err := os.MkdirAll(outputSysDir, 0744); err != nil {
		log.Println("Error creating sys output directory:", err)
		return
	}

	if err := os.WriteFile(jsonOutputPath, jsonObj, 0644); err != nil {
		log.Println("Error writing JSON output:", err)
		return
	}

	if err := os.MkdirAll(outputScriptsDir, 0744); err != nil {
		log.Println("Error creating scripts output directory:", err)
		return
	}

	if err := output.GenerateInstallScript(baseDistro, packages, scriptOutputPath); err != nil {
		log.Println("Error generating install script:", err)
	} else {
		fmt.Println("Script generated successfully at:", scriptOutputPath)
	}
}
