package system

import (
	"fmt"
	"log"

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

	//Output to dist directory
	outputPath := "dist/dotfile-backup.tar.gz"

	// Run the backup
	err := manager.CreateDotfileBackup(outputPath)
	if err != nil {
		fmt.Printf("Backup failed: %v\n", err)
		return
	}

	fmt.Println("Backup complete!")
}
