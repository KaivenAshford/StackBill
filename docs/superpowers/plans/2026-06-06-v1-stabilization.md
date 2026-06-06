# V1 Stabilization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix all known defects, unify migrations, productionize Docker, add Swagger API docs, add backend tests, add CI — ship a stable V1.

**Architecture:** Five independent work packages executed sequentially. Each package produces a self-contained, committable change. No cross-package dependencies except Task 1 (frontend fixes) must land before Task 5 (CI) to ensure `npm run build` passes.

**Tech Stack:** Go 1.25 / Gin / GORM / swaggo/swag / SQLite (tests) / GitHub Actions / nginx / Docker Compose

---

## Task 1: Fix Frontend Defects

**Files:**
- Modify: `frontend/src/utils/request.ts`
- Modify: `frontend/src/stores/dashboard.ts`
- Modify: `frontend/src/stores/subscription.ts`
- Modify: `frontend/src/stores/asset.ts`
- Modify: `frontend/src/stores/category.ts`
- Modify: `frontend/src/stores/user.ts`
- Modify: `frontend/src/locales/zh-CN.ts`
- Modify: `frontend/src/locales/en-US.ts`
- Modify: `frontend/src/router/index.ts`

### 1.1 Fix Store data access — the interceptor already unwraps `{code, message, data}`

- [ ] **Step 1: Fix `stores/dashboard.ts`**

The interceptor in `request.ts` returns `response.data` (the inner `{code,message,data}` object). When `code === 0`, it returns `data` directly. So `res` in the store is already `DashboardData`, not `AxiosResponse`. Remove the extra `.data`.

```ts
// frontend/src/stores/dashboard.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getDashboard, type DashboardData } from '@/api/dashboard'

export const useDashboardStore = defineStore('dashboard', () => {
  const data = ref<DashboardData | null>(null)
  const loaded = ref(false)

  async function ensureLoaded() {
    if (loaded.value) return
    const res = await getDashboard()
    data.value = res.data as unknown as DashboardData
    loaded.value = true
  }

  function invalidate() {
    data.value = null
    loaded.value = false
  }

  async function refresh() {
    invalidate()
    await ensureLoaded()
  }

  return { data, loaded, ensureLoaded, invalidate, refresh }
})
```

Wait — actually re-reading the interceptor: it returns `data` which is `response.data` (the full `{code,message,data}` object). Then it checks `data.code !== 0` and rejects. If `data.code === 0`, it returns `data` (which is `{code:0, message:"success", data: {...}}`). So `res` in the store IS `{code:0, message:"success", data: {...}}`, meaning `res.data` is correct.

Let me re-verify:

```ts
// request.ts interceptor
(response) => {
    const data = response.data          // data = {code:0, message:"success", data:{...}}
    if (data.code !== undefined && data.code !== 0) {
      return Promise.reject(...)
    }
    return data                          // returns {code:0, message:"success", data:{...}}
}
```

So `res` in store = `{code:0, message:"success", data:{...}}`, and `res.data` is the actual payload. The stores are CORRECT. No fix needed here.

- [ ] **Step 1 (revised): Verify stores are correct — no change needed**

After careful re-read: the interceptor returns `response.data` which is `{code, message, data}`. The stores access `res.data` to get the inner payload. This is correct. No store changes needed.

- [ ] **Step 2: Fix i18n duplicate key in `frontend/src/locales/zh-CN.ts`**

Replace the duplicate `cycle` key. Change line 64 (`cycle: '计费周期'`) to `billingCycle` and line 76 (`cycle: '自定义'`) to `customCycle`:

The `subscription` section should become:

```ts
  subscription: {
    name: '订阅名称',
    amount: '金额',
    currency: '币种',
    billingCycle: '计费周期',
    status: '状态',
    nextPayment: '下次付款',
    category: '分类',
    startDate: '开始日期',
    url: '网址',
    remark: '备注',
    weekly: '每周',
    monthly: '每月',
    quarterly: '每季度',
    yearly: '每年',
    oneTime: '一次性',
    customCycle: '自定义',
    active: '活跃',
    paused: '已暂停',
    cancelled: '已取消',
    expired: '已过期',
  },
```

- [ ] **Step 3: Fix i18n duplicate key in `frontend/src/locales/en-US.ts`**

Same change:

```ts
  subscription: {
    name: 'Name',
    amount: 'Amount',
    currency: 'Currency',
    billingCycle: 'Billing Cycle',
    status: 'Status',
    nextPayment: 'Next Payment',
    category: 'Category',
    startDate: 'Start Date',
    url: 'URL',
    remark: 'Remark',
    weekly: 'Weekly',
    monthly: 'Monthly',
    quarterly: 'Quarterly',
    yearly: 'Yearly',
    oneTime: 'One Time',
    customCycle: 'Custom',
    active: 'Active',
    paused: 'Paused',
    cancelled: 'Cancelled',
    expired: 'Expired',
  },
```

- [ ] **Step 4: Update any Vue templates that reference the old key**

