package output

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"os"
	"time"
)

//backupData structure for tarball creation
type BackupData struct {
	Timestamp     time.Time                `json:"timestamp"`
	SystemInfo    SystemInfo              `json:"system_info"`
	EncryptedKeys map[string]EncryptedKey `json:"encrypted_keys"`
	EncryptionKey []byte                  `json:"encryption_key"`
}

type SystemInfo struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	OS       string `json:"os"`
}

type EncryptedKey struct {
	OriginalPath  string `json:"original_path"`
	KeyType       string `json:"key_type"`
	EncryptedData string `json:"encrypted_data"`
	Permissions   uint32 `json:"permissions"`
}

type Dotfile struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Mode     uint32 `json:"mode"`
	IsBinary bool   `json:"is_binary"`
}
//create a compressed tarball with the backup data
func CreateBackupTarball(backupData *BackupData, tarballPath string) error {
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
