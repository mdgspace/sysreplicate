package automation

import (
	"archive/tar"
	"fmt"
	"path/filepath"
	"github.com/mdgspace/sysreplicate/internal/domain"
)

func (am *AutomationManager) BackupAutomation(data *domain.AutomationData, tarWriter *tar.Writer) error {
	fmt.Println("Adding automation files to backup...")
	
	for _, service := range data.SystemDServices {
		if err := am.addFileToTarball(service.Path, service.Content, "automation/systemd/", tarWriter); err != nil {
			fmt.Printf("Warning: Failed to add SystemD service %s: %v\n", service.Name, err)
		}
	}
	for _, timer := range data.SystemDTimers {
		if err := am.addFileToTarball(timer.Path, timer.Content, "automation/systemd/", tarWriter); err != nil {
			fmt.Printf("Warning: Failed to add SystemD timer %s: %v\n", timer.Name, err)
		}
	}
	
	for _, cronjob := range data.UserCronjobs {
		if err := am.addFileToTarball(cronjob.Path, cronjob.Content, "automation/cron/", tarWriter); err != nil {
			fmt.Printf("Warning: Failed to add user cronjob %s: %v\n", cronjob.Path, err)
		}
	}
	
	for _, cronjob := range data.SystemCronjobs {
		if err := am.addFileToTarball(cronjob.Path, cronjob.Content, "automation/cron/", tarWriter); err != nil {
			fmt.Printf("Warning: Failed to add system cronjob %s: %v\n", cronjob.Path, err)
		}
	}
	
	return nil
}
func (am *AutomationManager) addFileToTarball(originalPath, content, tarballPrefix string, tarWriter *tar.Writer) error {
	tarballPath := tarballPrefix + filepath.Base(originalPath)
	header := &tar.Header{
		Name: tarballPath,
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write tar header for %s: %w", tarballPath, err)
	}
	if _, err := tarWriter.Write([]byte(content)); err != nil {
		return fmt.Errorf("failed to write content for %s: %w", tarballPath, err)
	}
	
	return nil
}

func (am *AutomationManager) ValidateAutomationData(data *domain.AutomationData) error {
	unitNames := make(map[string]bool)
	
	for _, service := range data.SystemDServices {
		if unitNames[service.Name] {
			return fmt.Errorf("duplicate SystemD unit name: %s", service.Name)
		}
		unitNames[service.Name] = true
	}
	
	for _, timer := range data.SystemDTimers {
		if unitNames[timer.Name] {
			return fmt.Errorf("duplicate SystemD unit name: %s", timer.Name)
		}
		unitNames[timer.Name] = true
	}
	for _, service := range data.SystemDServices {
		if service.UnitType != "service" && service.UnitType != "target" {
			return fmt.Errorf("invalid unit type for service %s: %s", service.Name, service.UnitType)
		}
	}
	
	for _, timer := range data.SystemDTimers {
		if timer.UnitType != "timer" {
			return fmt.Errorf("invalid unit type for timer %s: %s", timer.Name, timer.UnitType)
		}
	}
	
	return nil
}
