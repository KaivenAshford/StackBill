# User System Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the complete user authentication and profile management system — registration, login, JWT auth, profile editing, and password change — with strict three-layer architecture (API → Service → Repository).

**Architecture:** Each feature flows through API handler (HTTP binding/response) → Service (business logic) → Repository (GORM queries). DTOs isolate models from HTTP layer. Repos are concrete structs, not interfaces. Default categories are created on registration.

**Tech Stack:** Go, Gin, GORM, PostgreSQL, bcrypt, JWT; Vue 3, Naive UI, Pinia, Axios

---

### Task 1: Install bcrypt dependency

**Files:**
- Modify: `backend/go.mod`

- [ ] **Step 1: Add bcrypt package**

Run:
```bash
cd /home/kingqaquuu/StackBill/backend
go get golang.org/x/crypto/bcrypt
go mod tidy
```

Expected: `go.mod` updated with `golang.org/x/crypto` dependency, `go build ./...` succeeds.

---

### Task 2: Add user DTOs

**Files:**
- Create: `backend/internal/dto/user.go`

- [ ] **Step 1: Create dto/user.go**

Create `backend/internal/dto/user.go`:

```go
package dto

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=500"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}
```

- [ ] **Step 2: Verify build**

Run: `cd /home/kingqaquuu/StackBill/backend && go build ./...`
Expected: compiles without error.

---

### Task 3: Create repositories

**Files:**
- Create: `backend/internal/repository/user.go`
- Create: `backend/internal/repository/category.go`

- [ ] **Step 1: Create repository/user.go**

Create `backend/internal/repository/user.go`:

```go
package repository

import (
	"github.com/kingqaquuu/stackbill/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
```

- [ ] **Step 2: Create repository/category.go**

Create `backend/internal/repository/category.go`:

```go
package repository

import (
	"github.com/kingqaquuu/stackbill/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) BatchCreate(categories []model.Category) error {
	return r.db.Create(&categories).Error
}
```

- [ ] **Step 3: Verify build**

Run: `cd /home/kingqaquuu/StackBill/backend && go build ./...`
Expected: compiles without error.

---

### Task 4: Create services

**Files:**
- Create: `backend/internal/service/auth.go`
- Create: `backend/internal/service/user.go`

- [ ] **Step 1: Create service/auth.go**

Create `backend/internal/service/auth.go`:

```go
package service

import (
	"errors"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/middleware"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	categoryRepo *repository.CategoryRepository
	jwtSecret    string
	jwtExpire    int
}

func NewAuthService(userRepo *repository.UserRepository, categoryRepo *repository.CategoryRepository, jwtSecret string, jwtExpire int) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
		jwtSecret:    jwtSecret,
		jwtExpire:    jwtExpire,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, NewServiceError(409, 40901, "username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, NewServiceError(409, 40902, "email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	s.initDefaultCategories(user.ID)

	token, err := middleware.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email},
	}, nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewServiceError(401, 40101, "invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, NewServiceError(401, 40101, "invalid credentials")
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email, Nickname: user.Nickname, Avatar: user.Avatar},
	}, nil
}

func (s *AuthService) initDefaultCategories(userID uint) {
	categories := []model.Category{
		{UserID: userID, Name: "AI 工具", Type: "subscription", Color: "#7c3aed", Icon: "robot", SortOrder: 1},
		{UserID: userID, Name: "开发工具", Type: "subscription", Color: "#2563eb", Icon: "code", SortOrder: 2},
		{UserID: userID, Name: "云服务", Type: "subscription", Color: "#0891b2", Icon: "cloud", SortOrder: 3},
		{UserID: userID, Name: "域名", Type: "subscription", Color: "#059669", Icon: "globe", SortOrder: 4},
		{UserID: userID, Name: "服务器", Type: "subscription", Color: "#d97706", Icon: "server", SortOrder: 5},
		{UserID: userID, Name: "娱乐", Type: "subscription", Color: "#dc2626", Icon: "game-controller", SortOrder: 6},
		{UserID: userID, Name: "办公", Type: "subscription", Color: "#4f46e5", Icon: "briefcase", SortOrder: 7},
		{UserID: userID, Name: "其他", Type: "subscription", Color: "#6b7280", Icon: "ellipsis", SortOrder: 8},
	}
	_ = s.categoryRepo.BatchCreate(categories)
}
```

