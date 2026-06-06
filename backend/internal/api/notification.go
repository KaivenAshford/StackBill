package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type NotificationHandler struct {
	notificationSvc *service.NotificationService
}

func NewNotificationHandler(notificationSvc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationSvc: notificationSvc}
}

// GetNotificationSetting godoc
// @Summary Get notification settings
// @Description Get the current user's notification preferences
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.NotificationSettingResponse}
// @Router /notification-settings [get]
func (h *NotificationHandler) GetNotificationSetting(c *gin.Context) {
	userID := c.GetUint("user_id")
	result, err := h.notificationSvc.Get(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, result)
}

// UpdateNotificationSetting godoc
// @Summary Update notification settings
// @Description Update the current user's notification preferences
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.UpdateNotificationSettingRequest true "Notification settings"
// @Success 200 {object} response.Response{data=dto.NotificationSettingResponse}
// @Router /notification-settings [put]
func (h *NotificationHandler) UpdateNotificationSetting(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateNotificationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid parameters")
		return
	}
	result, err := h.notificationSvc.Update(userID, &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, result)
}

type WebhookHandler struct {
	webhookSvc *service.WebhookService
}

func NewWebhookHandler(webhookSvc *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{webhookSvc: webhookSvc}
}

// ListWebhooks godoc
// @Summary List webhooks
// @Description Get all webhooks for the current user
// @Tags webhooks
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]dto.WebhookResponse}
// @Router /webhooks [get]
func (h *WebhookHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	result, err := h.webhookSvc.List(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, result)
}

// CreateWebhook godoc
// @Summary Create webhook
// @Description Create a new webhook endpoint
// @Tags webhooks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateWebhookRequest true "Webhook data"
// @Success 200 {object} response.Response{data=dto.WebhookResponse}
// @Router /webhooks [post]
func (h *WebhookHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid parameters")
		return
	}
	result, err := h.webhookSvc.Create(userID, &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, result)
}

// UpdateWebhook godoc
// @Summary Update webhook
// @Description Update an existing webhook
// @Tags webhooks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Webhook ID"
// @Param body body dto.UpdateWebhookRequest true "Webhook data"
// @Success 200 {object} response.Response{data=dto.WebhookResponse}
// @Router /webhooks/{id} [put]
func (h *WebhookHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid parameters")
		return
	}
	result, err := h.webhookSvc.Update(userID, uint(id), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, result)
}

// DeleteWebhook godoc
// @Summary Delete webhook
// @Description Delete a webhook endpoint
// @Tags webhooks
// @Produce json
// @Security BearerAuth
// @Param id path int true "Webhook ID"
// @Success 200 {object} response.Response
// @Router /webhooks/{id} [delete]
func (h *WebhookHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.webhookSvc.Delete(userID, uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, nil)
}
