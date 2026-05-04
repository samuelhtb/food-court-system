package dto

// RegisterRequest
type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required"`
	TenantName string `json:"tenant_name"`
}

// LoginRequest 
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 
type UserResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	TenantName string `json:"tenant_name,omitempty"`
}