package backup

import (
	"archive/tar"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdgspace/sysreplicate/internal/core/generator"
	"github.com/mdgspace/sysreplicate/internal/domain"
)

type RestoreManager struct {
	backupData *UnifiedBackupData
}

func NewRestoreManager() *RestoreManager {
	return &RestoreManager{}
}

// restoring the complete system from a unified backup
func (rm *RestoreManager) RestoreFromBackup(tarballPath string, passphrase string) error {
	fmt.Printf("Starting system restore from: %s\n", tarballPath)

	err := rm.extractBackupData(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to extract backup data: %w", err)
	}

	fmt.Printf("Backup created on: %s\n", rm.backupData.Timestamp.Format("2006-01-09 15:04:05"))
	fmt.Printf("Original system: %s@%s (%s)\n",
		rm.backupData.SystemInfo.Username,
		rm.backupData.SystemInfo.Hostname,
		rm.backupData.Distro)

	fmt.Println("Restoring SSH/GPG keys...")
	err = rm.restoreKeys(passphrase)
	if err != nil {
		fmt.Printf("Warning: Key restoration failed: %v\n", err)
	}
	fmt.Println() 

	// 2. Restore dotfiles
	fmt.Println("Restoring dotfiles...")
	err = rm.restoreDotfiles(tarballPath)
	if err != nil {
		fmt.Printf("Warning: Dotfile restoration failed: %v\n", err)
	}
	fmt.Println() 

	// 3. Generate package installation script
	fmt.Println("Generating package installation script...")
	err = rm.generateInstallScript()
	if err != nil {
		fmt.Printf("Warning: Package script generation failed: %v\n", err)
	}
	fmt.Println()

	fmt.Println("System restore completed successfully!")
	fmt.Println()
	fmt.Printf("Restore Summary:\n")
	fmt.Printf("  Keys restored: %d\n", len(rm.backupData.EncryptedKeys))
	fmt.Printf("  Dotfiles restored: %d\n", len(rm.backupData.Dotfiles))
	fmt.Printf("  Package categories: %d\n", len(rm.backupData.Packages))
	
	if rm.backupData.Automation != nil {
		automationCount := len(rm.backupData.Automation.SystemDServices) + len(rm.backupData.Automation.SystemDTimers) + 
			len(rm.backupData.Automation.UserCronjobs) + len(rm.backupData.Automation.SystemCronjobs)
		if automationCount > 0 {
			fmt.Printf("  Automation files: %d (%d services, %d timers, %d user cronjobs, %d system cronjobs)\n",
				automationCount, len(rm.backupData.Automation.SystemDServices), len(rm.backupData.Automation.SystemDTimers),
				len(rm.backupData.Automation.UserCronjobs), len(rm.backupData.Automation.SystemCronjobs))
		}
	}
	
	fmt.Println()
	fmt.Println("Note: Run the generated install script to restore packages and automation:")
	fmt.Println("  chmod +x dist/restored_packages_install.sh")
	fmt.Println("  ./dist/restored_packages_install.sh")
	
	return nil
}

// extractBackupData extracts and parses the main backup JSON from tarball
func (rm *RestoreManager) extractBackupData(tarballPath string) error {
	file, err := os.Open(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to open tarball: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	var (
		jsonData []byte
		hashStr  string
	)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		switch header.Name {
		case "integrity.hash":
			data, err := io.ReadAll(tarReader)
			if err != nil {
				return fmt.Errorf("failed to read integrity.hash: %w", err)
			}
			hashStr = string(data)
		case "unified_backup.json":
			data, err := io.ReadAll(tarReader)
			if err != nil {
				return fmt.Errorf("failed to read backup data: %w", err)
			}
			jsonData = data
		}
	}

	if jsonData == nil {
		return fmt.Errorf("backup data not found in tarball")
	}

	if hashStr != "" && !VerifyHash(jsonData, hashStr) {
		return fmt.Errorf("integrity check failed: backup data may be corrupted or tampered")
	}

	rm.backupData = &UnifiedBackupData{}
	if err := json.Unmarshal(jsonData, rm.backupData); err != nil {
		return fmt.Errorf("failed to parse backup data: %w", err)
	}

	return nil
}

