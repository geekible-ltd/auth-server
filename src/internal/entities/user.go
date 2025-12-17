package entities

import "time"

type User struct {
	ID                              uint       `json:"id"`
	TenantID                        uint       `json:"tenant_id"`
	FirstName                       string     `json:"first_name"`
	LastName                        string     `json:"last_name"`
	Email                           string     `json:"email"`
	PasswordHash                    string     `json:"password_hash"`
	FailedLoginAttempts             int        `json:"failed_login_attempts"`
	IsActive                        bool       `json:"is_active"`
	Role                            string     `json:"role"`
	LastLoginAt                     *time.Time `json:"last_login_at"`
	LastLoginIP                     string     `json:"last_login_ip"`
	ResetPasswordToken              string     `json:"reset_password_token"`
	ResetPasswordTokenExpiresAt     *time.Time `json:"reset_password_token_expires_at"`
	IsEmailVerified                 bool       `json:"is_email_verified"`
	EmailVerificationToken          string     `json:"email_verification_token"`
	EmailVerificationTokenExpiresAt *time.Time `json:"email_verification_token_expires_at"`
	CreatedAt                       time.Time  `json:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at"`
	DeletedAt                       *time.Time `json:"deleted_at"`

	Tenant Tenant `json:"tenant" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
