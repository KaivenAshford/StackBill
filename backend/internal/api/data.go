package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type ExportHandler struct {
	exportSvc *service.ExportService
}

func NewExportHandler(exportSvc *service.ExportService) *ExportHandler {
	return &ExportHandler{exportSvc: exportSvc}
}

// ExportSubscriptions godoc
// @Summary Export subscriptions as CSV
// @Description Download all subscriptions as a CSV file
// @Tags subscriptions
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {file} file
// @Failure 401 {object} response.Response
// @Router /subscriptions/export [get]
func (h *ExportHandler) ExportSubscriptions(c *gin.Context) {
	userID := c.GetUint("user_id")
	data, err := h.exportSvc.ExportSubscriptionsCSV(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=subscriptions.csv")
	c.Data(200, "text/csv", data)
}

// ExportAssets godoc
// @Summary Export assets as CSV
// @Description Download all assets as a CSV file
// @Tags assets
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {file} file
// @Failure 401 {object} response.Response
// @Router /assets/export [get]
func (h *ExportHandler) ExportAssets(c *gin.Context) {
	userID := c.GetUint("user_id")
	data, err := h.exportSvc.ExportAssetsCSV(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=assets.csv")
	c.Data(200, "text/csv", data)
}

type ImportHandler struct {
	importSvc *service.ImportService
}

func NewImportHandler(importSvc *service.ImportService) *ImportHandler {
	return &ImportHandler{importSvc: importSvc}
}

// ImportSubscriptions godoc
// @Summary Import subscriptions from CSV
// @Description Upload a CSV file to batch create subscriptions
// @Tags subscriptions
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "CSV file"
// @Success 200 {object} response.Response{data=service.ImportResult}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /subscriptions/import [post]
func (h *ImportHandler) ImportSubscriptions(c *gin.Context) {
	userID := c.GetUint("user_id")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, fmt.Sprintf("failed to read file: %v", err))
		return
	}
	defer file.Close()

	result, err := h.importSvc.ImportSubscriptionsCSV(userID, file)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, result)
}

// ImportAssets godoc
// @Summary Import assets from CSV
// @Description Upload a CSV file to batch create assets
// @Tags assets
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "CSV file"
// @Success 200 {object} response.Response{data=service.ImportResult}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /assets/import [post]
func (h *ImportHandler) ImportAssets(c *gin.Context) {
	userID := c.GetUint("user_id")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, fmt.Sprintf("failed to read file: %v", err))
		return
	}
	defer file.Close()

	result, err := h.importSvc.ImportAssetsCSV(userID, file)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, result)
}
