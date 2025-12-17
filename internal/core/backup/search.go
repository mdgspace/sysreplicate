package backup

import (
    "os"
    "path/filepath"
    "strings"

    "github.com/mdgspace/sysreplicate/internal/domain"
)

// any saved keylcoation
type KeyLocation struct {
        Path        string

        Type        string // "ssh", "gpg", "custom" strongs can be stored
        
        Files       []string
        
        IsDirectory bool
}

// searches for keys in standard locations
func searchStandardLocations() ([]KeyLocation, error) {
    var locations []KeyLocation
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    for _, location := range domain.StandardKeyLocations {
        // replace ~operator with actual home directory
        fullPath := strings.Replace(location, "~", homeDir, 1)
        
        if _, err := os.Stat(fullPath); os.IsNotExist(err) {
            continue // skip on dir invalid
        }

        
        keyType := determineKeyType(fullPath)
        files, err := discoverKeyFiles(fullPath)
        if err != nil {
            continue

        }

        if len(files) > 0 {
            locations = append(locations, KeyLocation{
                
                Path:        fullPath,
                Type:        keyType,
                Files:       files,
                IsDirectory: true,
            })
        }
    }

    return locations, nil
}

// determineKeyType identifies the type of keys based on directory path
func determineKeyType(path string) string {
    if strings.Contains(path, ".ssh") {
        return "ssh"
    }
    if strings.Contains(path, ".gnupg") {
        return "gpg"
    }
    return "custom"
}

// discoverKeyFiles finds all key files in a directory
func discoverKeyFiles(dirPath string) ([]string, error) {
    var keyFiles []string
    
    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !info.IsDir() && isKeyFile(path, info) {
            keyFiles = append(keyFiles, path)
        }
        
        return nil
    })
    
    return keyFiles, err
}

// isKeyFile determines if a file is likely a key file
func isKeyFile(path string, info os.FileInfo) bool {
    name := info.Name()
    
    // check SSH 
    for _, pattern := range domain.SshPatterns {
        if strings.Contains(name, pattern) {
            return true
        }
    }
    
    // check GPG
    for _, pattern := range domain.GpgPatterns {
        if strings.Contains(name, pattern) {
            return true
        }
    }
    
    // check for private key extensions - more can be added
    if strings.HasSuffix(name, ".pub") || 
       strings.HasSuffix(name, ".pem") || 
       strings.HasSuffix(name, ".key") {
        return true
    }
    
    return false
}
