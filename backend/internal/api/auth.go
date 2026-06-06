package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "Registration data"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "参数校验失败"
		if strings.Contains(err.Error(), "Username") {
			msg = "用户名需 3-50 个字符"
		} else if strings.Contains(err.Error(), "Email") {
			msg = "邮箱格式不正确"
		} else if strings.Contains(err.Error(), "Password") {
			msg = "密码需 6-50 个字符"
		}
		response.Fail(c, 400, service.ErrCodeInvalidParams, msg)
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

// Login godoc
// @Summary Login
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "用户名和密码不能为空")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get the currently authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := h.userService.GetCurrentUser(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}