Search all `.vue` files for `t('subscription.cycle')` and update to `t('subscription.billingCycle')` or `t('subscription.customCycle')` as appropriate based on context.

- [ ] **Step 5: Fix 401 hard redirect in `frontend/src/utils/request.ts`**

Replace `window.location.href = '/login'` with Vue Router navigation:

```ts
import axios from 'axios'
import router from '@/router'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
})

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

request.interceptors.response.use(
  (response) => {
    const data = response.data
    if (data.code !== undefined && data.code !== 0) {
      return Promise.reject(new Error(data.message || 'request failed'))
    }
    return data
  },
  (error) => {
    const msg = error.response?.data?.message || error.message || 'request failed'
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      router.push('/login')
    }
    return Promise.reject(new Error(msg))
  },
)

export default request
```

- [ ] **Step 6: Add 404 catch-all route in `frontend/src/router/index.ts`**

Add a catch-all as the last child of the main layout route:

```ts
        { path: 'settings', name: 'Settings', component: () => import('@/views/settings/Index.vue') },
        { path: ':pathMatch(.*)*', name: 'NotFound', redirect: '/' },
```

- [ ] **Step 7: Verify frontend builds**

Run: `cd frontend && npm run build`
Expected: Build succeeds with no type errors.

- [ ] **Step 8: Commit**

```bash
git add frontend/
git commit -m "fix: frontend defects — i18n duplicate key, 401 redirect, 404 route"
```

---

## Task 2: Unify Migrations — Remove SQL Files, Clean Up Unused Reminder Model

**Files:**
- Delete: `backend/migrations/` (entire directory)
- Delete: `backend/internal/model/reminder.go`
- Modify: `backend/pkg/database/database.go`

- [ ] **Step 1: Delete SQL migration files**

```bash
rm -rf backend/migrations/
```

- [ ] **Step 2: Delete unused `model/reminder.go`**

```bash
rm backend/internal/model/reminder.go
```

- [ ] **Step 3: Update `database.go` to remove Reminder from AutoMigrate**

```go
// backend/pkg/database/database.go
package database

import (
	"fmt"
	"log"

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

	log.Println("database connected")
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
	)
}
```

- [ ] **Step 4: Verify backend builds**

Run: `cd backend && go build ./...`
Expected: Build succeeds.

- [ ] **Step 5: Update README — remove migration file references**

