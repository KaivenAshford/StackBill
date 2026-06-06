package dto

type NotificationSettingResponse struct {
	ID               uint `json:"id"`
	EmailEnabled     bool `json:"email_enabled"`
	RemindDaysBefore int  `json:"remind_days_before"`
}

type UpdateNotificationSettingRequest struct {
	EmailEnabled     *bool `json:"email_enabled"`
	RemindDaysBefore *int  `json:"remind_days_before" binding:"omitempty,min=1,max=30"`
}

type WebhookResponse struct {
	ID     uint   `json:"id"`
	URL    string `json:"url"`
	Events string `json:"events"`
	Active bool   `json:"active"`
}

type CreateWebhookRequest struct {
	URL    string `json:"url" binding:"required,max=500"`
	Secret string `json:"secret" binding:"max=200"`
	Events string `json:"events" binding:"required,max=500"`
}

type UpdateWebhookRequest struct {
	URL    string `json:"url" binding:"max=500"`
	Secret string `json:"secret" binding:"max=200"`
	Events string `json:"events" binding:"max=500"`
	Active *bool  `json:"active"`
}
