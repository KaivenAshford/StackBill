package dto

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=500"`
}

type UpdatePasswordRequest struct {
	OldPassword    string `json:"old_password" binding:"required,min=6,max=50"`
	NewPassword    string `json:"new_password" binding:"required,min=6,max=50"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=50"`
}