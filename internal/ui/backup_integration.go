package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdgspace/sysreplicate/internal/domain"
	"github.com/mdgspace/sysreplicate/internal/core/backup"
)

// handle backup integration
func RunUnifiedBackup() {
	fmt.Println("=== Unified System Backup (Keys + Dotfiles + Packages) ===")

	// unified backup manager
	ubm := backup.NewUnifiedBackupManager()

	// gert all the custom key paths from user
	fmt.Println("\nOptional: Add custom key locations")
	customPaths := backup.GetCustomPaths()

	// Create unified backup
	err := ubm.CreateUnifiedBackup(customPaths)
	if err != nil {
		log.Printf("Unified backup failed: %v", err)
		return
	}

	fmt.Println("Complete system backup completed successfully!")
	fmt.Println()
	fmt.Println("Your backup includes:")
	fmt.Println("- SSH/GPG keys (encrypted)")
	fmt.Println("- Dotfiles (.bashrc, .vimrc, .gitconfig, etc.)")
	fmt.Println("- Package lists for reinstallation")
	fmt.Println("- System automation files (SystemD services, timers, cronjobs)")
}

// system restoration from backup
func RunRestore() {
	fmt.Println("=== System Restore from Backup ===")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter backup tarball path: ")

	if !scanner.Scan() {
		fmt.Println("Failed to read input")
		return
	}

	tarballPath := strings.TrimSpace(scanner.Text())
	if tarballPath == "" {
		fmt.Println("No tarball path provided")
		return
	}

	// Normalize path separators (handle both Windows and Unix paths)
	normalizedPath := strings.ReplaceAll(tarballPath, "\\", "/")

	// Check if file exists and is a file (not directory)
	fileInfo, err := os.Stat(normalizedPath)
	if os.IsNotExist(err) {
		fmt.Printf("Backup file does not exist: %s\n", normalizedPath)
		return
	}
	if err != nil {
		fmt.Printf("Error checking backup file: %v\n", err)
		return
	}
	if fileInfo.IsDir() {
		fmt.Printf("Path is a directory, not a file: %s\n", normalizedPath)
		return
	}

	// Use normalized path for restoration
	tarballPath = normalizedPath

	// Confirm restoration
	fmt.Printf("This will restore your system from: %s\n", tarballPath)
	fmt.Print("WARNING: This will overwrite existing files. Continue? (y/N): ")

	if !scanner.Scan() {
		return
	}

	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("Restoration cancelled")
		return
	}

	// creating previously defined restore manager and run restoration
	rm := backup.NewRestoreManager()
	err = rm.RestoreFromBackup(tarballPath)
	if err != nil {
		log.Printf("Restoration failed: %v", err)
		return
	}

	fmt.Println("\nSystem restoration completed!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Run the generated package installation script")
	fmt.Println("2. Restart your shell or run 'source ~/.bashrc' (or ~/.zshrc)")
	fmt.Println("3. Check that your SSH keys work: 'ssh-add -l'")
	fmt.Println("4. Verify automation files were restored correctly")
}

// rest of the options
func RunBackup() {
	fmt.Println("=== Key Backup Process ===")

	//create backup manager
	backupManager := backup.NewBackupManager()

	//get custom paths from user
	customPaths := backup.GetCustomPaths()

	//create backup
	err := backupManager.CreateBackup(customPaths)
	if err != nil {
		log.Printf("Backup failed: %v", err)
		return
	}

	fmt.Println("Key backup completed successfully!")
}

func RunDotfileBackup() {
	fmt.Println("=== SysReplicate: Distro Dotfile Backup ===")

	// Create a backup manager
	manager := backup.NewDotfileBackupManager()

	// Ensure "dist" directory exists
	if err := os.MkdirAll(domain.OutputScriptsDirPath, os.ModePerm); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	// Run the backup
	err := manager.CreateDotfileBackup(domain.DotfileOutputPath)
	if err != nil {
		fmt.Printf("Backup failed: %v\n", err)
		return
	}

	fmt.Println("Backup complete!")
}
