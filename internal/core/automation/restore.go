package automation

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mdgspace/sysreplicate/internal/domain"
)

///TODO(@jaadu): IMPROVE THE RESTORE LOGIC AND RESTORATION COMMAND
func (am *AutomationManager) GenerateRestorationCommands(data *domain.AutomationData) []string {
	var commands []string
	if len(data.SystemDServices) > 0 || len(data.SystemDTimers) > 0 {
		commands = append(commands, "echo 'Restoring SystemD units...'")
		
		for _, service := range data.SystemDServices {
			commands = append(commands, fmt.Sprintf("sudo cp automation/systemd/%s %s", 
				filepath.Base(service.Path), service.Path))
		}
		for _, timer := range data.SystemDTimers {
			commands = append(commands, fmt.Sprintf("sudo cp automation/systemd/%s %s", 
				filepath.Base(timer.Path), timer.Path))
		}
		


		// Reload SystemD daemon
		commands = append(commands, "sudo systemctl daemon-reload")
		// Enable and start services
		for _, service := range data.SystemDServices {
			if service.UnitType == "service" {
				commands = append(commands, fmt.Sprintf("sudo systemctl enable --now %s || true", 
					strings.TrimSuffix(service.Name, ".service")))
			}
		}
		// Enable and start timers
		for _, timer := range data.SystemDTimers {
			commands = append(commands, fmt.Sprintf("sudo systemctl enable --now %s || true", 
				strings.TrimSuffix(timer.Name, ".timer")))
		}
	}
	if len(data.UserCronjobs) > 0 || len(data.SystemCronjobs) > 0 {
		commands = append(commands, "echo 'Restoring cronjobs...'")
		for _, cronjob := range data.UserCronjobs {
			if cronjob.Type == "user" {
				commands = append(commands, fmt.Sprintf("crontab automation/cron/%s || true", 
					filepath.Base(cronjob.Path)))
			}
		}
		for _, cronjob := range data.SystemCronjobs {
			if cronjob.Type == "system" {
				commands = append(commands, fmt.Sprintf("sudo cp automation/cron/%s %s", 
					filepath.Base(cronjob.Path), cronjob.Path))
			} else if cronjob.Type == "cron_d" {
				commands = append(commands, fmt.Sprintf("sudo cp automation/cron/%s %s", 
					filepath.Base(cronjob.Path), cronjob.Path))
			}
		}
	}
	
	return commands
}