package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/api"
	"github.com/kingqaquuu/stackbill/internal/middleware"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"github.com/kingqaquuu/stackbill/internal/service"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, jwtSecret string, jwtExpireHours int) {
	r.Use(middleware.CORSMiddleware())

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	authService := service.NewAuthService(userRepo, categoryRepo, jwtSecret, jwtExpireHours)
	userService := service.NewUserService(userRepo)

	authHandler := api.NewAuthHandler(authService, userService)
	userHandler := api.NewUserHandler(userService)

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"code": 0, "message": "ok"})
		})

		auth := apiGroup.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	authorized := apiGroup.Group("")
	authorized.Use(middleware.JWTAuth(jwtSecret))
	{
		authorized.GET("/auth/me", authHandler.GetCurrentUser)
		authorized.PUT("/users/profile", userHandler.UpdateProfile)
		authorized.PUT("/users/password", userHandler.UpdatePassword)
	}
}
