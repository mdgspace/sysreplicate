package automation

import (
	"fmt"
	"os"
	"strings"

	"github.com/mdgspace/sysreplicate/internal/domain"
)

type AutomationManager struct {
	username string
}
func NewAutomationManager() *AutomationManager {
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME")
	}
	
	return &AutomationManager{
		username: username,
	}
}
func (am *AutomationManager) DetectAutomation() (*domain.AutomationData, error) {
	fmt.Println("Detecting automation files...")
	
	data := &domain.AutomationData{
		SystemDServices: make([]domain.SystemDUnit, 0),
		SystemDTimers:   make([]domain.SystemDUnit, 0),
		UserCronjobs:    make([]domain.Cronjob, 0),
		SystemCronjobs:  make([]domain.Cronjob, 0),
	}
	
	systemdServices, systemdTimers, err := am.detectSystemDUnits()
	if err != nil {
		fmt.Printf("Warning: Failed to detect SystemD units: %v\n", err)
	} else {
		data.SystemDServices = systemdServices
		data.SystemDTimers = systemdTimers
	}

	// usercustom, systemCronjobs, err := am.detectCronjobs()
	// if err != nil {
	// 	fmt.Printf("Warning: Failed to detect customs: %v\n", err)
	// } else {
	// 	data.UserCronjobs = usercustoms
	// 	data.SystemCronjobs = 
	// }
	
	userCronjobs, systemCronjobs, err := am.detectCronjobs()
	if err != nil {
		fmt.Printf("Warning: Failed to detect cronjobs: %v\n", err)
	} else {
		data.UserCronjobs = userCronjobs
		data.SystemCronjobs = systemCronjobs
	}
	
	fmt.Printf("Detected %d SystemD services, %d SystemD timers, %d user cronjobs, %d system cronjobs\n",
		len(data.SystemDServices), len(data.SystemDTimers), len(data.UserCronjobs), len(data.SystemCronjobs))



	if len(data.SystemDServices) > 0 {
		fmt.Println("  SystemD Services found:")
		for _, service := range data.SystemDServices {
			fmt.Printf("    - %s (%s)\n", service.Name, service.Path)
		}
	}
	
	if len(data.SystemDTimers) > 0 {
		fmt.Println("  SystemD Timers found:")
		for _, timer := range data.SystemDTimers {
			fmt.Printf("    - %s (%s)\n", timer.Name, timer.Path)
		}
	}
	
	if len(data.UserCronjobs) > 0 {
		fmt.Println("  User Cronjobs found:")
		for _, cronjob := range data.UserCronjobs {
			fmt.Printf("    - %s\n", cronjob.Path)
		}
	}
	
	if len(data.SystemCronjobs) > 0 {
		fmt.Println("  System Cronjobs found:")
		for _, cronjob := range data.SystemCronjobs {
			fmt.Printf("    - %s\n", cronjob.Path)
		}
	}
	
	return data, nil
}
/////symlink logicc
func (am *AutomationManager) isCustomSystemDUnit(filePath string) bool {
	// Check if it's a symlink
	linkInfo, err := os.Lstat(filePath)
	if err != nil {
		return false
	}
	
	// If it's not a symlink, it's custom
	if linkInfo.Mode()&os.ModeSymlink == 0 {
		return true
	}
	
	// If it's a symlink, check if it points to package-managed directory
	target, err := os.Readlink(filePath)
	if err != nil {
		return false
	}
	
	for _, dir := range domain.PackageManagedDirs {
		if strings.HasPrefix(target, dir) {
			return false
		}
	}
	
	return true
}

func (am *AutomationManager) readFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// getSystemDUnitStatus checks if a SystemD unit is enabled and active
func (am *AutomationManager) getSystemDUnitStatus(unitName string) (bool, bool) {
	///TODO(@jaadu): getSystemDUnitStatus checks if a SystemD unit is enabled and active: THIS IS A Simplified implementation you should use SYSTEMCTL COMMANDS in better implementation
	//AS enabled services logic is different
	return false, false
}