In `README.md`, remove the section about SQL migration files and the `backend/migrations/` reference in the project structure. Update the "数据库迁移" section to only mention AutoMigrate.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: remove SQL migrations and unused Reminder model, use AutoMigrate only"
```

---

## Task 3: Docker Compose Production Readiness

**Files:**
- Modify: `docker-compose.yml`
- Modify: `frontend/nginx.conf`
- Create: `backend/.dockerignore`
- Create: `frontend/.dockerignore`

- [ ] **Step 1: Create `backend/.dockerignore`**

```
.git
.gitignore
*.md
config.yaml
.env
```

- [ ] **Step 2: Create `frontend/.dockerignore`**

```
.git
.gitignore
node_modules
dist
*.md
.env
```

- [ ] **Step 3: Update `frontend/nginx.conf` with gzip, security headers, caching**

```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript image/svg+xml;
    gzip_min_length 256;

    # Security headers
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-XSS-Protection "1; mode=block" always;

    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # Static asset caching
    location ~* \.(js|css|svg|png|jpg|jpeg|gif|ico|woff2?)$ {
        expires 7d;
        add_header Cache-Control "public, immutable";
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

- [ ] **Step 4: Update `docker-compose.yml`**

```yaml
services:
  postgres:
    image: postgres:17-alpine
    environment:
      POSTGRES_USER: ${DB_USER:-stackbill}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-stackbill_password}
      POSTGRES_DB: ${DB_NAME:-stackbill}
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-stackbill}"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "${SERVER_PORT:-8080}:8080"
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER:-stackbill}
      DB_PASSWORD: ${DB_PASSWORD:-stackbill_password}
      DB_NAME: ${DB_NAME:-stackbill}
      JWT_SECRET: ${JWT_SECRET:-change-me-in-production}
      SERVER_PORT: 8080
    depends_on:
      postgres:
        condition: service_healthy
    volumes:
      - ./backend/config.yaml:/app/config.yaml
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "wget -qO- http://localhost:8080/api/v1/health || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 10s
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "80:80"
    depends_on:
      backend:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "wget -qO- http://localhost:80/ || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 3
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  pgdata:
```

Note: The health endpoint `GET /api/v1/health` already exists in `router.go` line 42-44 — it returns `{"code":0,"message":"ok"}`. No backend change needed.

- [ ] **Step 5: Verify backend and frontend still build**

Run: `cd backend && go build ./... && cd ../frontend && npm run build`
Expected: Both succeed.

- [ ] **Step 6: Commit**

```bash
git add docker-compose.yml frontend/nginx.conf backend/.dockerignore frontend/.dockerignore
git commit -m "fix: productionize Docker Compose — health checks, restart, security headers, .dockerignore"
```

---

## Task 4: Swagger API Docs + Error Code Constants

**Files:**
- Modify: `backend/go.mod` (add swaggo/swag dependency)
- Create: `backend/internal/service/error_codes.go`
- Modify: `backend/internal/service/errors.go`
- Modify: `backend/internal/service/auth.go`
- Modify: `backend/internal/service/user.go`
- Modify: `backend/internal/service/subscription.go`
- Modify: `backend/internal/service/asset.go`
- Modify: `backend/internal/service/category.go`
- Modify: `backend/internal/service/reminder.go`
- Modify: `backend/internal/router/router.go`
- Modify: `backend/cmd/server/main.go`
- Create: `backend/docs/` (auto-generated by swag)
- Create: `backend/internal/api/health.go`

### 4.1 Add swaggo/swag dependency

- [ ] **Step 1: Install swaggo packages**

```bash
cd backend
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### 4.2 Create error code constants

- [ ] **Step 2: Create `backend/internal/service/error_codes.go`**

```go
package service

// Error code constants — replace all hardcoded error codes with these.
const (
	ErrCodeInvalidParams       = 40001
	ErrCodeInvalidReminderID   = 40002
	ErrCodeIncorrectPassword   = 40003
	ErrCodeUnauthorized        = 40100
	ErrCodeInvalidCredentials  = 40101
	ErrCodeForbidden           = 40301
	ErrCodeNotFound            = 40400
	ErrCodeDuplicateUsername   = 40901
	ErrCodeDuplicateEmail      = 40902
	ErrCodeDuplicateCategory   = 40903
)
```

### 4.3 Update service files to use constants

- [ ] **Step 3: Update `service/auth.go`** — replace literal codes

In `auth.go`, change:
- `NewServiceError(409, 40901, ...)` → `NewServiceError(409, ErrCodeDuplicateUsername, ...)`
- `NewServiceError(409, 40902, ...)` → `NewServiceError(409, ErrCodeDuplicateEmail, ...)`
- `NewServiceError(401, 40101, ...)` → `NewServiceError(401, ErrCodeInvalidCredentials, ...)`

- [ ] **Step 4: Update `service/user.go`** — replace literal codes

Change:
- `NewServiceError(404, 40400, ...)` → `NewServiceError(404, ErrCodeNotFound, ...)`
- `NewServiceError(400, 40002, ...)` → `NewServiceError(400, ErrCodeIncorrectPassword, ...)`

- [ ] **Step 5: Update `service/subscription.go`** — replace literal codes

Change all occurrences:
- `NewServiceError(404, 40400, ...)` → `NewServiceError(404, ErrCodeNotFound, ...)`
- `NewServiceError(403, 40301, ...)` → `NewServiceError(403, ErrCodeForbidden, ...)`

- [ ] **Step 6: Update `service/asset.go`** — replace literal codes

Same pattern:
- `NewServiceError(404, 40400, ...)` → `NewServiceError(404, ErrCodeNotFound, ...)`
- `NewServiceError(403, 40301, ...)` → `NewServiceError(403, ErrCodeForbidden, ...)`

- [ ] **Step 7: Update `service/category.go`** — replace literal codes

Change:
- `NewServiceError(404, 40400, ...)` → `NewServiceError(404, ErrCodeNotFound, ...)`
- `NewServiceError(403, 40301, ...)` → `NewServiceError(403, ErrCodeForbidden, ...)`
- `NewServiceError(409, 40901, ...)` → `NewServiceError(409, ErrCodeDuplicateCategory, ...)`

- [ ] **Step 8: Update `service/reminder.go`** — replace literal codes

Change:
- `NewServiceError(400, 40001, ...)` → `NewServiceError(400, ErrCodeInvalidReminderID, ...)`

### 4.4 Add Swagger annotations and wire up

- [ ] **Step 9: Add Swagger annotations to all API handlers**

Add `swag` import and annotations to each handler file. Example for `api/auth.go`:

```go
package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

// Register godoc
// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "Registration data"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := "参数校验失败"
		if strings.Contains(err.Error(), "Username") {
			msg = "用户名需 3-50 个字符"
		} else if strings.Contains(err.Error(), "Email") {
			msg = "邮箱格式不正确"
		} else if strings.Contains(err.Error(), "Password") {
			msg = "密码需 6-50 个字符"
		}
		response.Fail(c, 400, service.ErrCodeInvalidParams, msg)
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

// Login godoc
// @Summary User login
// @Description Authenticate and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, service.ErrCodeInvalidParams, "用户名和密码不能为空")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

// GetCurrentUser godoc
// @Summary Get current user info
// @Description Get the authenticated user's profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := h.userService.GetCurrentUser(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}
```

Apply similar Swagger annotations to all handlers in:
- `api/user.go` — UpdateProfile, UpdatePassword
- `api/subscription.go` — List, GetByID, Create, Update, Delete
- `api/asset.go` — List, GetByID, Create, Update, Delete
- `api/category.go` — List, GetByID, Create, Update, Delete
- `api/reminder.go` — List, MarkRead, MarkAllRead, Delete
- `api/dashboard.go` — GetDashboard

Each annotation follows the pattern:
```go
// MethodName godoc
// @Summary <short summary>
// @Description <longer description>
// @Tags <tag name>
// @Accept json
// @Produce json
// @Security BearerAuth  (for authenticated endpoints)
// @Param <param> <location> <type> <description>
// @Success 200 {object} response.Response{data=<ResponseType>}
// @Failure <code> {object} response.Response
// @Router <path> [<method>]
```

Also update all `response.Fail(c, 400, 40001, ...)` calls in API handlers to use `service.ErrCodeInvalidParams` instead of `40001`.

- [ ] **Step 10: Add main Swagger annotation to `cmd/server/main.go`**

```go
// Package main StackBill API server
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
```

- [ ] **Step 11: Wire Swagger UI in `router.go`**

```go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/api"
	"github.com/kingqaquuu/stackbill/internal/middleware"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"github.com/kingqaquuu/stackbill/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/kingqaquuu/stackbill/docs"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, jwtSecret string, jwtExpireHours int) {
	r.Use(middleware.CORSMiddleware())

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	dashboardService := service.NewDashboardService(subscriptionRepo, assetRepo, reminderRepo, categoryRepo, subscriptionService)

	// Handlers
	authHandler := api.NewAuthHandler(authService, userService)
	userHandler := api.NewUserHandler(userService)
	categoryHandler := api.NewCategoryHandler(categoryService)
	subscriptionHandler := api.NewSubscriptionHandler(subscriptionService)
	assetHandler := api.NewAssetHandler(assetService)
	reminderHandler := api.NewReminderHandler(reminderService)
	dashboardHandler := api.NewDashboardHandler(dashboardService)

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
		authorized.DELETE("/reminders/:id", reminderHandler.Delete)
	}
}
```

- [ ] **Step 12: Add RegisterRequest/LoginRequest/UserResponse DTOs with json names for Swagger**

The existing DTOs in `dto/user.go` need to be extended with the request/response types for auth (currently inlined). Add to `dto/user.go`:

```go
package dto

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=500"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}
```

- [ ] **Step 13: Generate Swagger docs**

```bash
cd backend
~/go/bin/swag init -g cmd/server/main.go -o docs/
```

- [ ] **Step 14: Update `.dockerignore` to exclude generated docs from builder**

Add to `backend/.dockerignore`:
```
docs/
```

Actually the docs need to be in the Docker image for Swagger UI. So do NOT exclude them. Keep `.dockerignore` as-is.

- [ ] **Step 15: Verify backend builds with Swagger**

Run: `cd backend && go build ./...`
Expected: Success.

- [ ] **Step 16: Commit**

```bash
git add -A
git commit -m "feat: add Swagger API docs and error code constants"
```

---

## Task 5: Backend Tests

**Files:**
- Create: `backend/internal/testutil/testutil.go`
- Create: `backend/internal/service/auth_test.go`
- Create: `backend/internal/service/subscription_test.go`
- Create: `backend/internal/service/asset_test.go`
- Create: `backend/internal/service/category_test.go`
- Create: `backend/internal/service/reminder_test.go`
- Create: `backend/internal/service/dashboard_test.go`
- Create: `backend/internal/api/middleware_test.go`

### 5.1 Add SQLite test dependency and test helpers

- [ ] **Step 1: Add SQLite driver dependency**

```bash
cd backend
go get gorm.io/driver/sqlite
```

- [ ] **Step 2: Create `backend/internal/testutil/testutil.go`**

```go
package testutil

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/config"
	"github.com/kingqaquuu/stackbill/internal/middleware"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"github.com/kingqaquuu/stackbill/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const TestJWTSecret = "test-secret-key"
