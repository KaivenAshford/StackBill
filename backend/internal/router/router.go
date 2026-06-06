package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/api"
	"github.com/kingqaquuu/stackbill/internal/config"
	"github.com/kingqaquuu/stackbill/internal/middleware"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/internal/task"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/kingqaquuu/stackbill/docs"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, jwtSecret string, jwtExpireHours int, cfg *config.Config) {
	r.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	reminderRepo := repository.NewReminderRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	webhookRepo := repository.NewWebhookRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, categoryRepo, jwtSecret, jwtExpireHours)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	assetService := service.NewAssetService(assetRepo)
	reminderService := service.NewReminderService(reminderRepo)
	dashboardService := service.NewDashboardService(subscriptionRepo, assetRepo, reminderRepo, categoryRepo, subscriptionService)
	exportService := service.NewExportService(subscriptionRepo, assetRepo)
	importService := service.NewImportService(subscriptionRepo, assetRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	webhookService := service.NewWebhookService(webhookRepo)

	// Handlers
	authHandler := api.NewAuthHandler(authService, userService)
	userHandler := api.NewUserHandler(userService)
	categoryHandler := api.NewCategoryHandler(categoryService)
	subscriptionHandler := api.NewSubscriptionHandler(subscriptionService)
	assetHandler := api.NewAssetHandler(assetService)
	reminderHandler := api.NewReminderHandler(reminderService)
	dashboardHandler := api.NewDashboardHandler(dashboardService)
	exportHandler := api.NewExportHandler(exportService)
	importHandler := api.NewImportHandler(importService)
	notificationHandler := api.NewNotificationHandler(notificationService)
	webhookHandler := api.NewWebhookHandler(webhookService)

	// Start email reminder scheduler
	scheduler := task.NewScheduler(notificationRepo, subscriptionRepo, assetRepo, userRepo, cfg.SMTP)
	scheduler.Start()

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

		// Dashboard
		authorized.GET("/dashboard", dashboardHandler.GetDashboard)

		// Categories
		authorized.GET("/categories", categoryHandler.List)
		authorized.GET("/categories/:id", categoryHandler.GetByID)
		authorized.POST("/categories", categoryHandler.Create)
		authorized.PUT("/categories/:id", categoryHandler.Update)
		authorized.DELETE("/categories/:id", categoryHandler.Delete)

		// Subscriptions
		authorized.GET("/subscriptions", subscriptionHandler.List)
		authorized.GET("/subscriptions/export", exportHandler.ExportSubscriptions)
		authorized.POST("/subscriptions/import", importHandler.ImportSubscriptions)
		authorized.GET("/subscriptions/:id", subscriptionHandler.GetByID)
		authorized.POST("/subscriptions", subscriptionHandler.Create)
		authorized.PUT("/subscriptions/:id", subscriptionHandler.Update)
		authorized.DELETE("/subscriptions/:id", subscriptionHandler.Delete)

		// Assets
		authorized.GET("/assets", assetHandler.List)
		authorized.GET("/assets/export", exportHandler.ExportAssets)
		authorized.POST("/assets/import", importHandler.ImportAssets)
		authorized.GET("/assets/:id", assetHandler.GetByID)
		authorized.POST("/assets", assetHandler.Create)
		authorized.PUT("/assets/:id", assetHandler.Update)
		authorized.DELETE("/assets/:id", assetHandler.Delete)

		// Reminders
		authorized.GET("/reminders", reminderHandler.List)
		authorized.PUT("/reminders/:id/read", reminderHandler.MarkRead)
		authorized.PUT("/reminders/read-all", reminderHandler.MarkAllRead)
		authorized.DELETE("/reminders/:id", reminderHandler.Delete)

		// Notification Settings
		authorized.GET("/notification-settings", notificationHandler.GetNotificationSetting)
		authorized.PUT("/notification-settings", notificationHandler.UpdateNotificationSetting)

		// Webhooks
		authorized.GET("/webhooks", webhookHandler.List)
		authorized.POST("/webhooks", webhookHandler.Create)
		authorized.PUT("/webhooks/:id", webhookHandler.Update)
		authorized.DELETE("/webhooks/:id", webhookHandler.Delete)
	}
}
