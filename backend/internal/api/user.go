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

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the currently authenticated user's nickname and avatar
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.UpdateProfileRequest true "Profile data"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid parameters")
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

// UpdatePassword godoc
// @Summary Update password
// @Description Change the currently authenticated user's password
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.UpdatePasswordRequest true "Password data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /users/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid parameters")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.userService.UpdatePassword(userID, &req); err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, nil)
}
