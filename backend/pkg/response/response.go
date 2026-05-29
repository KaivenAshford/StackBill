package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageResult struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func OKWithPage(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageResult{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func Fail(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, 401, message)
}

func Forbidden(c *gin.Context, message string) {
	Fail(c, http.StatusForbidden, 403, message)
}

func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, 400, message)
}

func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, 404, message)
}

func InternalError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, 500, message)
}
