package dto

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
	Status       string  `json:"status" binding:"omitempty,oneof=active inactive expired warning"`
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
	Status       string  `json:"status" binding:"omitempty,oneof=active inactive expired warning"`
	Description  string  `json:"description" binding:"max=500"`
	Remark       string  `json:"remark" binding:"max=500"`
}

type AssetListQuery struct {
	Page         int    `form:"page,default=1"`
	PageSize     int    `form:"page_size,default=20"`
	AssetType    string `form:"asset_type"`
	Status       string `form:"status"`
	ExpiringDays int    `form:"expiring_days"`
}
