package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

func handleServiceError(c *gin.Context, err error) {
	if svcErr, ok := err.(*service.ServiceError); ok {
		response.Fail(c, svcErr.HTTPCode, svcErr.Code, svcErr.Message)
		return
	}
	response.InternalError(c, "internal server error")
}
