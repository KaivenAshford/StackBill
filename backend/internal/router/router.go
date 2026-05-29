package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/middleware"
)

func Setup(r *gin.Engine, jwtSecret string) {
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"code": 0, "message": "ok"})
		})
	}

	auth := api.Group("")
	auth.Use(middleware.JWTAuth(jwtSecret))
	_ = auth
}
