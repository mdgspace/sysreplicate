package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/mdgspace/sysreplicate/internal/core/generator"
	"github.com/mdgspace/sysreplicate/internal/domain"
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

	// Convert []Dotfile to []domain.Dotfile
    outputFiles := make([]domain.Dotfile, len(files))
    for i, file := range files {
        outputFiles[i] = domain.Dotfile{
            Path:     file.Path,
            RealPath: file.RealPath,
            IsDir:    file.IsDir,
            IsBinary: file.IsBinary,
            Mode:     file.Mode,
            Content:  file.Content,
        }
    }


	// Create backup metadata
	// struct from output
	meta := &domain.BackupMetadata{
		Timestamp: time.Now(),
		Hostname:  hostname,
		Files:     outputFiles,
	}

    if err := generator.CreateDotfilesBackupTarball(meta, outputTar); err != nil {
        return fmt.Errorf("failed to create backup tarball: %w", err)
    }

	fmt.Println("Backup complete:", outputTar)
	return nil
}
