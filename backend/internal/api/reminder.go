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

func (h *ReminderHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var query dto.ReminderListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
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

func (h *ReminderHandler) MarkRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, 40001, "invalid id")
		return
	}
	if err := h.svc.MarkRead(userID, uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *ReminderHandler) MarkAllRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.svc.MarkAllRead(userID); err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, nil)
}
