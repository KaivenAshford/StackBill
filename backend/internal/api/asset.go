package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type AssetHandler struct {
	svc *service.AssetService
}

func NewAssetHandler(svc *service.AssetService) *AssetHandler {
	return &AssetHandler{svc: svc}
}

// List godoc
// @Summary List assets
// @Description Get paginated list of digital assets with optional filters
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param asset_type query string false "Filter by asset type (domain, server, docker_service, ssl_certificate, api_key, repository, other)"
// @Param status query string false "Filter by status (active, inactive, expired, warning)"
// @Param expiring_days query int false "Filter assets expiring within N days"
// @Param keyword query string false "Search by name"
// @Success 200 {object} response.Response{data=response.PageResult{items=[]dto.AssetResponse}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /assets [get]
func (h *AssetHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var query dto.AssetListQuery
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
// @Summary Get asset by ID
// @Description Get a single digital asset by its ID
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Asset ID"
// @Success 200 {object} response.Response{data=dto.AssetResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /assets/{id} [get]
func (h *AssetHandler) GetByID(c *gin.Context) {
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
// @Summary Create asset
// @Description Create a new digital asset
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateAssetRequest true "Asset data"
// @Success 200 {object} response.Response{data=dto.AssetResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /assets [post]
func (h *AssetHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreateAssetRequest
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
// @Summary Update asset
// @Description Update an existing digital asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Asset ID"
// @Param body body dto.UpdateAssetRequest true "Asset data"
// @Success 200 {object} response.Response{data=dto.AssetResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /assets/{id} [put]
func (h *AssetHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "invalid id")
		return
	}
	var req dto.UpdateAssetRequest
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
// @Summary Delete asset
// @Description Delete a digital asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Asset ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /assets/{id} [delete]
func (h *AssetHandler) Delete(c *gin.Context) {
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
