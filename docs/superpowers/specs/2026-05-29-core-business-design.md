# StackBill Phase 3: Core Business Modules Design

## Overview

Implement the four core business modules: Category, Subscription, Asset, and Reminder management. Each follows the three-layer architecture established in Phase 2 (Repository → Service → API Handler) with DTO isolation.

**Architecture:** Reuse the exact pattern from Phase 2. Each module gets its own repository, service, API handler, and DTO files. The router wires them into protected route groups behind JWT middleware.

**Reminder strategy:** Dynamic generation on query — no background tasks or pre-seeded records. Read status tracked via a lightweight `reminder_reads` table.

## API Routes

All routes are under `/api/v1` and require JWT authentication.

### Categories

| Method | Path | Description |
|--------|------|-------------|
| GET | `/categories` | List (optional `type` filter: subscription/asset) |
| GET | `/categories/:id` | Detail |
| POST | `/categories` | Create |
| PUT | `/categories/:id` | Update |
| DELETE | `/categories/:id` | Delete |

### Subscriptions

| Method | Path | Description |
|--------|------|-------------|
| GET | `/subscriptions` | List (paginated, filters: category_id, status, upcoming_renewal) |
| GET | `/subscriptions/:id` | Detail |
| POST | `/subscriptions` | Create |
| PUT | `/subscriptions/:id` | Update |
| DELETE | `/subscriptions/:id` | Delete |

### Assets

| Method | Path | Description |
|--------|------|-------------|
| GET | `/assets` | List (paginated, filters: asset_type, status, expiring_days) |
| GET | `/assets/:id` | Detail |
| POST | `/assets` | Create |
| PUT | `/assets/:id` | Update |
| DELETE | `/assets/:id` | Delete |

### Reminders

| Method | Path | Description |
|--------|------|-------------|
| GET | `/reminders` | List (dynamically generated from subscriptions + assets) |
| PUT | `/reminders/:id/read` | Mark single reminder as read |
| PUT | `/reminders/read-all` | Mark all reminders as read |

## Error Codes

Extend the existing pattern. Each module uses a base code:

| Module | Base Code | Scenarios |
|--------|-----------|-----------|
| Common | 400/40001 | Invalid parameters |
| Common | 404/40400 | Resource not found |
| Common | 403/40301 | Resource belongs to another user |
| Category | 409/40901 | Category name already exists for this user |
| Subscription | N/A | Uses common codes |
| Asset | N/A | Uses common codes |

## Backend File Structure

### New Files per Module

```
internal/
├── repository/
│   ├── category.go
│   ├── subscription.go
│   ├── asset.go
│   └── reminder.go
├── service/
│   ├── category.go
│   ├── subscription.go
│   ├── asset.go
│   └── reminder.go
├── api/
│   ├── category.go
│   ├── subscription.go
│   ├── asset.go
│   └── reminder.go
├── dto/
│   ├── category.go
│   ├── subscription.go
│   ├── asset.go
│   └── reminder.go
```

### Modified Files

- `internal/router/router.go` — Add 4 new route groups with handlers
- `internal/model/reminder.go` — Add ReminderRead model for tracking read status

## DTO Definitions

### Category DTOs

```go
type CategoryResponse struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    Type      string `json:"type"`
    Color     string `json:"color"`
    Icon      string `json:"icon"`
    SortOrder int    `json:"sort_order"`
}

type CreateCategoryRequest struct {
    Name      string `json:"name" binding:"required,max=50"`
    Type      string `json:"type" binding:"required,oneof=subscription asset"`
    Color     string `json:"color" binding:"max=20"`
    Icon      string `json:"icon" binding:"max=50"`
    SortOrder int    `json:"sort_order"`
}

type UpdateCategoryRequest struct {
    Name      string `json:"name" binding:"required,max=50"`
    Type      string `json:"type" binding:"required,oneof=subscription asset"`
    Color     string `json:"color" binding:"max=20"`
    Icon      string `json:"icon" binding:"max=50"`
    SortOrder int    `json:"sort_order"`
}

type CategoryListQuery struct {
    Type string `form:"type"`
}
```

### Subscription DTOs

