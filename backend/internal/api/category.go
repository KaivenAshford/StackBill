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

// List godoc
// @Summary List categories
// @Description Get all categories for the current user, optionally filtered by type
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filter by type (subscription, asset)"
// @Success 200 {object} response.Response{data=[]dto.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var query dto.CategoryListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid parameters")
		return
	}
	resp, err := h.svc.List(userID, &query)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}

// GetByID godoc
// @Summary Get category by ID
// @Description Get a single category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response{data=dto.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid id")
		return
	}
	resp, err := h.svc.GetByID(userID, uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}

// Create godoc
// @Summary Create category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateCategoryRequest true "Category data"
// @Success 200 {object} response.Response{data=dto.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid parameters")
		return
	}
	resp, err := h.svc.Create(userID, &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}

// Update godoc
// @Summary Update category
// @Description Update an existing category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param body body dto.UpdateCategoryRequest true "Category data"
// @Success 200 {object} response.Response{data=dto.CategoryResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid id")
		return
	}
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid parameters")
		return
	}
	resp, err := h.svc.Update(userID, uint(id), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}

// Delete godoc
// @Summary Delete category
// @Description Delete a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid id")
		return
	}
	if err := h.svc.Delete(userID, uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, nil)
}
