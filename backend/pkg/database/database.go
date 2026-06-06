package database

import (
	"fmt"
	"log/slog"

	"github.com/kingqaquuu/stackbill/internal/config"
	"github.com/kingqaquuu/stackbill/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	slog.Info("database connected", "host", cfg.Host, "dbname", cfg.DBName)
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Subscription{},
		&model.Asset{},
		&model.ReminderRead{},
		&model.ReminderDismissed{},
		&model.NotificationSetting{},
		&model.Webhook{},
	)
}