```go
type SubscriptionResponse struct {
    ID              uint       `json:"id"`
    Name            string     `json:"name"`
    Description     string     `json:"description"`
    CategoryID      uint       `json:"category_id"`
    Amount          float64    `json:"amount"`
    Currency        string     `json:"currency"`
    BillingCycle    string     `json:"billing_cycle"`
    BillingInterval int        `json:"billing_interval"`
    NextPaymentDate *string    `json:"next_payment_date"`
    StartDate       *string    `json:"start_date"`
    PaymentMethod   string     `json:"payment_method"`
    AutoRenew       bool       `json:"auto_renew"`
    Status          string     `json:"status"`
    WebsiteURL      string     `json:"website_url"`
    Remark          string     `json:"remark"`
    CreatedAt       string     `json:"created_at"`
    UpdatedAt       string     `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
    Name            string  `json:"name" binding:"required,max=100"`
    Description     string  `json:"description" binding:"max=500"`
    CategoryID      uint    `json:"category_id"`
    Amount          float64 `json:"amount" binding:"required"`
    Currency        string  `json:"currency" binding:"max=10"`
    BillingCycle    string  `json:"billing_cycle" binding:"required,oneof=weekly monthly quarterly yearly custom one_time"`
    BillingInterval int     `json:"billing_interval"`
    StartDate       *string `json:"start_date"`
    PaymentMethod   string  `json:"payment_method" binding:"max=50"`
    AutoRenew       *bool   `json:"auto_renew"`
    Status          string  `json:"status" binding:"oneof=active paused cancelled expired"`
    WebsiteURL      string  `json:"website_url" binding:"max=500"`
    Remark          string  `json:"remark" binding:"max=500"`
}

type UpdateSubscriptionRequest struct {
    Name            string  `json:"name" binding:"max=100"`
    Description     string  `json:"description" binding:"max=500"`
    CategoryID      uint    `json:"category_id"`
    Amount          float64 `json:"amount"`
    Currency        string  `json:"currency" binding:"max=10"`
    BillingCycle    string  `json:"billing_cycle" binding:"oneof=weekly monthly quarterly yearly custom one_time"`
    BillingInterval int     `json:"billing_interval"`
    StartDate       *string `json:"start_date"`
    PaymentMethod   string  `json:"payment_method" binding:"max=50"`
    AutoRenew       *bool   `json:"auto_renew"`
    Status          string  `json:"status" binding:"oneof=active paused cancelled expired"`
    WebsiteURL      string  `json:"website_url" binding:"max=500"`
    Remark          string  `json:"remark" binding:"max=500"`
}

type SubscriptionListQuery struct {
    Page            int    `form:"page,default=1"`
    PageSize        int    `form:"page_size,default=20"`
    CategoryID      *uint  `form:"category_id"`
    Status          string `form:"status"`
    UpcomingRenewal bool   `form:"upcoming_renewal"`
}
```

### Asset DTOs

```go
type AssetResponse struct {
    ID           uint    `json:"id"`
    Name         string  `json:"name"`
    AssetType    string  `json:"asset_type"`
    Provider     string  `json:"provider"`
    Identifier   string  `json:"identifier"`
    URL          string  `json:"url"`
    ExpireDate   *string `json:"expire_date"`
    CostAmount   float64 `json:"cost_amount"`
    CostCurrency string  `json:"cost_currency"`
    BillingCycle string  `json:"billing_cycle"`
    Status       string  `json:"status"`
    Description  string  `json:"description"`
    Remark       string  `json:"remark"`
    CreatedAt    string  `json:"created_at"`
    UpdatedAt    string  `json:"updated_at"`
}

type CreateAssetRequest struct {
    Name         string  `json:"name" binding:"required,max=100"`
    AssetType    string  `json:"asset_type" binding:"required,oneof=domain server docker_service ssl_certificate api_key repository other"`
    Provider     string  `json:"provider" binding:"max=100"`
    Identifier   string  `json:"identifier" binding:"max=200"`
    URL          string  `json:"url" binding:"max=500"`
    ExpireDate   *string `json:"expire_date"`
    CostAmount   float64 `json:"cost_amount"`
    CostCurrency string  `json:"cost_currency" binding:"max=10"`
    BillingCycle string  `json:"billing_cycle" binding:"max=20"`
    Status       string  `json:"status" binding:"oneof=active inactive expired warning"`
    Description  string  `json:"description" binding:"max=500"`
    Remark       string  `json:"remark" binding:"max=500"`
}

