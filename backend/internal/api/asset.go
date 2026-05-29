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

func (h *AssetHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var query dto.AssetListQuery
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

func (h *AssetHandler) GetByID(c *gin.Context) {
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

func (h *AssetHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreateAssetRequest
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

func (h *AssetHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, 40001, "invalid id")
		return
	}
	var req dto.UpdateAssetRequest
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

func (h *AssetHandler) Delete(c *gin.Context) {
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
