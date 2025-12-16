package domain

// import (
// 	"fmt"
// 	"time"
// )

const (
	OutputSysDirPath         = "dist/sys-info"
	OutputScriptsDirPath     = "dist"
	JsonOutputPath           = OutputSysDirPath + "/package.json"
	ScriptOutputPath         = OutputScriptsDirPath + "/setup.sh"
	DotfileOutputPath        = "dist/dotfile-backup.tar.gz"
	OsReleasePath            = "/etc/os-release"
	RestoreScriptPath        = "dist/restored_packages_install.sh"
	SystemdDirPath           = "/etc/systemd/system"
	UserCronPathTemplatePath = "/var/spool/cron/crontabs/%s"
	SystemCrontabDefaultPath = "/etc/crontab"
	CronDDirPath             = "/etc/cron.d"
	UnifiedTarballBasePath   = "dist/unified-backup-%s.tar.gz"

	DebianInstallCmd = "sudo apt-get install -y"
	ArchInstallCmd   = "sudo pacman -S --noconfirm"
	RhelInstallCmd   = "sudo dnf install -y"
	VoidInstallCmd   = "sudo xbps-install -y"

	DebianFetchCmd = `dpkg-query -W -f='${Package}\n' | sort > /tmp/all.txt
apt-mark showmanual | sort > /tmp/manual.txt
comm -12 /tmp/all.txt /tmp/manual.txt | xargs -r dpkg-query -W -f='${Package}=${Version}\n'
rm /tmp/all.txt /tmp/manual.txt`
	ArchOfficialFetchCmd = `pacman -Qen | cut -d' ' -f1`
	ArchYayFetchCmd      = `pacman -Qem | cut -d' ' -f1`
	RhelFetchCmd         = "rpm -qa"
	VoidFetchCmd         = "xbps-query -l"

	FlatpakFetchCmd = "flatpak list --app --columns=origin,application"
	SnapFetchCmd    = "snap list | awk 'NR>1 {print $1}'"
)

var DotfilePaths = []string{
	"~/.bashrc",
	"~/.zshrc",
	"~/.vimrc",
	"~/.config",
	"~/.bash_history",
	"~/.zsh_history",
	"~/.gitconfig",
	"~/.profile",
	"~/.npmrc",
}

var StandardKeyLocations = []string{
	"~/.ssh/",
	"~/.gnupg/",
}

var SshPatterns = []string{
	"id_rsa", "id_dsa", "id_ecdsa", "id_ed25519",
	"authorized_keys", "known_hosts", "config",
}

var GpgPatterns = []string{
	"pubring.gpg", "secring.gpg", "trustdb.gpg",
	"gpg.conf", "gpg-agent.conf",
}

var PackageManagedDirs = []string{
	"/usr/lib/systemd/system/",
	"/lib/systemd/system/",
	"/usr/share/systemd/",
}

