package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/mdgspace/sysreplicate/internal/core/generator"
	"github.com/mdgspace/sysreplicate/internal/domain"
)

type DotfileBackupManager struct{}

func NewDotfileBackupManager() *DotfileBackupManager {
	return &DotfileBackupManager{}
}

func (db *DotfileBackupManager) CreateDotfileBackup(outputTar string) error {
	files, err := ScanDotfiles()
	if err != nil {
		return fmt.Errorf("error scanning dotfiles: %w", err)
	}

	hostname, _ := os.Hostname()

	meta := &domain.BackupMetadata{
		Timestamp: time.Now(),
		Hostname:  hostname,
		Files:     files,
	}

	if err := generator.CreateDotfilesBackupTarball(meta, outputTar); err != nil {
		return fmt.Errorf("failed to create backup tarball: %w", err)
	}

	fmt.Println("Backup complete:", outputTar)
	return nil
}
