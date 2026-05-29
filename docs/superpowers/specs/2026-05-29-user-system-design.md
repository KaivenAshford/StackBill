# StackBill Phase 2: User System Design

## Overview

Implement the complete user authentication and profile management system for StackBill. This includes registration, login, JWT authentication, user info retrieval, profile editing, and password change functionality.

**Architecture:** Strict three-layer separation (API → Service → Repository) with DTO isolation. This pattern will be reused for all subsequent modules.

**Password hashing:** bcrypt (Go `golang.org/x/crypto/bcrypt`, cost factor 10).

## API Design

### Routes

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | Login | No |
| GET | `/api/v1/auth/me` | Get current user info | Yes |
| PUT | `/api/v1/users/profile` | Update nickname/avatar | Yes |
| PUT | `/api/v1/users/password` | Change password (requires old password) | Yes |

### Request/Response

**POST /auth/register**
```json
Request:  { "username": "test", "email": "test@example.com", "password": "123456" }
Response: { "code": 0, "data": { "token": "jwt-token", "user": { "id": 1, "username": "test", "email": "test@example.com", "nickname": "", "avatar": "" } } }
```

**POST /auth/login**
```json
Request:  { "username": "test", "password": "123456" }
Response: { "code": 0, "data": { "token": "jwt-token", "user": { "id": 1, "username": "test", "email": "test@example.com", "nickname": "", "avatar": "" } } }
```

**GET /auth/me**
```json
Response: { "code": 0, "data": { "id": 1, "username": "test", "email": "test@example.com", "nickname": "Test", "avatar": "" } }
```

**PUT /users/profile**
```json
Request:  { "nickname": "Test User", "avatar": "" }
Response: { "code": 0, "data": { "id": 1, "username": "test", "email": "test@example.com", "nickname": "Test User", "avatar": "" } }
```

**PUT /users/password**
```json
Request:  { "old_password": "123456", "new_password": "654321" }
Response: { "code": 0, "data": null }
```

## Error Codes

| Scenario | HTTP | Code | Message |
|----------|------|------|---------|
| Validation failure | 400 | 40001 | invalid parameters |
| Username already exists | 409 | 40901 | username already exists |
| Email already exists | 409 | 40902 | email already exists |
| Wrong username/password | 401 | 40101 | invalid credentials |
| Wrong old password | 400 | 40002 | incorrect old password |
| Not logged in | 401 | 40100 | unauthorized |
| Invalid/expired token | 401 | 40100 | invalid or expired token |

## Backend Layer Structure

```
internal/
├── api/
│   ├── auth.go          # Register, Login handlers
│   └── user.go          # GetCurrentUser, UpdateProfile, UpdatePassword handlers
├── dto/
│   ├── auth.go          # RegisterRequest, LoginRequest, LoginResponse (already exists, extend)
│   └── user.go          # UpdateProfileRequest, UpdatePasswordRequest, UserResponse
├── middleware/
│   └── auth.go          # Already implemented (JWT + CORS)
├── model/
│   └── user.go          # Already implemented
├── repository/
│   └── user.go          # FindByUsername, FindByEmail, Create, FindByID, Update
├── service/
│   ├── auth.go          # Register (hash password + create user + init categories + generate token)
│   └── user.go          # GetCurrentUser, UpdateProfile, UpdatePassword
├── router/
│   └── router.go        # Update: register public + protected routes
```

### Layer Responsibilities

**API layer (`internal/api/`):**
- Bind and validate HTTP request parameters using `gin.ShouldBindJSON`
- Call service methods
- Return unified response via `pkg/response`
- Extract `user_id` from gin context for authenticated endpoints

**Service layer (`internal/service/`):**
- Business logic: password hashing, token generation, category initialization
- Validate business rules (username/email uniqueness, old password correctness)
- Coordinate between repository calls
- Never directly access `gin.Context`

