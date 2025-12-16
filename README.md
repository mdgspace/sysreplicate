# SysReplicate(DistroHopper) - Linux System Backup and Restore Tool

## Features

### Core Functionality
- **Unified System Backup**: Complete system backup including keys, dotfiles, packages, and automation files
- **System Restore**: Full system restoration from backup tarballs
- **Individual Component Backup**: Separate backup options for keys, dotfiles, or packages only
- **Package Replication**: Generate installation scripts for package managers
- **Encryption**: AES-256-GCM encryption for sensitive data (SSH/GPG keys)

### Supported Components
- **SSH Keys**: Automatic detection and encryption of SSH keys from standard locations
- **GPG Keys**: Backup and encryption of GPG keyrings
- **Dotfiles**: Configuration files (.bashrc, .vimrc, .gitconfig, etc.)
- **Package Lists**: Support for multiple package managers (apt, pacman, dnf, xbps)
- **Automation Files**: SystemD services, timers, and cronjobs
- **Custom Key Locations**: User-defined paths for additional key storage

### Package Manager Support
- **Debian/Ubuntu**: apt-get package manager
- **Arch Linux**: pacman and AUR packages (via yay)
- **Fedora/RHEL**: dnf package manager
- **Void Linux**: xbps package manager
- **Flatpak**: Universal package format
- **Snap**: Canonical's package format

## Installation

### Prerequisites
- Go 1.24.3 or later
- Linux operating system (Windows and macOS are not supported)

### Build Instructions
```bash
git clone <repository-url>
cd sysreplicate
go run main.go
go build -o sysreplicate main.go
```

## Usage

### Running the Application
```bash
./sysreplicate
```

The application will display a menu with the following options:

1. **Create Complete System Backup (Recommended)**
2. **Restore System from Backup**
3. **Generate package replication files only**
4. **Backup SSH/GPG keys only**
5. **Backup dotfiles only**
6. **Exit**

### Menu Options Explained

#### 1. Create Complete System Backup
Creates a unified backup containing:
- SSH/GPG keys (encrypted with AES-256-GCM)
- Dotfiles from home directory
- Package lists for all supported package managers
- SystemD services, timers, and cronjobs
- System metadata (hostname, username, distribution)

The backup is stored as a compressed tarball in the `dist/` directory with timestamp naming.

#### 2. Restore System from Backup
Restores a complete system from a previously created backup:
- Prompts for backup tarball path
- Extracts and decrypts SSH/GPG keys to original locations
- Restores dotfiles to home directory
- Generates package installation script
- Provides restoration summary and next steps

#### 3. Generate Package Replication Files
Creates package lists and installation scripts without backing up keys or dotfiles:
- Detects Linux distribution and base distribution
- Fetches installed packages by category
- Generates JSON metadata file
- Creates installation script for package restoration

#### 4. Backup SSH/GPG Keys Only
Creates encrypted backup of SSH and GPG keys:
- Searches standard locations (~/.ssh/, ~/.gnupg/)
- Accepts custom key locations from user input
- Encrypts keys with AES-256-GCM
- Stores backup as compressed tarball

#### 5. Backup Dotfiles Only
Creates backup of configuration files:
- Scans home directory for dotfiles
- Excludes binary files and directories
- Creates metadata with file information
- Generates compressed tarball

## Project Structure

```
sysreplicate/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── internal/
│   ├── domain/
│   │   └── constants.go
│   └── platform/
│       ├── distro.go
│       └── packages.go
├── system/                    # Core system package
│   ├── run.go                 # Main menu and orchestration
│   ├── settings.go            # Configuration constants
│   ├── backup_integration.go  # Backup/restore integration
│   ├── backup/                # Backup functionality
│   │   ├── unified_backup.go  # Complete system backup
│   │   ├── restore.go         # System restoration
│   │   ├── key.go             # Key backup management
│   │   ├── dotfiles_backup.go # Dotfile backup
│   │   ├── encrypt.go         # Encryption utilities
│   │   ├── search.go          # Key location discovery
│   │   └── dotfile_scanner.go # Dotfile scanning
│   ├── automation/            # Automation file handling
│   │   ├── automation.go      # SystemD and cron detection
│   │   ├── backup.go          # Automation backup
│   │   └── detect.go          # Automation detection
│   └── output/                # Output generation
│       ├── json.go            # JSON metadata generation
│       ├── script.go          # Installation script generation
│       └── tarball.go         # Tarball creation
└── dist/                      # Output directory
    ├── unified-backup-*.tar.gz # Complete system backups
    ├── key-backup-*.tar.gz    # Key-only backups
    ├── dotfile-backup.tar.gz  # Dotfile backups
    └── restored_packages_install.sh # Generated install script
```

## Technical Details

### Encryption
- **Algorithm**: AES-256-GCM
- **Key Generation**: Cryptographically secure random 32-byte keys
- **Data Format**: Base64-encoded encrypted data with nonce
- **Scope**: Only SSH/GPG keys are encrypted; dotfiles and packages are stored in plaintext

### Backup Format
- **Container**: Compressed tarball (.tar.gz)
- **Metadata**: JSON file containing system information and file lists
- **Structure**: 
  - `unified_backup.json`: Main metadata and encrypted keys
  - `dotfiles/`: Directory containing dotfile contents
  - `automation/`: Directory containing automation files

### Distribution Detection
The tool detects Linux distributions by reading `/etc/os-release` and extracting:
- `ID`: Specific distribution name
- `ID_LIKE`: Base distribution family

### Package Detection
Package lists are fetched based on the base distribution:
- **Debian-based**: Uses `dpkg` and `apt` commands
- **Arch-based**: Uses `pacman` and `yay` commands
- **Red Hat-based**: Uses `dnf` commands
- **Void**: Uses `xbps` commands

### Automation Detection
The tool detects and backs up:
- **SystemD Services**: Custom services in `/etc/systemd/system/`
- **SystemD Timers**: Custom timers in `/etc/systemd/system/`
- **User Cronjobs**: User-specific cron jobs
- **System Cronjobs**: System-wide cron jobs

## Security Considerations

- SSH and GPG keys are encrypted with AES-256-GCM before storage
- Encryption keys are generated randomly for each backup
- Backup files should be stored securely as they contain sensitive data
- The tool does not require user passwords for encryption (uses random keys)

## Limitations

- Only supports Linux operating systems
- Requires appropriate permissions to read system files
- Package restoration may fail if packages are not available in target distribution repositories
- Some automation files may require manual configuration after restoration

## Output Files

### Backup Files
- `unified-backup-YYYY-MM-DD-HH-MM-SS.tar.gz`: Complete system backup
- `key-backup-YYYY-MM-DD-HH-MM-SS.tar.gz`: SSH/GPG keys only
- `dotfile-backup.tar.gz`: Dotfiles only

### Generated Scripts
- `restored_packages_install.sh`: Package installation script for restoration
- `setup.sh`: Package installation script for replication only

### Metadata Files
- `package.json`: System and package information in JSON format

