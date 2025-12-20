package domain

type AutomationData struct {
	SystemDServices []SystemDUnit `json:"systemd_services"`
	SystemDTimers   []SystemDUnit `json:"systemd_timers"`
	UserCronjobs    []Cronjob     `json:"user_cronjobs"`
	SystemCronjobs  []Cronjob     `json:"system_cronjobs"`
}
type SystemDUnit struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Content     string `json:"content"`
	UnitType    string `json:"unit_type"` ////saare service, timer and target available ere
	IsEnabled   bool   `json:"is_enabled"`
	IsActive    bool   `json:"is_active"`
}
type Cronjob struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Type    string `json:"type"` //.//user, system, cron_d
}