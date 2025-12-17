package dto

import "time"

type UserResponseDTO struct {
	ID          uint       `json:"id"`
	TenantID    uint       `json:"tenant_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type UserUpdateRequestDTO struct {
	ID          uint       `json:"id"`
	TenantID    uint       `json:"tenant_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
}
