package generator

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"github.com/mdgspace/sysreplicate/internal/domain"
)

// create a compressed tarball with the backup data
func CreateBackupTarball(backupData *domain.BackupData, tarballPath string) error {
	//create tarball file
	file, err := os.Create(tarballPath)
	if err != nil {
		return err
	}
	defer file.Close()

	//gzip writer
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	//tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	//convertto JSON
	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return err
	}

	//add JSON file to tarball
	header := &tar.Header{
		Name: "backup.json",
		Mode: 0644,
		Size: int64(len(jsonData)),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err := tarWriter.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func CreateDotfilesBackupTarball(meta *domain.BackupMetadata, tarballPath string) error {
	// Create the tarball file
	file, err := os.Create(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to create tarball: %w", err)
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	jsonData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Add metadata as backup.json
	header := &tar.Header{
		Name: "backup.json",
		Mode: 0644,
		Size: int64(len(jsonData)),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header for metadata: %w", err)
	}
	if _, err := tarWriter.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write metadata to tar: %w", err)
	}

	// Add dotfiles
	for _, f := range meta.Files {
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
		hdr.Name = f.RealPath
		tarWriter.WriteHeader(hdr)
		io.Copy(tarWriter, file)
	}

	return nil
}