const TestJWTExpire = 72

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
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
	return db
}

type TestServices struct {
	DB                 *gorm.DB
	UserRepo           *repository.UserRepository
	CategoryRepo       *repository.CategoryRepository
	SubscriptionRepo   *repository.SubscriptionRepository
	AssetRepo          *repository.AssetRepository
	ReminderRepo       *repository.ReminderRepository
	AuthService        *service.AuthService
	UserService        *service.UserService
	CategoryService    *service.CategoryService
	SubscriptionService *service.SubscriptionService
	AssetService       *service.AssetService
	ReminderService    *service.ReminderService
	DashboardService   *service.DashboardService
}

func NewTestServices(t *testing.T) *TestServices {
	t.Helper()
	db := SetupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	reminderRepo := repository.NewReminderRepository(db)

	authService := service.NewAuthService(userRepo, categoryRepo, TestJWTSecret, TestJWTExpire)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	assetService := service.NewAssetService(assetRepo)
	reminderService := service.NewReminderService(reminderRepo)
	dashboardService := service.NewDashboardService(subscriptionRepo, assetRepo, reminderRepo, categoryRepo, subscriptionService)

	return &TestServices{
		DB:                 db,
		UserRepo:           userRepo,
		CategoryRepo:       categoryRepo,
		SubscriptionRepo:   subscriptionRepo,
		AssetRepo:          assetRepo,
		ReminderRepo:       reminderRepo,
		AuthService:        authService,
		UserService:        userService,
		CategoryService:    categoryService,
		SubscriptionService: subscriptionService,
		AssetService:       assetService,
		ReminderService:    reminderService,
		DashboardService:   dashboardService,
	}
}

