package backup

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mdgspace/sysreplicate/internal/domain"
	"github.com/mdgspace/sysreplicate/system/automation"
	"github.com/mdgspace/sysreplicate/system/output"
	"github.com/mdgspace/sysreplicate/internal/platform"
)

// all backup information in one structure
type UnifiedBackupData struct {
	Timestamp     time.Time                   `json:"timestamp"`
	SystemInfo    output.SystemInfo           `json:"system_info"`
	EncryptedKeys map[string]output.EncryptedKey `json:"encrypted_keys"`
	Dotfiles      []output.Dotfile            `json:"dotfiles"`
	Packages      map[string][]string         `json:"packages"`
	Automation    *automation.AutomationData `json:"automation"`
	EncryptionKey []byte                      `json:"encryption_key"`
	Distro        string                      `json:"distro"`
	BaseDistro    string                      `json:"base_distro"`
}

// complete system backup
type UnifiedBackupManager struct {
	config *EncryptionConfig
}

func NewUnifiedBackupManager() *UnifiedBackupManager {
	return &UnifiedBackupManager{}
}

// complete system backup including keys, dotfiles, and packages
func (ubm *UnifiedBackupManager) CreateUnifiedBackup(customPaths []string) error {
	fmt.Println("Starting unified system backup...")

	// Generate encryption key
	key, err := GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	ubm.config = &EncryptionConfig{
		Key: key,
	}

	// Get system information
	hostname, _ := os.Hostname()
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}

	// Detect distro and get packages
	distro, baseDistro := platform.DetectDistro()
	packages := platform.FetchPackages(baseDistro)
	
	fmt.Printf("Detected distro: %s (%s)\n", distro, baseDistro)
	totalPackages := 0
	for repo, pkgs := range packages {
		if len(pkgs) > 0 {
			fmt.Printf("  %s packages: %d\n", repo, len(pkgs))
			totalPackages += len(pkgs)
		}
	}
	fmt.Printf("  Total packages to restore: %d\n", totalPackages)
	fmt.Println()

	// Create unified backup data
	backupData := &UnifiedBackupData{
		Timestamp:     time.Now(),
		SystemInfo: output.SystemInfo{
			Hostname: hostname,
			Username: username,
			OS:       "linux",
		},
		EncryptedKeys: make(map[string]output.EncryptedKey),
		EncryptionKey: key,
		Packages:      packages,
		Distro:        distro,
		BaseDistro:    baseDistro,
	}

	// 1. Backup SSH/GPG keys
	fmt.Println("Backing up SSH/GPG keys...")
	err = ubm.backupKeys(customPaths, backupData)
	if err != nil {
		fmt.Printf("Warning: Key backup failed: %v\n", err)
	}
	fmt.Println() 

	// 2. Backup dotfiles
	fmt.Println("Backing up dotfiles...")
	err = ubm.backupDotfiles(backupData)
	if err != nil {
		fmt.Printf("Warning: Dotfile backup failed: %v\n", err)
	}
	fmt.Println() 

	// 3. Backup automation files
	fmt.Println("Backing up automation files...")
	err = ubm.backupAutomation(backupData)
	if err != nil {
		fmt.Printf("Warning: Automation backup failed: %v\n", err)
	}
	fmt.Println() 

	// 4. Create unified tarball
	fmt.Println("Creating unified backup tarball...")

	var tarballPath = fmt.Sprintf(domain.UnifiedTarballBasePath,
	time.Now().Format("2006-01-02-15-04-05"))
	
	err = ubm.createUnifiedTarball(backupData, tarballPath)
	if err != nil {
		return fmt.Errorf("failed to create unified tarball: %w", err)
	}

	fmt.Printf("Unified backup completed successfully: %s\n", tarballPath)
	fmt.Println()
	fmt.Printf("Backup Summary:\n")
	fmt.Printf("  Keys: %d files\n", len(backupData.EncryptedKeys))
	fmt.Printf("  Dotfiles: %d files\n", len(backupData.Dotfiles))
	fmt.Printf("  Packages: %d categories\n", len(backupData.Packages))
	
	if backupData.Automation != nil {
		automationCount := len(backupData.Automation.SystemDServices) + len(backupData.Automation.SystemDTimers) + 
			len(backupData.Automation.UserCronjobs) + len(backupData.Automation.SystemCronjobs)
		fmt.Printf("  Automation: %d files (%d services, %d timers, %d user cronjobs, %d system cronjobs)\n",
			automationCount, len(backupData.Automation.SystemDServices), len(backupData.Automation.SystemDTimers),
			len(backupData.Automation.UserCronjobs), len(backupData.Automation.SystemCronjobs))
	}
	
	return nil
}

