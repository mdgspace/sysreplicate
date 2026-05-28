package automation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdgspace/sysreplicate/internal/domain"
)

type SystemDUnit = domain.SystemDUnit
type Cronjob = domain.Cronjob
type AutomationData = domain.AutomationData

func (am *AutomationManager) detectSystemDUnits() ([]SystemDUnit, []SystemDUnit, error) {
	var services []SystemDUnit
	var timers []SystemDUnit

	if _, err := os.Stat(domain.SystemdDirPath); os.IsNotExist(err) {
		Log.Warn("SystemD directory %s does not exist, skipping SystemD detection", domain.SystemdDirPath)
		return services, timers, nil
	}
	err := filepath.Walk(domain.SystemdDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !am.isCustomSystemDUnit(path) {
			return nil
		}

		ext := filepath.Ext(path)
		unitName := filepath.Base(path)

		content, err := am.readFileContent(path)
		if err != nil {
			Log.Warn("Could not read SystemD unit %s: %v", path, err)
			return nil
		}
		isEnabled, isActive := am.getSystemDUnitStatus(unitName)

		unit := domain.SystemDUnit{
			Name:      unitName,
			Path:      path,
			Content:   content,
			UnitType:  ext[1:], // Remove the dot
			IsEnabled: isEnabled,
			IsActive:  isActive,
		}

		switch ext {
		case ".service":
			services = append(services, unit)
		case ".timer":
			timers = append(timers, unit)
		case ".target":
			services = append(services, unit)
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to scan SystemD directory: %w", err)
	}

	return services, timers, nil
}

// ///detectCronjobs scans for cron job files
func (am *AutomationManager) detectCronjobs() ([]domain.Cronjob, []domain.Cronjob, error) {
	var userCronjobs []domain.Cronjob
	var systemCronjobs []domain.Cronjob

	userCronPath := fmt.Sprintf(domain.UserCronTemplatePath, am.username)
	if content, err := am.readFileContent(userCronPath); err == nil {
		lines := strings.Split(content, "\n")
		var filteredLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				filteredLines = append(filteredLines, line)
			}
		}

		if len(filteredLines) > 0 {
			userCronjobs = append(userCronjobs, domain.Cronjob{
				Path:    userCronPath,
				Content: content,
				Type:    "user",
			})
		}
	}

	systemCronPaths := []string{
		domain.SystemCrontabDefaultPath,
	}

	if cronDDir := domain.CronDDirPath; am.dirExists(cronDDir) {
		if files, err := filepath.Glob(filepath.Join(cronDDir, "*")); err == nil {
			for _, file := range files {
				if info, err := os.Stat(file); err == nil && !info.IsDir() {
					systemCronPaths = append(systemCronPaths, file)
				}
			}
		}
	}

	for _, cronPath := range systemCronPaths {
		if content, err := am.readFileContent(cronPath); err == nil {
			lines := strings.Split(content, "\n")
			var filteredLines []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					filteredLines = append(filteredLines, line)
				}
			}

			if len(filteredLines) > 0 {
				cronType := "system"
				if strings.Contains(cronPath, "/etc/cron.d/") {
					cronType = "cron_d"
				}

				systemCronjobs = append(systemCronjobs, domain.Cronjob{
					Path:    cronPath,
					Content: content,
					Type:    cronType,
				})
			}
		}
	}

	return userCronjobs, systemCronjobs, nil
}

func (am *AutomationManager) dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
