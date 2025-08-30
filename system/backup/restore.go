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

	"github.com/mdgspace/sysreplicate/system/output"
)

type RestoreManager struct {
	backupData *UnifiedBackupData
}

func NewRestoreManager() *RestoreManager {
	return &RestoreManager{}
}

// restoring the complete system from a unified backup
func (rm *RestoreManager) RestoreFromBackup(tarballPath string) error {
	fmt.Printf("Starting system restore from: %s\n", tarballPath)

	// extract and parse backup data
	err := rm.extractBackupData(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to extract backup data: %w", err)
	}

	// backup information and furthur instructions
	///s
	fmt.Printf("Backup created on: %s\n", rm.backupData.Timestamp.Format("2006-01-09 15:04:05"))
	fmt.Printf("Original system: %s@%s (%s)\n", 
		rm.backupData.SystemInfo.Username, 
		rm.backupData.SystemInfo.Hostname,
		rm.backupData.Distro)

	// 1. Restore SSH/GPG keys
	fmt.Println("Restoring SSH/GPG keys...")
	err = rm.restoreKeys()
	if err != nil {
		fmt.Printf("Warning: Key restoration failed: %v\n", err)
	}

	// 2. Restore dotfiles
	fmt.Println("Restoring dotfiles...")
	err = rm.restoreDotfiles(tarballPath)
	if err != nil {
		fmt.Printf("Warning: Dotfile restoration failed: %v\n", err)
	}

	// 3. Generate package installation script
	fmt.Println("Generating package installation script...")
	err = rm.generateInstallScript()
	if err != nil {
		fmt.Printf("Warning: Package script generation failed: %v\n", err)
	}

	fmt.Println("System restore completed successfully!")
	fmt.Println("Note: Run the generated install script to restore packages:")
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

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		if header.Name == "unified_backup.json" {
			data, err := io.ReadAll(tarReader)
			if err != nil {
				return fmt.Errorf("failed to read backup data: %w", err)
			}

			rm.backupData = &UnifiedBackupData{}
			err = json.Unmarshal(data, rm.backupData)
			if err != nil {
				return fmt.Errorf("failed to parse backup data: %w", err)
			}
			
			return nil
		}
	}

	return fmt.Errorf("backup data not found in tarball")
}

// decryptiug and restoring SSH/GPG keys to their original locations
func (rm *RestoreManager) restoreKeys() error {
	config := &EncryptionConfig{
		Key: rm.backupData.EncryptionKey,
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

	fmt.Printf("Successfully restored %d keys\n", restoredCount)
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

	fmt.Printf("Successfully restored %d dotfiles\n", restoredCount)
	return nil
}

// generateInstallScript creates a script to reinstall packages
func (rm *RestoreManager) generateInstallScript() error {
	scriptPath := "dist/restored_packages_install.sh"
	
	// dir check
	if err := os.MkdirAll(filepath.Dir(scriptPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return output.GenerateInstallScript(rm.backupData.BaseDistro, rm.backupData.Packages, scriptPath)
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