type UpdateAssetRequest struct {
    Name         string  `json:"name" binding:"max=100"`
    AssetType    string  `json:"asset_type" binding:"oneof=domain server docker_service ssl_certificate api_key repository other"`
    Provider     string  `json:"provider" binding:"max=100"`
    Identifier   string  `json:"identifier" binding:"max=200"`
    URL          string  `json:"url" binding:"max=500"`
    ExpireDate   *string `json:"expire_date"`
    CostAmount   float64 `json:"cost_amount"`
    CostCurrency string  `json:"cost_currency" binding:"max=10"`
    BillingCycle string  `json:"billing_cycle" binding:"max=20"`
    Status       string  `json:"status" binding:"oneof=active inactive expired warning"`
    Description  string  `json:"description" binding:"max=500"`
    Remark       string  `json:"remark" binding:"max=500"`
}

type AssetListQuery struct {
    Page          int    `form:"page,default=1"`
    PageSize      int    `form:"page_size,default=20"`
    AssetType     string `form:"asset_type"`
    Status        string `form:"status"`
    ExpiringDays  int    `form:"expiring_days"`
}
```

### Reminder DTOs

```go
type ReminderResponse struct {
    ID         uint   `json:"id"`
    TargetType string `json:"target_type"`
    TargetID   uint   `json:"target_id"`
    RemindType string `json:"remind_type"`
    RemindDate string `json:"remind_date"`
    Title      string `json:"title"`
    Content    string `json:"content"`
    IsRead     bool   `json:"is_read"`
}

type ReminderListQuery struct {
    Page     int    `form:"page,default=1"`
    PageSize int    `form:"page_size,default=20"`
    Type     string `form:"type"`
    IsRead   *bool  `form:"is_read"`
}
```

## Repository Layer

### CategoryRepository

```go
func (r *CategoryRepository) List(userID uint, categoryType string) ([]model.Category, error)
func (r *CategoryRepository) FindByID(id uint) (*model.Category, error)
func (r *CategoryRepository) Create(cat *model.Category) error
func (r *CategoryRepository) Update(cat *model.Category) error
func (r *CategoryRepository) Delete(id uint) error
```

### SubscriptionRepository

```go
func (r *SubscriptionRepository) List(userID uint, query *dto.SubscriptionListQuery) ([]model.Subscription, int64, error)
func (r *SubscriptionRepository) FindByID(id uint) (*model.Subscription, error)
func (r *SubscriptionRepository) Create(sub *model.Subscription) error
func (r *SubscriptionRepository) Update(sub *model.Subscription) error
func (r *SubscriptionRepository) Delete(id uint) error
func (r *SubscriptionRepository) GetActiveByUserID(userID uint) ([]model.Subscription, error)
```

### AssetRepository

```go
func (r *AssetRepository) List(userID uint, query *dto.AssetListQuery) ([]model.Asset, int64, error)
func (r *AssetRepository) FindByID(id uint) (*model.Asset, error)
func (r *AssetRepository) Create(asset *model.Asset) error
func (r *AssetRepository) Update(asset *model.Asset) error
func (r *AssetRepository) Delete(id uint) error
func (r *AssetRepository) GetByUserID(userID uint) ([]model.Asset, error)
```

### ReminderRepository

```go
func (r *ReminderRepository) GetReadKeys(userID uint) (map[string]bool, error)
func (r *ReminderRepository) MarkRead(userID uint, targetType string, targetID uint) error
func (r *ReminderRepository) MarkAllRead(userID uint) error
```

## Service Layer

### CategoryService

- `List(userID uint, query *dto.CategoryListQuery) ([]dto.CategoryResponse, error)` — list with optional type filter
- `GetByID(userID uint, id uint) (*dto.CategoryResponse, error)` — find + verify ownership
- `Create(userID uint, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)` — check name uniqueness per user
- `Update(userID uint, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)` — find + verify ownership + update
- `Delete(userID uint, id uint) error` — find + verify ownership + delete

### SubscriptionService

- `List(userID uint, query *dto.SubscriptionListQuery) (*response.PageResult, error)` — paginated list with filters
- `GetByID(userID uint, id uint) (*dto.SubscriptionResponse, error)` — find + verify ownership
- `Create(userID uint, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error)` — create + calculate next_payment_date
- `Update(userID uint, id uint, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error)` — find + verify + update + recalc next_payment_date
- `Delete(userID uint, id uint) error` — find + verify + delete
- `CalculateMonthlyExpense(userID uint) (float64, error)` — sum active subscriptions normalized to monthly
- `CalculateYearlyExpense(userID uint) (float64, error)` — sum active subscriptions normalized to yearly

**Next payment date calculation:**
- Given `start_date`, `billing_cycle`, `billing_interval`
- Find the next occurrence after today
- Cycles: weekly (×7 days), monthly (×1 month), quarterly (×3 months), yearly (×12 months)
- If `one_time`: set next_payment_date to null

### AssetService

- `List(userID uint, query *dto.AssetListQuery) (*response.PageResult, error)` — paginated with filters
- `GetByID(userID uint, id uint) (*dto.AssetResponse, error)` — find + verify ownership
- `Create(userID uint, req *dto.CreateAssetRequest) (*dto.AssetResponse, error)` — create
- `Update(userID uint, id uint, req *dto.UpdateAssetRequest) (*dto.AssetResponse, error)` — find + verify + update
- `Delete(userID uint, id uint) error` — find + verify + delete

### ReminderService

- `List(userID uint, query *dto.ReminderListQuery) (*response.PageResult, error)` — dynamically generate from subscriptions + assets
- `MarkRead(userID uint, id uint) error` — mark a reminder as read by composite key (target_type + target_id)
- `MarkAllRead(userID uint) error` — mark all as read

**Dynamic generation logic:**
1. Query active subscriptions with `next_payment_date` within 7 days → `subscription_renewal` reminders
2. Query assets with `expire_date` within 30 days → `asset_expiration` reminders
3. Query assets with status `warning` or `inactive` → `service_warning` reminders
4. Cross-reference with `reminder_reads` table to set `is_read` flag
5. Each reminder gets a synthetic ID (hash of target_type + target_id) for consistency

## ReminderRead Model

Add a new model for tracking read status:

```go
type ReminderRead struct {
    ID         uint   `gorm:"primarykey" json:"id"`
    UserID     uint   `gorm:"index;not null" json:"user_id"`
    TargetType string `gorm:"size:30;not null" json:"target_type"`
    TargetID   uint   `gorm:"not null" json:"target_id"`
    CreatedAt  time.Time `json:"created_at"`
}
```

Unique index on `(user_id, target_type, target_id)`.

## Router Changes

Add to the `authorized` group in `router.go`:

```go
// Categories
authorized.GET("/categories", categoryHandler.List)
authorized.GET("/categories/:id", categoryHandler.GetByID)
authorized.POST("/categories", categoryHandler.Create)
authorized.PUT("/categories/:id", categoryHandler.Update)
authorized.DELETE("/categories/:id", categoryHandler.Delete)