- [ ] **Step 2: Create service/errors.go**

Create `backend/internal/service/errors.go`:

```go
package service

type ServiceError struct {
	HTTPCode int
	Code     int
	Message  string
}

func NewServiceError(httpCode, code int, message string) *ServiceError {
	return &ServiceError{HTTPCode: httpCode, Code: code, Message: message}
}

func (e *ServiceError) Error() string {
	return e.Message
}
```

- [ ] **Step 3: Create service/user.go**

Create `backend/internal/service/user.go`:

```go
package service

import (
	"errors"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetCurrentUser(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

func (s *UserService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	user.Nickname = req.Nickname
	user.Avatar = req.Avatar

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

func (s *UserService) UpdatePassword(userID uint, req *dto.UpdatePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewServiceError(404, 40400, "user not found")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return NewServiceError(400, 40002, "incorrect old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}
```

- [ ] **Step 4: Verify build**

Run: `cd /home/kingqaquuu/StackBill/backend && go build ./...`
Expected: compiles without error.

---

### Task 5: Create API handlers

**Files:**
- Create: `backend/internal/api/auth.go`
- Create: `backend/internal/api/user.go`

- [ ] **Step 1: Create api/auth.go**

Create `backend/internal/api/auth.go`:

```go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.OK(c, resp)
}
```

Note: `GetCurrentUser` is placed on `AuthHandler` because it's under the `/auth/me` route. It reuses the same handler pattern and calls through to `UserService` internally. However, to keep layer separation clean, `AuthHandler.GetCurrentUser` delegates to `UserService`. We need to add a `userService` field to `AuthHandler`.

Revised `api/auth.go` — use this version instead:

```go
package api

import (
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

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

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

- [ ] **Step 2: Create api/user.go**

Create `backend/internal/api/user.go`:

```go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/service"
	"github.com/kingqaquuu/stackbill/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}

	userID := c.GetUint("user_id")
	resp, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, resp)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, 40001, "invalid parameters")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.userService.UpdatePassword(userID, &req); err != nil {
		handleServiceError(c, err)
		return
	}

	response.OK(c, nil)
}
```

- [ ] **Step 3: Create api/helpers.go**

Create `backend/internal/api/helpers.go`:

```go
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
```

- [ ] **Step 4: Verify build**

Run: `cd /home/kingqaquuu/StackBill/backend && go build ./...`
Expected: compiles without error.

---

### Task 6: Update router and main.go

**Files:**
- Modify: `backend/internal/router/router.go` (full rewrite)
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Rewrite router.go**

Replace the entire content of `backend/internal/router/router.go` with:

```go
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
```

- [ ] **Step 2: Update main.go**

Replace the entire content of `backend/cmd/server/main.go` with:

```go
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

- [ ] **Step 3: Verify build**

Run: `cd /home/kingqaquuu/StackBill/backend && go build ./...`
Expected: compiles without error.

- [ ] **Step 4: Commit backend**

```bash
cd /home/kingqaquuu/StackBill
git add backend/
git commit -m "feat: implement user system - registration, login, JWT auth, profile management"
```

---

### Task 7: Backend smoke test with curl

This task requires a running PostgreSQL instance. If Docker is available:

- [ ] **Step 1: Start PostgreSQL**

```bash
cd /home/kingqaquuu/StackBill
docker compose up -d postgres
```

Wait for it to be healthy: `docker compose ps`

- [ ] **Step 2: Copy config and start backend**

```bash
cd /home/kingqaquuu/StackBill/backend
cp config.example.yaml config.yaml
# Edit config.yaml: set host to "localhost" instead of "postgres"
sed -i 's/host: postgres/host: localhost/' config.yaml
JWT_SECRET=test-secret-key go run ./cmd/server &
```

