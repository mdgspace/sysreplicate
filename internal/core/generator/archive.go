package generator

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mdgspace/sysreplicate/internal/domain"
)

// create a compressed tarball with the backup data
func CreateBackupTarball(backupData *domain.BackupData, tarballPath string) error {
	file, err := os.Create(tarballPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return err
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(jsonData))
	hashEntry := &tar.Header{
		Name: "integrity.hash",
		Mode: 0644,
		Size: int64(len(hash)),
	}
	if err := tarWriter.WriteHeader(hashEntry); err != nil {
		return err
	}
	if _, err := tarWriter.Write([]byte(hash)); err != nil {
		return err
	}

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

		if err := func() error {
			file, err := os.Open(f.Path)
			if err != nil {
				return err
			}
			defer file.Close()

			info, _ := file.Stat()
			hdr, _ := tar.FileInfoHeader(info, "")
			hdr.Name = f.RealPath
			if err := tarWriter.WriteHeader(hdr); err != nil {
				return err
			}
			_, err = io.Copy(tarWriter, file)
			return err
		}(); err != nil {
			continue
		}
	}

	return nil
}