// Subscriptions
authorized.GET("/subscriptions", subHandler.List)
authorized.GET("/subscriptions/:id", subHandler.GetByID)
authorized.POST("/subscriptions", subHandler.Create)
authorized.PUT("/subscriptions/:id", subHandler.Update)
authorized.DELETE("/subscriptions/:id", subHandler.Delete)

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
```

## Frontend Changes

### New API Files

- `src/api/category.ts` — list, getByID, create, update, delete
- `src/api/subscription.ts` — list, getByID, create, update, delete
- `src/api/asset.ts` — list, getByID, create, update, delete
- `src/api/reminder.ts` — list, markRead, markAllRead

### New Stores

- `src/stores/category.ts` — category list with caching
- `src/stores/subscription.ts` — subscription list state
- `src/stores/asset.ts` — asset list state

### View Implementations

Replace all placeholder views with functional implementations:

- **subscription/Index.vue** — DataTable with columns (name, amount, cycle, next payment, status), filters (category, status), create button
- **subscription/Detail.vue** — Card layout showing all subscription details
- **subscription/Edit.vue** — Form for create/edit with validation, category dropdown, billing cycle selector
- **asset/Index.vue** — DataTable with columns (name, type, provider, expire date, status), filters (type, status), create button
- **asset/Detail.vue** — Card layout showing all asset details
- **asset/Edit.vue** — Form for create/edit with validation, asset type selector
- **category/Index.vue** — List of categories with inline create/edit/delete, color picker
- **reminder/Index.vue** — List of reminders grouped by type, mark as read buttons, unread badge

### i18n Additions

Add keys for all 4 modules in both zh-CN and en-US covering: field labels, status values, billing cycles, asset types, button text, validation messages.

## Security

- All queries scoped by `user_id` from JWT token
- Ownership verification before update/delete operations
- Parameter validation via Gin binding tags
- Soft delete via GORM `DeletedAt` field (already in base Model)
