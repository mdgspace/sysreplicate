package system

import (
	"fmt"
	"log"
	"os"
	"github.com/mdgspace/sysreplicate/system/backup"
)

// handle backup integration
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

	// Output path
	outputPath := "dist/dotfile-backup.tar.gz"

	// Ensure "dist" directory exists
	if err := os.MkdirAll("dist", os.ModePerm); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	// Run the backup
	err := manager.CreateDotfileBackup(outputPath)
	if err != nil {
		fmt.Printf("Backup failed: %v\n", err)
		return
	}

	fmt.Println("Backup complete!")
}
