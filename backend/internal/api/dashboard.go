package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := h.svc.GetDashboard(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}