func CreateTestUser(t *testing.T, svc *TestServices, username string) (uint, string) {
	t.Helper()
	resp, err := svc.AuthService.Register(&config.RegisterRequest{
		Username: username,
		Email:    username + "@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	return resp.User.ID, resp.Token
}

// Helper to generate a valid JWT token for testing
func GenerateTestToken(t *testing.T, userID uint, username string) string {
	t.Helper()
	token, err := middleware.GenerateToken(userID, username, TestJWTSecret, TestJWTExpire)
	if err != nil {
		t.Fatalf("generate test token: %v", err)
	}
	return token
}
```

Wait — `CreateTestUser` references `config.RegisterRequest` but that type is in `dto`, not `config`. Let me fix:

```go
package testutil

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/middleware"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"github.com/kingqaquuu/stackbill/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const TestJWTSecret = "test-secret-key"
const TestJWTExpire = 72

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
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
	return db
}

type TestServices struct {
	DB                  *gorm.DB
	UserRepo            *repository.UserRepository
	CategoryRepo        *repository.CategoryRepository
	SubscriptionRepo    *repository.SubscriptionRepository
	AssetRepo           *repository.AssetRepository
	ReminderRepo        *repository.ReminderRepository
	AuthService         *service.AuthService
	UserService         *service.UserService
	CategoryService     *service.CategoryService
	SubscriptionService *service.SubscriptionService
	AssetService        *service.AssetService
	ReminderService     *service.ReminderService
	DashboardService    *service.DashboardService
}

func NewTestServices(t *testing.T) *TestServices {
	t.Helper()
	db := SetupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	reminderRepo := repository.NewReminderRepository(db)

	authService := service.NewAuthService(userRepo, categoryRepo, TestJWTSecret, TestJWTExpire)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	assetService := service.NewAssetService(assetRepo)
	reminderService := service.NewReminderService(reminderRepo)
	dashboardService := service.NewDashboardService(subscriptionRepo, assetRepo, reminderRepo, categoryRepo, subscriptionService)

	return &TestServices{
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

// RegisterTestUser creates a user via the AuthService and returns userID.
// Also returns the default categories created during registration.
func RegisterTestUser(t *testing.T, svc *TestServices, username string) uint {
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

func GenerateTestToken(t *testing.T, userID uint, username string) string {
	t.Helper()
	token, err := middleware.GenerateToken(userID, username, TestJWTSecret, TestJWTExpire)
	if err != nil {
		t.Fatalf("generate test token: %v", err)
	}
	return token
}
```

- [ ] **Step 3: Create `backend/internal/service/auth_test.go`**

```go
package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/testutil"
)

func TestAuthService_Register_Success(t *testing.T) {
	svc := testutil.NewTestServices(t)

	resp, err := svc.AuthService.Register(&dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token")
	}
	if resp.User.Username != "testuser" {
		t.Errorf("username = %q, want %q", resp.User.Username, "testuser")
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("email = %q, want %q", resp.User.Email, "test@example.com")
	}
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	svc := testutil.NewTestServices(t)

	_, _ = svc.AuthService.Register(&dto.RegisterRequest{
		Username: "testuser",
		Email:    "first@example.com",
		Password: "password123",
	})

	_, err := svc.AuthService.Register(&dto.RegisterRequest{
		Username: "testuser",
		Email:    "second@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
	svcErr, ok := err.(*ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if svcErr.Code != ErrCodeDuplicateUsername {
		t.Errorf("code = %d, want %d", svcErr.Code, ErrCodeDuplicateUsername)
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	svc := testutil.NewTestServices(t)

	_, _ = svc.AuthService.Register(&dto.RegisterRequest{
		Username: "user1",
		Email:    "same@example.com",
		Password: "password123",
	})

	_, err := svc.AuthService.Register(&dto.RegisterRequest{
		Username: "user2",
		Email:    "same@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
	svcErr, ok := err.(*ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if svcErr.Code != ErrCodeDuplicateEmail {
		t.Errorf("code = %d, want %d", svcErr.Code, ErrCodeDuplicateEmail)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	svc := testutil.NewTestServices(t)

	_, _ = svc.AuthService.Register(&dto.RegisterRequest{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: "password123",
	})

	resp, err := svc.AuthService.Login(&dto.LoginRequest{
		Username: "loginuser",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token after login")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc := testutil.NewTestServices(t)

	_, _ = svc.AuthService.Register(&dto.RegisterRequest{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: "password123",
	})

	_, err := svc.AuthService.Login(&dto.LoginRequest{
		Username: "loginuser",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
	svcErr, ok := err.(*ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if svcErr.Code != ErrCodeInvalidCredentials {
		t.Errorf("code = %d, want %d", svcErr.Code, ErrCodeInvalidCredentials)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc := testutil.NewTestServices(t)

	_, err := svc.AuthService.Login(&dto.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestAuthService_Register_CreatesDefaultCategories(t *testing.T) {
	svc := testutil.NewTestServices(t)

	resp, _ := svc.AuthService.Register(&dto.RegisterRequest{
		Username: "catuser",
		Email:    "cat@example.com",
		Password: "password123",
	})

	cats, err := svc.CategoryService.List(resp.User.ID, &dto.CategoryListQuery{})
	if err != nil {
		t.Fatalf("list categories: %v", err)
	}
	if len(cats) != 8 {
		t.Errorf("got %d categories, want 8", len(cats))
	}
}
```

- [ ] **Step 4: Create `backend/internal/service/subscription_test.go`**

```go
package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/testutil"
)

func TestSubscriptionService_CreateAndGet(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "subuser")

	startDate := "2026-01-01"
	resp, err := svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name:         "Netflix",
		Amount:       15.99,
		Currency:     "USD",
		BillingCycle: "monthly",
		StartDate:    &startDate,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.Name != "Netflix" {
		t.Errorf("name = %q, want %q", resp.Name, "Netflix")
	}
	if resp.Status != "active" {
		t.Errorf("status = %q, want %q", resp.Status, "active")
	}

	got, err := svc.SubscriptionService.GetByID(userID, resp.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Name != "Netflix" {
		t.Errorf("got name = %q, want %q", got.Name, "Netflix")
	}
}

func TestSubscriptionService_Update(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "subuser2")

	created, _ := svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name:         "Spotify",
		Amount:       9.99,
		Currency:     "USD",
		BillingCycle: "monthly",
	})

	updated, err := svc.SubscriptionService.Update(userID, created.ID, &dto.UpdateSubscriptionRequest{
		Name:   "Spotify Premium",
		Amount: 14.99,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Spotify Premium" {
		t.Errorf("name = %q, want %q", updated.Name, "Spotify Premium")
	}
	if updated.Amount != 14.99 {
		t.Errorf("amount = %f, want %f", updated.Amount, 14.99)
	}
}

func TestSubscriptionService_Delete(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "subuser3")

	created, _ := svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name:         "ToDelete",
		Amount:       5.00,
		Currency:     "USD",
		BillingCycle: "yearly",
	})

	err := svc.SubscriptionService.Delete(userID, created.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = svc.SubscriptionService.GetByID(userID, created.ID)
	if err == nil {
		t.Error("expected error getting deleted subscription")
	}
}

func TestSubscriptionService_UserIsolation(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userA := testutil.RegisterTestUser(t, svc, "userA")
	userB := testutil.RegisterTestUser(t, svc, "userB")

	created, _ := svc.SubscriptionService.Create(userA, &dto.CreateSubscriptionRequest{
		Name:         "Private Sub",
		Amount:       10.00,
		Currency:     "USD",
		BillingCycle: "monthly",
	})

	_, err := svc.SubscriptionService.GetByID(userB, created.ID)
	if err == nil {
		t.Error("user B should not access user A's subscription")
	}

	err = svc.SubscriptionService.Delete(userB, created.ID)
	if err == nil {
		t.Error("user B should not delete user A's subscription")
	}
}

func TestSubscriptionService_CalculateMonthlyExpense(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "expenseuser")

	svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name:         "Monthly Sub",
		Amount:       10.00,
		Currency:     "USD",
		BillingCycle: "monthly",
	})
	svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name:         "Yearly Sub",
		Amount:       120.00,
		Currency:     "USD",
		BillingCycle: "yearly",
	})

	monthly, err := svc.SubscriptionService.CalculateMonthlyExpense(userID)
	if err != nil {
		t.Fatalf("monthly expense: %v", err)
	}
	// 10.00 (monthly) + 10.00 (120/12 yearly) = 20.00
	if monthly != 20.00 {
		t.Errorf("monthly expense = %f, want %f", monthly, 20.00)
	}
}
```

- [ ] **Step 5: Create `backend/internal/service/asset_test.go`**

```go
package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/testutil"
)

func TestAssetService_CreateAndGet(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "assetuser")

	expDate := "2027-01-01"
	resp, err := svc.AssetService.Create(userID, &dto.CreateAssetRequest{
		Name:      "example.com",
		AssetType: "domain",
		ExpireDate: &expDate,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.Name != "example.com" {
		t.Errorf("name = %q, want %q", resp.Name, "example.com")
	}

	got, err := svc.AssetService.GetByID(userID, resp.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "example.com" {
		t.Errorf("got name = %q, want %q", got.Name, "example.com")
	}
}

func TestAssetService_Update(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "assetuser2")

	created, _ := svc.AssetService.Create(userID, &dto.CreateAssetRequest{
		Name:      "myserver",
		AssetType: "server",
	})

	updated, err := svc.AssetService.Update(userID, created.ID, &dto.UpdateAssetRequest{
		Name:      "myserver-v2",
		AssetType: "server",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "myserver-v2" {
		t.Errorf("name = %q, want %q", updated.Name, "myserver-v2")
	}
}

func TestAssetService_Delete(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "assetuser3")

	created, _ := svc.AssetService.Create(userID, &dto.CreateAssetRequest{
		Name:      "to-delete",
		AssetType: "ssl_certificate",
	})

	err := svc.AssetService.Delete(userID, created.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = svc.AssetService.GetByID(userID, created.ID)
	if err == nil {
		t.Error("expected error getting deleted asset")
	}
}

func TestAssetService_UserIsolation(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userA := testutil.RegisterTestUser(t, svc, "userA_asset")
	userB := testutil.RegisterTestUser(t, svc, "userB_asset")

	created, _ := svc.AssetService.Create(userA, &dto.CreateAssetRequest{
		Name:      "private-domain",
		AssetType: "domain",
	})

	_, err := svc.AssetService.GetByID(userB, created.ID)
	if err == nil {
		t.Error("user B should not access user A's asset")
	}

	err = svc.AssetService.Delete(userB, created.ID)
	if err == nil {
		t.Error("user B should not delete user A's asset")
	}
}
```

- [ ] **Step 6: Create `backend/internal/service/category_test.go`**

```go
package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/testutil"
)

func TestCategoryService_CreateAndGet(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "catuser")

	resp, err := svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name:  "Custom Category",
		Type:  "subscription",
		Color: "#ff0000",
		Icon:  "star",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.Name != "Custom Category" {
		t.Errorf("name = %q, want %q", resp.Name, "Custom Category")
	}

	got, err := svc.CategoryService.GetByID(userID, resp.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Name != "Custom Category" {
		t.Errorf("got name = %q, want %q", got.Name, "Custom Category")
	}
}

func TestCategoryService_DuplicateName(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "catuser2")

	_, _ = svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name: "SameName", Type: "subscription",
	})

	_, err := svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name: "SameName", Type: "subscription",
	})
	if err == nil {
		t.Fatal("expected error for duplicate category name")
	}
}

func TestCategoryService_Update(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "catuser3")

	created, _ := svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name: "Original", Type: "asset", Color: "#000000",
	})

	updated, err := svc.CategoryService.Update(userID, created.ID, &dto.UpdateCategoryRequest{
		Name: "Updated", Type: "asset", Color: "#ffffff",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("name = %q, want %q", updated.Name, "Updated")
	}
}

func TestCategoryService_Delete(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "catuser4")

	created, _ := svc.CategoryService.Create(userID, &dto.CreateCategoryRequest{
		Name: "ToDelete", Type: "subscription",
	})

	err := svc.CategoryService.Delete(userID, created.ID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = svc.CategoryService.GetByID(userID, created.ID)
	if err == nil {
		t.Error("expected error getting deleted category")
	}
}

func TestCategoryService_UserIsolation(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userA := testutil.RegisterTestUser(t, svc, "catuserA")
	userB := testutil.RegisterTestUser(t, svc, "catuserB")

	created, _ := svc.CategoryService.Create(userA, &dto.CreateCategoryRequest{
		Name: "PrivateCat", Type: "subscription",
	})

	_, err := svc.CategoryService.GetByID(userB, created.ID)
	if err == nil {
		t.Error("user B should not access user A's category")
	}
}
```

- [ ] **Step 7: Create `backend/internal/service/reminder_test.go`**

```go
package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/testutil"
)

func TestReminderService_ListEmpty(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "remuser")

	result, err := svc.ReminderService.List(userID, &dto.ReminderListQuery{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("total = %d, want 0", result.Total)
	}
}

func TestReminderService_ListWithRenewals(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "remuser2")

	// Create a subscription with next_payment_date within 7 days
	// Since calculateNextPayment sets it based on start_date, we need a recent start date
	startDate := "2026-06-01"
	svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name:         "Renewing Soon",
		Amount:       10.00,
		Currency:     "USD",
		BillingCycle: "weekly",
		StartDate:    &startDate,
	})

	result, err := svc.ReminderService.List(userID, &dto.ReminderListQuery{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total < 1 {
		t.Error("expected at least 1 reminder")
	}
}

func TestReminderService_MarkRead(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "remuser3")

	startDate := "2026-06-01"
	sub, _ := svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name:         "ReadTest",
		Amount:       5.00,
		Currency:     "USD",
		BillingCycle: "weekly",
		StartDate:    &startDate,
	})

	// Mark the renewal reminder as read (ID = sub.ID + reminderOffsetRenewal = sub.ID)
	err := svc.ReminderService.MarkRead(userID, sub.ID)
	if err != nil {
		t.Fatalf("mark read: %v", err)
	}
}
```

- [ ] **Step 8: Create `backend/internal/service/dashboard_test.go`**

```go
package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/testutil"
)

