package ui

import (
	"runtime"
	"strings"

	"github.com/mdgspace/sysreplicate/internal/tui"
)

var string_list = []string{
	"Create Complete System Backup (Recommended)",
	"Restore System from Backup",
	"Generate package replication files only",
	"Backup SSH/GPG keys only",
	"Backup dotfiles only",
}

// Run is the entry point for the system orchestrator.
func Run() {
	osType := runtime.GOOS
	tui.PrintTitle([]string{"Detected OS Type: " + strings.ToUpper(osType)})

	// Map of OS types to their respective actions. We only support Linux.
	osActions := map[string]func(){
        "linux": showMenu,
        "darwin": func() {
            tui.PrintError([]string{"MacOS is not supported"})
        },
        "windows": func() {
            tui.PrintError([]string{"Windows is not supported"})
        },
    }

	// Execute the action based on detected OS type
	if action, exists := osActions[osType]; exists {
		action()
	} else {
		tui.PrintError([]string{"OS not supported"})
	}
}

// showMenu displays the main menu for Linux users
func showMenu() {
	// Display the menu options to the user
	choice := tui.PublicOptions(string_list)

	// Map of user choices to their respective functions
	choiceMap := map[int]func(){
		1: RunUnifiedBackup,
		2: RunRestore,
		3: RunPackageReplication,
		4: RunKeysBackup,
		5: RunDotfileBackup,
	}

	// Execute the selected action based on user choice
	if action, exists := choiceMap[choice]; exists {
		action()
	} else {
		tui.PrintError([]string{"Invalid choice. Please select 1-6."})
	}
}