// SSH/GPG key backup
func (ubm *UnifiedBackupManager) backupKeys(customPaths []string, backupData *UnifiedBackupData) error {
	// Search standard locations
	standardLocations, err := searchStandardLocations()
	if err != nil {
		return fmt.Errorf("failed to search standard locations: %w", err)
	}

	// process the custom paths user might have given while backup
	bm := &BackupManager{}
	customLocations := bm.processCustomPaths(customPaths)

	// Combine all locations
	allLocations := append(standardLocations, customLocations...)
	
	// encrypt and store keys
	keyCount := 0
	for _, location := range allLocations {
		if len(location.Files) > 0 {
			fmt.Printf("  %s keys found:\n", location.Type)
			for _, filePath := range location.Files {
				fileInfo, err := os.Stat(filePath)
				if err != nil {
					continue
				}

				encryptedData, err := EncryptFile(filePath, ubm.config)
				if err != nil {
					fmt.Printf("    Warning: Failed to encrypt %s: %v\n", filePath, err)
					continue
				}

				fmt.Printf("    - %s\n", filePath)
				keyID := filepath.Base(filePath) + "_" + strings.ReplaceAll(filePath, "/", "_")
				backupData.EncryptedKeys[keyID] = output.EncryptedKey{
					OriginalPath:  filePath,
					KeyType:       location.Type,
					EncryptedData: encryptedData,
					Permissions:   uint32(fileInfo.Mode()),
				}
				keyCount++
			}
		}
	}
	
	if keyCount == 0 {
		fmt.Println("  No SSH/GPG keys found")
	} else {
		fmt.Printf("  Total keys backed up: %d\n", keyCount)
	}
	
	return nil
}

// dotfile backup logic
func (ubm *UnifiedBackupManager) backupDotfiles(backupData *UnifiedBackupData) error {
	files, err := ScanDotfiles()
	if err != nil {
		return fmt.Errorf("error scanning dotfiles: %w", err)
	}

	// Convert to output format and show details
	outputFiles := make([]output.Dotfile, len(files))
	dotfileCount := 0
	
	for i, file := range files {
		if !file.IsDir && !file.IsBinary {
			fmt.Printf("  - %s\n", file.Path)
			dotfileCount++
		}
		
		outputFiles[i] = output.Dotfile{
			Path:     file.Path,
			RealPath: file.RealPath,
			IsDir:    file.IsDir,
			IsBinary: file.IsBinary,
			Mode:     file.Mode,
			Content:  file.Content,
		}
	}

	if dotfileCount == 0 {
		fmt.Println("  No dotfiles found")
	} else {
		fmt.Printf("  Total dotfiles backed up: %d\n", dotfileCount)
	}

	backupData.Dotfiles = outputFiles
	return nil
}

func (ubm *UnifiedBackupManager) backupAutomation(backupData *UnifiedBackupData) error {
	am := automation.NewAutomationManager()
	
	data, err := am.DetectAutomation()
	if err != nil {
		return fmt.Errorf("failed to detect automation: %w", err)
	}
	
	if err := am.ValidateAutomationData(data); err != nil {
		return fmt.Errorf("invalid automation data: %w", err)
	}
	
	backupData.Automation = data
	return nil
}

// creating one single tarball containing all backup data
func (ubm *UnifiedBackupManager) createUnifiedTarball(backupData *UnifiedBackupData, tarballPath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(tarballPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to create tarball: %w", err)
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	//// Add main backup metadata
	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %w", err)
	}

	header := &tar.Header{
		Name: "unified_backup.json",
		Mode: 0644,
		Size: int64(len(jsonData)),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if _, err := tarWriter.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write backup data: %w", err)
	}

	// Add dotfiles as separate entries (non-binary only)
	for _, dotfile := range backupData.Dotfiles {
		if dotfile.IsDir || dotfile.IsBinary {
			continue
		}

		file, err := os.Open(dotfile.Path)
		if err != nil {
			fmt.Printf("Warning: Could not open dotfile %s: %v\n", dotfile.Path, err)
			continue
		}

		info, err := file.Stat()
		if err != nil {
			file.Close()
			continue
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			file.Close()
			continue
		}

		hdr.Name = "dotfiles/" + dotfile.RealPath
		
		if err := tarWriter.WriteHeader(hdr); err != nil {
			file.Close()
			continue
		}

		_, err = io.Copy(tarWriter, file)
		file.Close()
		
		if err != nil {
			fmt.Printf("Warning: Failed to add dotfile %s to tarball: %v\n", dotfile.Path, err)
		}
	}

	if backupData.Automation != nil {
		am := automation.NewAutomationManager()
		err := am.BackupAutomation(backupData.Automation, tarWriter)
		if err != nil {
			fmt.Printf("Warning: Failed to add automation files to tarball: %v\n", err)
		}
	}

	return nil
}
