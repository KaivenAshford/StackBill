package dto

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
