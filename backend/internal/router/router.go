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

	// Repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	reminderRepo := repository.NewReminderRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, categoryRepo, jwtSecret, jwtExpireHours)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	assetService := service.NewAssetService(assetRepo)
	reminderService := service.NewReminderService(reminderRepo)

	// Handlers
	authHandler := api.NewAuthHandler(authService, userService)
	userHandler := api.NewUserHandler(userService)
	categoryHandler := api.NewCategoryHandler(categoryService)
	subscriptionHandler := api.NewSubscriptionHandler(subscriptionService)
	assetHandler := api.NewAssetHandler(assetService)
	reminderHandler := api.NewReminderHandler(reminderService)

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

		// Categories
		authorized.GET("/categories", categoryHandler.List)
		authorized.GET("/categories/:id", categoryHandler.GetByID)
		authorized.POST("/categories", categoryHandler.Create)
		authorized.PUT("/categories/:id", categoryHandler.Update)
		authorized.DELETE("/categories/:id", categoryHandler.Delete)

		// Subscriptions
		authorized.GET("/subscriptions", subscriptionHandler.List)
		authorized.GET("/subscriptions/:id", subscriptionHandler.GetByID)
		authorized.POST("/subscriptions", subscriptionHandler.Create)
		authorized.PUT("/subscriptions/:id", subscriptionHandler.Update)
		authorized.DELETE("/subscriptions/:id", subscriptionHandler.Delete)

		// Assets
		authorized.GET("/assets", assetHandler.List)
		authorized.GET("/assets/:id", assetHandler.GetByID)
		authorized.POST("/assets", assetHandler.Create)
		authorized.PUT("/assets/:id", assetHandler.Update)
		authorized.DELETE("/assets/:id", assetHandler.Delete)

		// Reminders
		authorized.GET("/reminders", reminderHandler.List)
		authorized.PUT("/reminders/:id/read", reminderHandler.MarkRead)
		authorized.PUT("/reminders/read-all", reminderHandler.MarkAllRead)
	}
}