func TestDashboardService_GetDashboard(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userID := testutil.RegisterTestUser(t, svc, "dashuser")

	svc.SubscriptionService.Create(userID, &dto.CreateSubscriptionRequest{
		Name:         "Test Sub",
		Amount:       10.00,
		Currency:     "USD",
		BillingCycle: "monthly",
	})
	svc.AssetService.Create(userID, &dto.CreateAssetRequest{
		Name:      "Test Asset",
		AssetType: "domain",
	})

	dash, err := svc.DashboardService.GetDashboard(userID)
	if err != nil {
		t.Fatalf("get dashboard: %v", err)
	}
	if dash.SubscriptionCount != 1 {
		t.Errorf("subscription count = %d, want 1", dash.SubscriptionCount)
	}
	if dash.AssetCount != 1 {
		t.Errorf("asset count = %d, want 1", dash.AssetCount)
	}
	if dash.MonthlyExpense != 10.00 {
		t.Errorf("monthly expense = %f, want 10.00", dash.MonthlyExpense)
	}
}

func TestDashboardService_UserIsolation(t *testing.T) {
	svc := testutil.NewTestServices(t)
	userA := testutil.RegisterTestUser(t, svc, "dashA")
	userB := testutil.RegisterTestUser(t, svc, "dashB")

	svc.SubscriptionService.Create(userA, &dto.CreateSubscriptionRequest{
		Name:         "Private",
		Amount:       100.00,
		Currency:     "USD",
		BillingCycle: "monthly",
	})

	dashB, _ := svc.DashboardService.GetDashboard(userB)
	if dashB.SubscriptionCount != 0 {
		t.Errorf("user B subscription count = %d, want 0", dashB.SubscriptionCount)
	}
	if dashB.MonthlyExpense != 0 {
		t.Errorf("user B monthly expense = %f, want 0", dashB.MonthlyExpense)
	}
}
```

- [ ] **Step 9: Create `backend/internal/api/middleware_test.go`**

```go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/middleware"
	"github.com/kingqaquuu/stackbill/internal/testutil"
)

