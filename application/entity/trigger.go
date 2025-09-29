package entity

type TriggerData struct {
	Type     string `json:"type" binding:"required,oneof=one-off cron"`
	DateTime string `json:"datetime,omitempty"`
	Cron     string `json:"cron,omitempty"`
}
