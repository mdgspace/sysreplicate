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
