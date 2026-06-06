package service

import (
	"fmt"
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testJWTSecret = "test-secret-key"
const testJWTExpire = 72

type testServices struct {
	DB                  *gorm.DB
	UserRepo            *repository.UserRepository
	CategoryRepo        *repository.CategoryRepository
	SubscriptionRepo    *repository.SubscriptionRepository
	AssetRepo           *repository.AssetRepository
	ReminderRepo        *repository.ReminderRepository
	AuthService         *AuthService
	UserService         *UserService
	CategoryService     *CategoryService
	SubscriptionService *SubscriptionService
	AssetService        *AssetService
	ReminderService     *ReminderService
	DashboardService    *DashboardService
}

var testDBCounter int

func newTestServices(t *testing.T) *testServices {
	t.Helper()
	testDBCounter++
	dsn := fmt.Sprintf("file:testdb%d?mode=memory&cache=shared", testDBCounter)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Subscription{},
		&model.Asset{},
		&model.ReminderRead{},
		&model.ReminderDismissed{},
	); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	reminderRepo := repository.NewReminderRepository(db)

	authService := NewAuthService(userRepo, categoryRepo, testJWTSecret, testJWTExpire)
	userService := NewUserService(userRepo)
	categoryService := NewCategoryService(categoryRepo)
	subscriptionService := NewSubscriptionService(subscriptionRepo)
	assetService := NewAssetService(assetRepo)
	reminderService := NewReminderService(reminderRepo)
	dashboardService := NewDashboardService(subscriptionRepo, assetRepo, reminderRepo, categoryRepo, subscriptionService)

	return &testServices{
		DB:                  db,
		UserRepo:            userRepo,
		CategoryRepo:        categoryRepo,
		SubscriptionRepo:    subscriptionRepo,
		AssetRepo:           assetRepo,
		ReminderRepo:        reminderRepo,
		AuthService:         authService,
		UserService:         userService,
		CategoryService:     categoryService,
		SubscriptionService: subscriptionService,
		AssetService:        assetService,
		ReminderService:     reminderService,
		DashboardService:    dashboardService,
	}
}

func registerTestUser(t *testing.T, svc *testServices, username string) uint {
	t.Helper()
	resp, err := svc.AuthService.Register(&dto.RegisterRequest{
		Username: username,
		Email:    username + "@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register test user %s: %v", username, err)
	}
	return resp.User.ID
}