// decryptiug and restoring SSH/GPG keys to their original locations
func (rm *RestoreManager) restoreKeys(passphrase string) error {
	if rm.backupData.KeyDerivationParams == nil {
		return fmt.Errorf("no key derivation parameters found in backup")
	}
	derivedKey := DeriveKey(passphrase, rm.backupData.KeyDerivationParams)
	config := &EncryptionConfig{
		Key: derivedKey,
	}

	restoredCount := 0
	for keyID, encKey := range rm.backupData.EncryptedKeys {
		fmt.Printf("Restoring key: %s -> %s\n", keyID, encKey.OriginalPath)

		// Decrypt the key data
		decryptedData, err := rm.decryptData(encKey.EncryptedData, config)
		if err != nil {
			fmt.Printf("Warning: Failed to decrypt key %s: %v\n", keyID, err)
			continue
		}

		//// Ensure directory exists
		dir := filepath.Dir(encKey.OriginalPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Warning: Failed to create directory %s: %v\n", dir, err)
			continue
		}

		// Write decrypted data to original location
		err = os.WriteFile(encKey.OriginalPath, decryptedData, os.FileMode(encKey.Permissions))
		if err != nil {
			fmt.Printf("Warning: Failed to write key to %s: %v\n", encKey.OriginalPath, err)
			continue
		}

		restoredCount++
	}

	if restoredCount > 0 {
		fmt.Printf("Successfully restored %d keys\n", restoredCount)
	} else {
		fmt.Println("No keys were restored")
	}
	return nil
}

// extract from tarbell
func (rm *RestoreManager) restoreDotfiles(tarballPath string) error {
	file, err := os.Open(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to open tarball: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	
	homeDir, _ := os.UserHomeDir()
	restoredCount := 0

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Process dotfiles
		if strings.HasPrefix(header.Name, "dotfiles/") {
			relativePath := strings.TrimPrefix(header.Name, "dotfiles/")
			targetPath := filepath.Join(homeDir, relativePath)

			fmt.Printf("Restoring dotfile: %s -> %s\n", header.Name, targetPath)

			// Ensure directory exists
			dir := filepath.Dir(targetPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Warning: Failed to create directory %s: %v\n", dir, err)
				continue
			}

			// creates thje target file and copies content to it
			targetFile, err := os.Create(targetPath)
			if err != nil {
				fmt.Printf("Warning: Failed to create file %s: %v\n", targetPath, err)
				continue
			}
			_, err = io.Copy(targetFile, tarReader)
			targetFile.Close()
			
			if err != nil {
				fmt.Printf("Warning: Failed to copy dotfile content: %v\n", err)
				continue
			}

			// Set permissions
			err = os.Chmod(targetPath, header.FileInfo().Mode())
			if err != nil {
				fmt.Printf("Warning: Failed to set permissions for %s: %v\n", targetPath, err)
			}

			restoredCount++
		}
	}

	if restoredCount > 0 {
		fmt.Printf("Successfully restored %d dotfiles\n", restoredCount)
	} else {
		fmt.Println("No dotfiles were restored")
	}
	return nil
}

// generateInstallScript creates a script to reinstall packages
func (rm *RestoreManager) generateInstallScript() error {
	scriptPath := domain.RestoreScriptPath
	
	// dir check
	if err := os.MkdirAll(filepath.Dir(scriptPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return generator.GenerateInstallScript(rm.backupData.BaseDistro, rm.backupData.Packages, rm.backupData.Automation, scriptPath)
}

// decryptData decrypts base64 encoded data using AES-GCM
func (rm *RestoreManager) decryptData(encryptedBase64 string, config *EncryptionConfig) ([]byte, error) {
	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// creating cipher
	block, err := aes.NewCipher(config.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// creating GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and encrypted data
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	encryptedData := ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}
