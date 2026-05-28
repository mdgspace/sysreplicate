package domain
import(
	"time"
	"os"
	// "github.com/mdgspace/sysreplicate/system/output"
)

type SystemInfo struct {
	Hostname   string `json:"hostname"`
	Username   string `json:"username"`
	OS         string `json:"os"`
	Distro     string `json:"distro"`
	BaseDistro string `json:"base_distro"`
}

type BackupData struct {
	Timestamp     time.Time               `json:"timestamp"`
	SystemInfo    SystemInfo              `json:"system_info"`
	EncryptedKeys map[string]EncryptedKey `json:"encrypted_keys"`
	EncryptionKey []byte                  `json:"encryption_key"`
}

type EncryptedKey struct {
	OriginalPath  string `json:"original_path"`
	KeyType       string `json:"key_type"`
	EncryptedData string `json:"encrypted_data"`
	Permissions   uint32 `json:"permissions"`
}

type Dotfile struct {
	Path     string      `json:"path"`
	RealPath string      `json:"real_path"`
	IsDir    bool        `json:"is_dir"`
	IsBinary bool        `json:"is_binary"`
	Mode     os.FileMode `json:"mode"`
	Content  string      `json:"content"` // ignore for the binary files
}

type BackupMetadata struct {
	Timestamp time.Time `json:"timestamp"`
	Hostname  string    `json:"hostname"`
	Files     []Dotfile `json:"files"`
}