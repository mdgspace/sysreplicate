package system

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"strings"

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

func restoreBackup() {
	fmt.Println("Restoring Backup")
	fmt.Println("Enter backup tarball path")

	reader := bufio.NewReader(os.Stdin)
    name, _ := reader.ReadString('\n') // reads until newline
    name = strings.TrimSpace(name)
    

}
