package backup

import (
    "archive/tar"
    "bufio"
    "compress/gzip"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "syscall"
    "time"

    "golang.org/x/term"
)

//structure of backed up keys
type BackupData struct {
    
	Timestamp     time.Time                `json:"timestamp"`
    
	SystemInfo    SystemInfo              `json:"system_info"`
    
	EncryptedKeys map[string]EncryptedKey `json:"encrypted_keys"`
    
	Salt          []byte                  `json:"salt"`
}

// basic system information
type SystemInfo struct {
    
	Hostname string `json:"hostname"`
    
	Username string `json:"username"`
    
	OS       string `json:"os"`
}

//encrypted key file
type EncryptedKey struct {
    
	OriginalPath string `json:"original_path"`
    
	KeyType      string `json:"key_type"`
    
	EncryptedData string `json:"encrypted_data"`
    
	Permissions  uint32 `json:"permissions"`
}

//handles the backup and restore operations
type BackupManager struct {
    config *EncryptionConfig
}
func NewBackupManager() *BackupManager {
    return &BackupManager{}
}

// create a complete backup of keys
func (bm *BackupManager) CreateBackup(customPaths []string) error {
    fmt.Println("Starting key backup process...")

    //password for encryption
    password, err := bm.getPassword()
    if err != nil {
        return fmt.Errorf("failed to get password: %w", err)
    }

    // generate salt
    salt, err := GenerateSalt()
    if err != nil {
        return fmt.Errorf("failed to generate salt: %w", err)
    }

    bm.config = &EncryptionConfig{
        Password: password,
        Salt:     salt,
    }

    // search standard locations
    fmt.Println("searching standard key locations...")
    standardLocations, err := searchStandardLocations()
    if err != nil {
        return fmt.Errorf("failed to search standard locations: %w", err)
    }

    //add custom paths
    customLocations := bm.processCustomPaths(customPaths)

    //combine all locations
    allLocations := append(standardLocations, customLocations...)

    if len(allLocations) == 0 {
        fmt.Println("No key locations found to backup.")
        return nil
    }

    //create backup data
    backupData := &BackupData{
        Timestamp:     time.Now(),
        SystemInfo:    bm.getSystemInfo(),
        EncryptedKeys: make(map[string]EncryptedKey),
        Salt:          salt,
    }

    //encrypt and store keys
    fmt.Println("Encrypting keys...")
    for _, location := range allLocations {
        err := bm.processLocation(location, backupData)
        if err != nil {
            fmt.Printf("Warning: Failed to process location %s: %v\n", location.Path, err)
            continue
        }
    }

    //creating tarball for the backup storing
    fmt.Println("Creating backup tarball...")
    tarballPath := fmt.Sprintf("dist/key-backup-%s.tar.gz", 
        time.Now().Format("2006-01-02-15-04-05"))
    
    err = bm.createTarball(backupData, tarballPath)
    if err != nil {
        return fmt.Errorf("failed to create tarball: %w", err)
    }

    fmt.Printf("Backup completed successfully: %s\n", tarballPath)
    fmt.Printf("Backed up %d key files\n", len(backupData.EncryptedKeys))
    
    return nil
}

// processLocation processes a single key location
func (bm *BackupManager) processLocation(location KeyLocation, backupData *BackupData) error {
    for _, filePath := range location.Files {
        //get file info for permissions
        fileInfo, err := os.Stat(filePath)
        if err != nil {
            
			continue
        
		}

        // call encryption of the file
        encryptedData, err := EncryptFile(filePath, bm.config)
        if err != nil {
        
			return fmt.Errorf("failed to encrypt %s: %w", filePath, err)
        
		}

        // store encrypted key
        keyID := filepath.Base(filePath) + "_" + strings.ReplaceAll(filePath, "/", "_")
        backupData.EncryptedKeys[keyID] = EncryptedKey{
            OriginalPath:  filePath,
            KeyType:       location.Type,
            EncryptedData: encryptedData,
            Permissions:   uint32(fileInfo.Mode()),
        }
    }
    return nil
}

// processCustomPaths converts custom paths to KeyLocation objects
func (bm *BackupManager) processCustomPaths(customPaths []string) []KeyLocation {
    var locations []KeyLocation
    
    for _, path := range customPaths {
        if path == "" {
        
			continue
        
		}

        // Expand home directory
        if strings.HasPrefix(path, "~/") {
        
			homeDir, _ := os.UserHomeDir()
			path = filepath.Join(homeDir, path[2:])
		}

        fileInfo, err := os.Stat(path)
        if err != nil {
            fmt.Printf("Warning: Custom path %s does not exist\n", path)
            continue
        }

        if fileInfo.IsDir() {
            // Either Process directory
            files, err := discoverKeyFiles(path)
            if err != nil {
                fmt.Printf("Warning: Failed to scan directory %s: %v\n", path, err)
                continue
            }
            
            if len(files) > 0 {
                locations = append(locations, KeyLocation{
                    Path:        path,
                    Type:        "custom",
                    Files:       files,
                    IsDirectory: true,
                })
            }
        } else {
            // Or Process single file
            locations = append(locations, KeyLocation{
                Path:        path,
                Type:        "custom",
                Files:       []string{path},
                IsDirectory: false,
            })
        }
    }
    
    return locations
}

// collect basic system information
func (bm *BackupManager) getSystemInfo() SystemInfo {
    hostname, _ := os.Hostname()
    username := os.Getenv("USER")
    if username == "" {
        username = os.Getenv("USERNAME")
    }
    
    return SystemInfo{
        Hostname: hostname,
        Username: username,
        OS:       "linux",
    }
}

// prompt users for encryption password
func (bm *BackupManager) getPassword() (string, error) {
    fmt.Print("Enter password for key encryption: ")
    bytePassword, err := term.ReadPassword(int(syscall.Stdin))
    if err != nil {
        return "", err
    }
    fmt.Println()
    
    password := string(bytePassword)
    if len(password) < 8 {
        return "", fmt.Errorf("password must be at least 8 characters long") ////just for better recurity - can add more such conditions
    }
    
    return password, nil
}

//compressed tarball with the backup data
func (bm *BackupManager) createTarball(backupData *BackupData, tarballPath string) error {
    // if err := os.MkdirAll(filepath.Dir(tarballPath), 0755); err != nil {
    //     return err
    // }

    // Create tarball file
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

    // Convert backup data to JSON
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

// custom key path prompt to the userss
func GetCustomPaths() []string {
    var paths []string
    scanner := bufio.NewScanner(os.Stdin)
    
    fmt.Println("\nEnter additional key locations (one per line, empty line to finish):")
    fmt.Println("Examples: ~/mykeys/, /opt/certificates/, ~/.config/app/keys")
    
    for {
        fmt.Print("Path: ")
        if !scanner.Scan() {
            break
        }
        
        path := strings.TrimSpace(scanner.Text())
        if path == "" {
            break
        }
        
        paths = append(paths, path)
    }
    
    return paths
}
