// StackBill API
//
// @title StackBill API
// @version 1.0
// @description Digital asset and subscription management platform API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/config"
	"github.com/kingqaquuu/stackbill/internal/router"
	"github.com/kingqaquuu/stackbill/pkg/database"
	"github.com/kingqaquuu/stackbill/pkg/logger"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("load config failed", "error", err)
		return
	}

	logger.Init(&cfg.Log)

	if err := database.Init(&cfg.Database); err != nil {
		slog.Error("init database failed", "error", err)
		return
	}

	if err := database.AutoMigrate(); err != nil {
		slog.Error("auto migrate failed", "error", err)
		return
	}

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	router.Setup(r, database.DB, cfg.JWT.Secret, cfg.JWT.ExpireHours, cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	slog.Info("server starting", "addr", addr)
	if err := r.Run(addr); err != nil {
		slog.Error("start server failed", "error", err)
	}
}