**Repository layer (`internal/repository/`):**
- Pure GORM database operations
- One struct per domain entity with methods for CRUD
- Always scope queries by user_id where applicable
- Return model structs, not DTOs

### DTO Definitions

**dto/auth.go** (extend existing):
```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6,max=50"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string       `json:"token"`
    User  UserResponse `json:"user"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Nickname string `json:"nickname"`
    Avatar   string `json:"avatar"`
}
```

**dto/user.go** (new):
```go
type UpdateProfileRequest struct {
    Nickname string `json:"nickname" binding:"max=50"`
    Avatar   string `json:"avatar" binding:"max=500"`
}

type UpdatePasswordRequest struct {
    OldPassword string `json:"old_password" binding:"required,min=6,max=50"`
    NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}
```

### Repository Interface

```go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error)
func (r *UserRepository) FindByEmail(email string) (*model.User, error)
func (r *UserRepository) FindByID(id uint) (*model.User, error)
func (r *UserRepository) Create(user *model.User) error
func (r *UserRepository) Update(user *model.User) error
```

### Service Methods

**auth.go:**
- `Register(req *dto.RegisterRequest) (*dto.LoginResponse, error)` — check uniqueness, hash password, create user, init default categories, generate token
- `Login(req *dto.LoginRequest) (*dto.LoginResponse, error)` — find user by username, verify password, generate token

**user.go:**
- `GetCurrentUser(userID uint) (*dto.UserResponse, error)` — find user by ID, convert to DTO
- `UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)` — find user, update fields, save
- `UpdatePassword(userID uint, req *dto.UpdatePasswordRequest) error` — find user, verify old password, hash new password, save

## Default Category Initialization

When a new user registers, the system creates 8 default categories with `type = "subscription"`:

| Name | Color | Icon | Sort Order |
|------|-------|------|------------|
| AI 工具 | #7c3aed | robot | 1 |
| 开发工具 | #2563eb | code | 2 |
| 云服务 | #0891b2 | cloud | 3 |
| 域名 | #059669 | globe | 4 |
| 服务器 | #d97706 | server | 5 |
| 娱乐 | #dc2626 | game-controller | 6 |
| 办公 | #4f46e5 | briefcase | 7 |
| 其他 | #6b7280 | ellipsis | 8 |

This happens inside `service/auth.go::Register` after user creation, using the same repository pattern.

## Frontend Changes

### API layer (`src/api/auth.ts`)
Add `getMe()`, `updateProfile()`, `updatePassword()` functions.

### Store (`src/stores/user.ts`)
Add `fetchUser()` method to restore user info from `/auth/me` on page refresh when token exists.

### Views
- **Login.vue / Register.vue** — connect to real API (already have forms, just wire up)
- **settings/Index.vue** — implement profile edit form and password change form

### Router
No changes needed — auth guard already implemented.

## Security Measures

- bcrypt with default cost factor (10) for password hashing
- Password minimum 6 characters, maximum 50
- Username 3-50 characters, alphanumeric
- Email format validation via binding tags
- Password never returned in any response (DTO filtering)
- Username and email uniqueness checked before registration
- Old password verification required for password change
- JWT token with configurable expiration (default 72h)
- All authenticated endpoints require valid Bearer token

## Files to Create/Modify

**New files:**
- `backend/internal/api/auth.go`
- `backend/internal/api/user.go`
- `backend/internal/service/auth.go`
- `backend/internal/service/user.go`
- `backend/internal/repository/user.go`
- `backend/internal/repository/category.go` (for default category init)
- `backend/internal/dto/user.go`

**Modified files:**
- `backend/internal/router/router.go` — add auth and user routes
- `backend/internal/dto/dto.go` — extend with UserResponse, LoginResponse
- `frontend/src/api/auth.ts` — add new API functions
- `frontend/src/stores/user.ts` — add fetchUser
- `frontend/src/views/auth/Login.vue` — wire to real API
- `frontend/src/views/auth/Register.vue` — wire to real API
- `frontend/src/views/settings/Index.vue` — implement forms