Expected: `server starting on :8080`

- [ ] **Step 3: Test registration**

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"123456"}' | jq .
```

Expected:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "...",
    "user": { "id": 1, "username": "testuser", "email": "test@example.com", "nickname": "", "avatar": "" }
  }
}
```

- [ ] **Step 4: Test duplicate registration**

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test2@example.com","password":"123456"}' | jq .
```

Expected: `{"code":40901,"message":"username already exists"}`

- [ ] **Step 5: Test login**

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}' | jq .
```

Expected: Same response as registration with a new token.

- [ ] **Step 6: Test /auth/me**

```bash
TOKEN="<token from step 5>"
curl -s http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Expected: `{"code":0,"data":{"id":1,"username":"testuser",...}}`

- [ ] **Step 7: Test update profile**

```bash
curl -s -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nickname":"Test User"}' | jq .
```

Expected: `{"code":0,"data":{"id":1,"nickname":"Test User",...}}`

- [ ] **Step 8: Test change password**

```bash
curl -s -X PUT http://localhost:8080/api/v1/users/password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"123456","new_password":"654321"}' | jq .
```

Expected: `{"code":0,"message":"success"}`

- [ ] **Step 9: Stop backend and PostgreSQL**

```bash
kill %1  # stop the background go server
docker compose -f /home/kingqaquuu/StackBill/docker-compose.yml down
```

---

### Task 8: Update frontend API layer and store

**Files:**
- Modify: `frontend/src/api/auth.ts`
- Modify: `frontend/src/stores/user.ts`

- [ ] **Step 1: Update api/auth.ts**

Replace the entire content of `frontend/src/api/auth.ts` with:

```ts
import request from '@/utils/request'
import type { User } from '@/types'

export function login(username: string, password: string) {
  return request.post<{ token: string; user: User }>('/auth/login', { username, password })
}

export function register(username: string, email: string, password: string) {
  return request.post<{ token: string; user: User }>('/auth/register', { username, email, password })
}

export function getMe() {
  return request.get<User>('/auth/me')
}

export function updateProfile(data: { nickname?: string; avatar?: string }) {
  return request.put<User>('/users/profile', data)
}

export function updatePassword(data: { old_password: string; new_password: string }) {
  return request.put<null>('/users/password', data)
}
```

- [ ] **Step 2: Update stores/user.ts**

Replace the entire content of `frontend/src/stores/user.ts` with:

```ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types'
import { getMe } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref<string>(localStorage.getItem('token') || '')

  function setToken(t: string) {
    token.value = t
    localStorage.setItem('token', t)
  }

  function setUser(u: User) {
    user.value = u
  }

  function logout() {
    user.value = null
    token.value = ''
    localStorage.removeItem('token')
  }

  function isLoggedIn() {
    return !!token.value
  }

  async function fetchUser() {
    if (!token.value) return
    try {
      const res = await getMe()
      user.value = res.data
    } catch {
      logout()
    }
  }

  return { user, token, setToken, setUser, logout, isLoggedIn, fetchUser }
})
```

- [ ] **Step 3: Verify build**

Run: `cd /home/kingqaquuu/StackBill/frontend && npm run build 2>&1 | tail -5`
Expected: `✓ built in ...ms`

---

### Task 9: Implement Settings page

**Files:**
- Modify: `frontend/src/views/settings/Index.vue`
- Modify: `frontend/src/locales/zh-CN.ts`
- Modify: `frontend/src/locales/en-US.ts`

- [ ] **Step 1: Add i18n keys to zh-CN.ts**

In `frontend/src/locales/zh-CN.ts`, add a `settings` section inside the default export object, after the `asset` key:

```ts
  settings: {
    profile: '个人资料',
    changePassword: '修改密码',
    nickname: '昵称',
    avatar: '头像',
    oldPassword: '旧密码',
    newPassword: '新密码',
    confirmPassword: '确认密码',
    passwordMismatch: '两次密码不一致',
  },
