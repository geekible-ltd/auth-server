package config

import "errors"

const (
	UserRoleSuperAdmin  = "super_admin"
	UserRoleAdmin       = "admin"
	UserRoleTenantAdmin = "tenant_admin"
	UserRoleTenantUser  = "tenant_user"
)

var (
	ErrFailedToCreateTenant        = errors.New("failed to create tenant")
	ErrFailedToCreateUser          = errors.New("failed to create user")
	ErrFailedToHashPassword        = errors.New("failed to hash password")
	ErrTenantAlreadyExists         = errors.New("tenant already exists")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrUserNotFound                = errors.New("user not found")
	ErrInvalidPassword             = errors.New("invalid password")
	ErrTenantNotFound              = errors.New("tenant not found")
	ErrTenantLicenceNotFound       = errors.New("tenant licence not found")
	ErrTenantLicenceExceeded       = errors.New("tenant licence exceeded")
	ErrTenantLicenceExpired        = errors.New("tenant licence expired")
	ErrFailedToCreateTenantLicence = errors.New("failed to create tenant licence")
)

const MaxFailedLoginAttempts = 3

