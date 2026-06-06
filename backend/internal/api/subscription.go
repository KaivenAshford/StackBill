package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type SubscriptionHandler struct {
	svc *service.SubscriptionService
}

func NewSubscriptionHandler(svc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

// List godoc
// @Summary List subscriptions
// @Description Get paginated list of subscriptions with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param category_id query int false "Filter by category ID"
// @Param status query string false "Filter by status (active, paused, cancelled, expired)"
// @Param upcoming_renewal query bool false "Filter upcoming renewals only"
// @Success 200 {object} response.Response{data=response.PageResult{items=[]dto.SubscriptionResponse}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var query dto.SubscriptionListQuery
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

// GetByID godoc
// @Summary Get subscription by ID
// @Description Get a single subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscription ID"
// @Success 200 {object} response.Response{data=dto.SubscriptionResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
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
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateSubscriptionRequest true "Subscription data"
// @Success 200 {object} response.Response{data=dto.SubscriptionResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreateSubscriptionRequest
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
// @Summary Update subscription
// @Description Update an existing subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscription ID"
// @Param body body dto.UpdateSubscriptionRequest true "Subscription data"
// @Success 200 {object} response.Response{data=dto.SubscriptionResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid id")
		return
	}
	var req dto.UpdateSubscriptionRequest
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
// @Summary Delete subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subscription ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
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
