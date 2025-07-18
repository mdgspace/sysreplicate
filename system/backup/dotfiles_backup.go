package backup

import (
	"fmt"
	"os"
	"time"
	"github.com/mdgspace/sysreplicate/system/output"
)

type BackupMetadata struct {
	Timestamp time.Time `json:"timestamp"`
	Hostname  string    `json:"hostname"`
	Files     []Dotfile `json:"files"`
}

type DotfileBackupManager struct{}

func NewDotfileBackupManager() *DotfileBackupManager {
	return &DotfileBackupManager{}
}

func (db *DotfileBackupManager) CreateDotfileBackup(outputTar string) error {
	// Scan dotfiles
	files, err := ScanDotfiles()
	if err != nil {
		return fmt.Errorf("error scanning dotfiles: %w", err)
	}

	hostname, _ := os.Hostname()

	// Create backup metadata
	meta := &BackupMetadata{
		Timestamp: time.Now(),
		Hostname:  hostname,
		Files:     files,
	}

	// Tarball creation 
	if err := CreateDotfilesBackupTarball(meta, outputTar); err != nil {
		return fmt.Errorf("failed to create backup tarball: %w", err)
	}

	fmt.Println("Backup complete:", outputTar)
	return nil
}
