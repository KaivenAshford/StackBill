package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}

	userID := c.GetUint("user_id")
	resp, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.userService.UpdatePassword(userID, &req); err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, nil)
}
