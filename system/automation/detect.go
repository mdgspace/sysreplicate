package automation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)
func (am *AutomationManager) detectSystemDUnits() ([]SystemDUnit, []SystemDUnit, error) {
	var services []SystemDUnit
	var timers []SystemDUnit
	
	systemdDir := "/etc/systemd/system"
	if _, err := os.Stat(systemdDir); os.IsNotExist(err) {
		fmt.Printf("SystemD directory %s does not exist, skipping SystemD detection\n", systemdDir)
		return services, timers, nil
	}
	err := filepath.Walk(systemdDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {return err}
		if info.IsDir() {return nil}
		if !am.isCustomSystemDUnit(path) {return nil}
		
		ext := filepath.Ext(path)
		unitName := filepath.Base(path)
		
		content, err := am.readFileContent(path)
		if err != nil {
			fmt.Printf("Warning: Could not read SystemD unit %s: %v\n", path, err)
			return nil
		}
		isEnabled, isActive := am.getSystemDUnitStatus(unitName)
		
		unit := SystemDUnit{
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

/////detectCronjobs scans for cron job files
func (am *AutomationManager) detectCronjobs() ([]Cronjob, []Cronjob, error) {
	var userCronjobs []Cronjob
	var systemCronjobs []Cronjob
	
	userCronPath := fmt.Sprintf("/var/spool/cron/crontabs/%s", am.username)
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
			userCronjobs = append(userCronjobs, Cronjob{
				Path:    userCronPath,
				Content: content,
				Type:    "user",
			})
		}
	}
	
	
	systemCronPaths := []string{
		"/etc/crontab",
	}
	
	if cronDDir := "/etc/cron.d"; am.dirExists(cronDDir) {
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
				
				systemCronjobs = append(systemCronjobs, Cronjob{
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

func (am *AutomationManager) GetAutomationSummary(data *AutomationData) string {
	var summary strings.Builder
	
	summary.WriteString("Automation Detection Summary:\n")
	summary.WriteString(fmt.Sprintf("- SystemD Services: %d\n", len(data.SystemDServices)))
	summary.WriteString(fmt.Sprintf("- SystemD Timers: %d\n", len(data.SystemDTimers)))
	summary.WriteString(fmt.Sprintf("- User Cronjobs: %d\n", len(data.UserCronjobs)))
	summary.WriteString(fmt.Sprintf("- System Cronjobs: %d\n", len(data.SystemCronjobs)))
	
	if len(data.SystemDServices) > 0 {
		summary.WriteString("\nSystemD Services found:\n")
		for _, service := range data.SystemDServices {
			summary.WriteString(fmt.Sprintf("  - %s (%s)\n", service.Name, service.Path))
		}
	}
	
	if len(data.SystemDTimers) > 0 {
		summary.WriteString("\nSystemD Timers found:\n")
		for _, timer := range data.SystemDTimers {
			summary.WriteString(fmt.Sprintf("  - %s (%s)\n", timer.Name, timer.Path))
		}
	}
	
	if len(data.UserCronjobs) > 0 {
		summary.WriteString("\nUser Cronjobs found:\n")
		for _, cronjob := range data.UserCronjobs {
			summary.WriteString(fmt.Sprintf("  - %s\n", cronjob.Path))
		}
	}
	
	if len(data.SystemCronjobs) > 0 {
		summary.WriteString("\nSystem Cronjobs found:\n")
		for _, cronjob := range data.SystemCronjobs {
			summary.WriteString(fmt.Sprintf("  - %s\n", cronjob.Path))
		}
	}
	
	return summary.String()
}