func TestJWTAuth_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.JWTAuth(testutil.TestJWTSecret))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.JWTAuth(testutil.TestJWTSecret))
	r.GET("/test", func(c *gin.Context) {
		userID := c.GetUint("user_id")
		c.JSON(200, gin.H{"user_id": userID})
	})

	token := testutil.GenerateTestToken(t, 42, "testuser")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("status = %d, want 200", w.Code)
	}
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.JWTAuth(testutil.TestJWTSecret))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// Generate a token that expires immediately (0 hours)
	token, _ := middleware.GenerateToken(1, "user", testutil.TestJWTSecret, 0)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.JWTAuth(testutil.TestJWTSecret))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("status = %d, want 401", w.Code)
	}
}
```

- [ ] **Step 10: Run all tests**

```bash
cd backend && go test ./... -v -count=1
```

Expected: All tests pass. Some tests may need adjustment because SQLite handles certain queries differently than PostgreSQL (e.g., raw SQL in `GetCategoryExpense` uses PostgreSQL-specific functions). If `GetCategoryExpense` fails, note it and adjust the test to skip that specific assertion.

- [ ] **Step 11: Commit**

```bash
git add -A
git commit -m "test: add backend unit tests — auth, subscriptions, assets, categories, reminders, dashboard, middleware"
```

---

## Task 6: GitHub Actions CI

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Create `.github/workflows/ci.yml`**

```yaml
name: CI

