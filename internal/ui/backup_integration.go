package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdgspace/sysreplicate/internal/core/backup"
	"github.com/mdgspace/sysreplicate/internal/core/generator"
	"github.com/mdgspace/sysreplicate/internal/domain"
	"github.com/mdgspace/sysreplicate/internal/platform"
	"github.com/mdgspace/sysreplicate/internal/tui"
)

// handle backup integration
func RunUnifiedBackup() {
	fmt.Println("=== Unified System Backup (Keys + Dotfiles + Packages) ===")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter encryption passphrase: ")
	if !scanner.Scan() {
		fmt.Println("Failed to read passphrase")
		return
	}
	passphrase := strings.TrimSpace(scanner.Text())
	if passphrase == "" {
		fmt.Println("Passphrase cannot be empty")
		return
	}

	fmt.Print("Confirm passphrase: ")
	if !scanner.Scan() {
		fmt.Println("Failed to read confirmation")
		return
	}
	confirm := strings.TrimSpace(scanner.Text())
	if passphrase != confirm {
		fmt.Println("Passphrases do not match")
		return
	}

	if platform.IsImmutableDistro() {
		fmt.Println()
		fmt.Println("WARNING: Immutable distro detected (Silverblue/Kinoite/SteamOS/MicroOS).")
		fmt.Println("  Standard package install commands do not apply. Use rpm-ostree or")
		fmt.Println("  transactional-update manually to layer packages on immutable systems.")
		fmt.Println()
	}

	ubm := backup.NewUnifiedBackupManager()

	fmt.Println("\nOptional: Add custom key locations")
	customPaths := backup.GetCustomPaths()

	err := ubm.CreateUnifiedBackup(customPaths, passphrase)
	if err != nil {
		log.Printf("Unified backup failed: %v", err)
		return
	}

	fmt.Println("Complete system backup completed successfully!")
	fmt.Println()
	fmt.Println("Your backup includes:")
	fmt.Println("- SSH/GPG keys (encrypted with passphrase)")
	fmt.Println("- Dotfiles (.bashrc, .vimrc, .gitconfig, etc.)")
	fmt.Println("- Package lists for reinstallation")
	fmt.Println("- System automation files (SystemD services, timers, cronjobs)")
	fmt.Println()
	fmt.Println("IMPORTANT: Remember your passphrase! You will need it to restore.")
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

	fmt.Print("Enter backup passphrase: ")
	if !scanner.Scan() {
		fmt.Println("Failed to read passphrase")
		return
	}
	passphrase := strings.TrimSpace(scanner.Text())
	if passphrase == "" {
		fmt.Println("Passphrase cannot be empty")
		return
	}

	rm := backup.NewRestoreManager()
	err = rm.RestoreFromBackup(tarballPath, passphrase)
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
func RunKeysBackup() {
	fmt.Println("=== Key Backup Process ===")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter encryption passphrase: ")
	if !scanner.Scan() {
		fmt.Println("Failed to read passphrase")
		return
	}
	passphrase := strings.TrimSpace(scanner.Text())
	if passphrase == "" {
		fmt.Println("Passphrase cannot be empty")
		return
	}

	fmt.Print("Confirm passphrase: ")
	if !scanner.Scan() {
		fmt.Println("Failed to read confirmation")
		return
	}
	if passphrase != strings.TrimSpace(scanner.Text()) {
		fmt.Println("Passphrases do not match")
		return
	}

	backupManager := backup.NewBackupManager()

	customPaths := backup.GetCustomPaths()

	err := backupManager.CreateBackup(customPaths, passphrase)
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

func RunPackageReplication() {
	distro, baseDistro := platform.DetectDistro()
	if distro == "unknown" && baseDistro == "unknown" {
		log.Println("Failed to fetch the details of your distro")
		return
	}

	tui.PrintText([]string{"Distribution:", distro})
	tui.PrintText([]string{"Built On:", baseDistro})

	packages := platform.FetchPackages(baseDistro)
	jsonObj, err := generator.GenerateMetadata(domain.SystemInfo{
		OS:         "linux",
		Distro:     distro,
		BaseDistro: baseDistro,
	}, packages)
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
		tui.PrintText([]string{"Script generated successfully at:", domain.ScriptOutputPath})
	}
}
