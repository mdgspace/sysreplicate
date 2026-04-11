package ui

import (
	// "bufio"
	// "fmt"
	"log"
	"os"
	"runtime"
	// "strings"

	"github.com/mdgspace/sysreplicate/internal/core/generator"
	"github.com/mdgspace/sysreplicate/internal/domain"
	"github.com/mdgspace/sysreplicate/internal/platform"
	"github.com/mdgspace/sysreplicate/internal/tui"
)

var string_list = []string{"1. Create Complete System Backup (Recommended)",
"2. Restore System from Backup",
"3. Generate package replication files only",
"4. Backup SSH/GPG keys only",
"5. Backup dotfiles only",
}

// Run is the entry point for the system orchestrator.
func Run() {
	osType := runtime.GOOS
	tui.PublicPrint([]string{"Detected OS Type:", osType})

	switch osType {
	case "darwin":
		tui.PublicPrint([]string{"MacOS is not supported"})
		return
	case "windows":
		tui.PublicPrint([]string{"Windows is not supported"})
		return
	case "linux":
		// tui.PublicOptions(string_list)
		showMenu() ////main menu component
	default:
		tui.PublicPrint([]string{"OS not supported"})
	}
}

// showMenu displays the main menu for Linux users
func showMenu() {
		choice := tui.PublicOptions(string_list)
		switch choice {
		case 1:
			RunUnifiedBackup()
		case 2:
			RunRestore()
		case 3:
			runPackageReplication()
		case 4:
			RunBackup()
		case 5:
			RunDotfileBackup()
		case 6:
			tui.PublicPrint([]string{"Goodbye Captain!"})
			return
		default:
			tui.PublicPrint([]string{"Invalid choice. Please select 1-6."})
		}
	
}

// this handles the original package replication functionality
func runPackageReplication() {
	distro, baseDistro := platform.DetectDistro()
	if distro == "unknown" && baseDistro == "unknown" {
		log.Println("Failed to fetch the details of your distro")
		return
	}

	tui.PublicPrint([]string{"Distribution:", distro})
	tui.PublicPrint([]string{"Built On:", baseDistro})

	packages := platform.FetchPackages(baseDistro)
	jsonObj, err := generator.GenerateMetadata("linux", distro, baseDistro, packages)
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

	if err := generator.GenerateInstallScript(baseDistro, packages, nil, domain.ScriptOutputPath); err != nil {
		log.Println("Error generating install script:", err)
	} else {
		tui.PublicPrint([]string{"Script generated successfully at:", domain.ScriptOutputPath})
	}
}
