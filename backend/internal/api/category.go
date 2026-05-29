package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var query dto.CategoryListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}
	resp, err := h.svc.List(userID, &query)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, 40001, "invalid id")
		return
	}
	resp, err := h.svc.GetByID(userID, uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}
	resp, err := h.svc.Create(userID, &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, 40001, "invalid id")
		return
	}
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}
	resp, err := h.svc.Update(userID, uint(id), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, 40001, "invalid id")
		return
	}
	if err := h.svc.Delete(userID, uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, nil)
}
