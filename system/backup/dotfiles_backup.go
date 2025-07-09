package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
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
	files, err := ScanDotfiles()
	if err != nil {
		return fmt.Errorf("error scanning dotfiles: %w", err)
	}

	hostname, _ := os.Hostname()

	meta := BackupMetadata{
		Timestamp: time.Now(),
		Hostname:  hostname,
		Files:     files,
	}

	tarFile, err := os.Create(outputTar)
	if err != nil {
		return fmt.Errorf("failed to create tar file: %w", err)
	}
	defer tarFile.Close()

	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Write metadata JSON
	metaBytes, _ := json.MarshalIndent(meta, "", "  ")
	tarWriter.WriteHeader(&tar.Header{
		Name: "backup.json",
		Mode: 0644,
		Size: int64(len(metaBytes)),
	})
	tarWriter.Write(metaBytes)

	// Add dotfiles
	for _, f := range files {
		if f.IsDir {
			continue
		}
		file, err := os.Open(f.Path)
		if err != nil {
			continue
		}
		defer file.Close()

		info, _ := file.Stat()
		hdr, _ := tar.FileInfoHeader(info, "")
		hdr.Name = f.RelPath
		tarWriter.WriteHeader(hdr)
		io.Copy(tarWriter, file)
	}

	fmt.Println("Backup complete:", outputTar)
	return nil
}
