package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/config"
	"github.com/kingqaquuu/stackbill/internal/router"
	"github.com/kingqaquuu/stackbill/pkg/database"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("init database: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	router.Setup(r, database.DB, cfg.JWT.Secret, cfg.JWT.ExpireHours)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
