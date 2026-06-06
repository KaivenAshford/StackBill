package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type ReminderHandler struct {
	svc *service.ReminderService
}

func NewReminderHandler(svc *service.ReminderService) *ReminderHandler {
	return &ReminderHandler{svc: svc}
}

// List godoc
// @Summary List reminders
// @Description Get paginated list of reminders with optional filters
// @Tags reminders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param type query string false "Filter by reminder type (subscription_renewal, asset_expiration, service_warning)"
// @Param is_read query bool false "Filter by read status"
// @Success 200 {object} response.Response{data=response.PageResult{items=[]dto.ReminderResponse}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /reminders [get]
func (h *ReminderHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var query dto.ReminderListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid parameters")
		return
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	result, err := h.svc.List(userID, &query)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, result)
}

// MarkRead godoc
// @Summary Mark reminder as read
// @Description Mark a single reminder as read by its ID
// @Tags reminders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reminder ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /reminders/{id}/read [put]
func (h *ReminderHandler) MarkRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid id")
		return
	}
	if err := h.svc.MarkRead(userID, uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, nil)
}

// MarkAllRead godoc
// @Summary Mark all reminders as read
// @Description Mark all reminders for the current user as read
// @Tags reminders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /reminders/read-all [put]
func (h *ReminderHandler) MarkAllRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.svc.MarkAllRead(userID); err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete godoc
// @Summary Delete (dismiss) reminder
// @Description Dismiss a reminder by its ID
// @Tags reminders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reminder ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /reminders/{id} [delete]
func (h *ReminderHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid id")
		return
	}
	if err := h.svc.Dismiss(userID, uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, nil)
}