on:
  push:
    branches: [master]
  pull_request:
    branches: [master]

jobs:
  backend:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: backend
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.25"
          cache-dependency-path: backend/go.sum

      - name: Install swag
        run: go install github.com/swaggo/swag/cmd/swag@latest

      - name: Generate Swagger docs
        run: ~/go/bin/swag init -g cmd/server/main.go -o docs/

      - name: Build
        run: go build ./...

      - name: Test
        run: go test ./... -race -count=1

  frontend:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - uses: actions/checkout@v4

      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: "22"
          cache: npm
          cache-dependency-path: frontend/package-lock.json

      - name: Install dependencies
        run: npm ci

      - name: Build
        run: npm run build
```

- [ ] **Step 2: Commit**

```bash
git add .github/
git commit -m "ci: add GitHub Actions workflow for backend and frontend"
```

---

## Task 7: Update README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README.md**

Update the README to reflect all changes:
- Remove `backend/migrations/` from the project structure
- Update "数据库迁移" section to only mention AutoMigrate
- Add Swagger UI info: accessible at `http://localhost/swagger/index.html` after starting
- Add health endpoint info
- Add CI badge placeholder
- Note that `npm run build` includes type checking

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: update README for V1 — Swagger, health endpoint, AutoMigrate only"
```

---

## Acceptance Criteria Checklist

- [ ] `go build ./...` passes (backend)
- [ ] `go test ./... -race` passes (all tests green)
- [ ] `npm run build` passes (frontend type-check + build)
- [ ] `docker compose up` starts all services
- [ ] Swagger UI accessible at `http://localhost/swagger/index.html`
- [ ] No SQL migration files remain
- [ ] No known frontend data-access bugs
- [ ] No duplicate i18n keys
- [ ] GitHub Actions CI passes on master