```

- [ ] **Step 2: Add i18n keys to en-US.ts**

In `frontend/src/locales/en-US.ts`, add the same section:

```ts
  settings: {
    profile: 'Profile',
    changePassword: 'Change Password',
    nickname: 'Nickname',
    avatar: 'Avatar',
    oldPassword: 'Old Password',
    newPassword: 'New Password',
    confirmPassword: 'Confirm Password',
    passwordMismatch: 'Passwords do not match',
  },
```

- [ ] **Step 3: Rewrite settings/Index.vue**

Replace the entire content of `frontend/src/views/settings/Index.vue` with:

```vue
<template>
  <div class="settings-page">
    <n-card :title="t('settings.profile')" style="margin-bottom: 24px;">
      <n-form :model="profileForm" @submit.prevent="handleUpdateProfile">
        <n-form-item :label="t('settings.nickname')" path="nickname">
          <n-input v-model:value="profileForm.nickname" :placeholder="t('settings.nickname')" />
        </n-form-item>
        <n-button type="primary" :loading="profileLoading" attr-type="submit">
          {{ t('common.save') }}
        </n-button>
      </n-form>
    </n-card>

    <n-card :title="t('settings.changePassword')">
      <n-form :model="passwordForm" @submit.prevent="handleChangePassword">
        <n-form-item :label="t('settings.oldPassword')" path="old_password">
          <n-input v-model:value="passwordForm.old_password" type="password" :placeholder="t('settings.oldPassword')" />
        </n-form-item>
        <n-form-item :label="t('settings.newPassword')" path="new_password">
          <n-input v-model:value="passwordForm.new_password" type="password" :placeholder="t('settings.newPassword')" />
        </n-form-item>
        <n-form-item :label="t('settings.confirmPassword')" path="confirm_password">
          <n-input v-model:value="passwordForm.confirm_password" type="password" :placeholder="t('settings.confirmPassword')" />
        </n-form-item>
        <n-button type="primary" :loading="passwordLoading" attr-type="submit">
          {{ t('common.save') }}
        </n-button>
      </n-form>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { useUserStore } from '@/stores/user'
import { updateProfile, updatePassword } from '@/api/auth'

const { t } = useI18n()
const message = useMessage()
const store = useUserStore()

const profileLoading = ref(false)
const passwordLoading = ref(false)

const profileForm = reactive({ nickname: '', avatar: '' })
const passwordForm = reactive({ old_password: '', new_password: '', confirm_password: '' })

onMounted(() => {
  if (store.user) {
    profileForm.nickname = store.user.nickname || ''
    profileForm.avatar = store.user.avatar || ''
  }
})

async function handleUpdateProfile() {
  profileLoading.value = true
  try {
    const res = await updateProfile({ nickname: profileForm.nickname, avatar: profileForm.avatar })
    store.setUser(res.data)
    message.success(t('common.success'))
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    profileLoading.value = false
  }
}

async function handleChangePassword() {
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    message.error(t('settings.passwordMismatch'))
    return
  }

  passwordLoading.value = true
  try {
    await updatePassword({ old_password: passwordForm.old_password, new_password: passwordForm.new_password })
    message.success(t('common.success'))
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    passwordLoading.value = false
  }
}
</script>

<style scoped>
.settings-page {
  max-width: 600px;
}
</style>
```

- [ ] **Step 4: Verify build**

Run: `cd /home/kingqaquuu/StackBill/frontend && npm run build 2>&1 | tail -5`
Expected: `✓ built in ...ms`

---

### Task 10: Final commit

- [ ] **Step 1: Stage and commit all changes**

```bash
cd /home/kingqaquuu/StackBill
git add frontend/src/api/auth.ts frontend/src/stores/user.ts frontend/src/views/settings/Index.vue frontend/src/locales/zh-CN.ts frontend/src/locales/en-US.ts
git commit -m "feat: connect frontend to user system API - profile settings page, auth store"
```
