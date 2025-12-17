# Auth Server

A comprehensive, multi-tenant authentication and authorization package for Go applications. This package provides robust user and tenant management with built-in security features, role-based access control, and easy integration with GORM-supported databases.

## Features

- 🏢 **Multi-tenant architecture** - Complete tenant isolation and management
- 📜 **Licence management** - Comprehensive software licensing system with seat tracking and expiry dates
- 👤 **User management** - Registration, authentication, and profile management
- 🔐 **Secure password handling** - BCrypt password hashing
- 🛡️ **Security features**:
  - Failed login attempt tracking
  - Account lockout after multiple failed attempts
  - Last login tracking with IP address
  - Password reset token support
  - Email verification workflow
- 🎭 **Role-based access control** - Pre-defined roles (Super Admin, Admin, Tenant Admin, Tenant User)
- 🗄️ **GORM integration** - Works with any GORM-supported database (PostgreSQL, MySQL, SQLite, etc.)
- 📦 **Clean architecture** - Repository pattern, service layer, and DTOs for maintainability

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Database Setup](#database-setup)
- [Usage Guide](#usage-guide)
  - [Initialize Auth Server](#initialize-auth-server)
  - [Tenant Registration](#tenant-registration)
  - [User Registration](#user-registration)
  - [User Login](#user-login)
  - [Tenant Management](#tenant-management)
  - [User Management](#user-management)
  - [Licence Management](#licence-management)
- [API Reference](#api-reference)
- [Error Handling](#error-handling)
- [Configuration](#configuration)
- [Complete Example](#complete-example)
- [Security Best Practices](#security-best-practices)

## Installation

```bash
go get github.com/geekible-ltd/auth-server
```

### Dependencies

This package requires:
- Go 1.24.5 or higher
- GORM v1.31.1
- golang.org/x/crypto (for BCrypt)

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/geekible-ltd/auth-server/src/auth-service-registration"
    "github.com/geekible-ltd/auth-server/src/internal/dto"
    "github.com/geekible-ltd/auth-server/src/internal/repositories"
    "github.com/geekible-ltd/auth-server/src/internal/services"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    // 1. Connect to your database
    dsn := "host=localhost user=myuser password=mypass dbname=mydb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // 2. Initialize and migrate database
    authServer := authserviceregistration.NewAuthServerRegistration(db)
    if err := authServer.MigrateDBModels(); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
    
    // 3. Initialize services
    userRepo := repositories.NewUserRepository(db)
    tenantRepo := repositories.NewTenantRepository(db)
    registrationService := services.NewUserRegistrationService(userRepo, tenantRepo)
    loginService := services.NewLoginService(userRepo, tenantRepo)
    
    // 4. Register a new tenant with admin user
    tenantDTO := &dto.TenantRegistrationDTO{
        Name:    "Acme Corporation",
        Email:   "contact@acme.com",
        Phone:   "+1234567890",
        Address: "123 Main St, City, Country",
    }
    tenantDTO.User.FirstName = "John"
    tenantDTO.User.LastName = "Doe"
    tenantDTO.User.Email = "john.doe@acme.com"
    tenantDTO.User.Password = "SecurePassword123!"
    
    if err := registrationService.RegisterTenant(tenantDTO); err != nil {
        log.Fatal("Failed to register tenant:", err)
    }
    
    // 5. Login
    loginDTO := dto.LoginDTO{
        Email:    "john.doe@acme.com",
        Password: "SecurePassword123!",
    }
    
    loginResponse, err := loginService.Login(loginDTO, "192.168.1.1")
    if err != nil {
        log.Fatal("Login failed:", err)
    }
    
    fmt.Printf("Login successful! User ID: %d, Role: %s\n", 
        loginResponse.UserID, loginResponse.Role)
}
```

## Database Setup

### Step 1: Connect to Your Database

This package supports any GORM-compatible database. Here are examples for common databases:

#### PostgreSQL
```go
import "gorm.io/driver/postgres"

dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

#### MySQL
```go
import "gorm.io/driver/mysql"

dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

#### SQLite (for development)
```go
import "gorm.io/driver/sqlite"

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
```

### Step 2: Run Database Migration

```go
authServer := authserviceregistration.NewAuthServerRegistration(db)
if err := authServer.MigrateDBModels(); err != nil {
    log.Fatal("Migration failed:", err)
}
```

This will create three tables:
- `tenants` - Stores tenant/organization information
- `users` - Stores user information with foreign key to tenants
- `tenant_licences` - Stores licence information with one-to-one relationship to tenants

## Usage Guide

### Initialize Auth Server

First, initialize the database and create repository and service instances:

```go
package main

import (
    "gorm.io/gorm"
    "github.com/geekible-ltd/auth-server/src/auth-service-registration"
    "github.com/geekible-ltd/auth-server/src/internal/repositories"
    "github.com/geekible-ltd/auth-server/src/internal/services"
)

type AuthServerApp struct {
    DB                      *gorm.DB
    UserRepository          *repositories.UserRepository
    TenantRepository        *repositories.TenantRepository
    TenantLicenceRepository *repositories.TenantLicenceRepository
    RegistrationService     *services.UserRegistrationService
    LoginService            *services.LoginService
    TenantService           *services.TenantService
    UserService             *services.UserService
    TenantLicenceService    *services.TenantLicenceService
}

func InitializeAuthServer(db *gorm.DB) (*AuthServerApp, error) {
    // Migrate database
    authServer := authserviceregistration.NewAuthServerRegistration(db)
    if err := authServer.MigrateDBModels(); err != nil {
        return nil, err
    }
    
    // Initialize repositories
    userRepo := repositories.NewUserRepository(db)
    tenantRepo := repositories.NewTenantRepository(db)
    tenantLicenceRepo := repositories.NewTenantLicenceRepository(db)
    
    // Initialize services
    return &AuthServerApp{
        DB:                      db,
        UserRepository:          userRepo,
        TenantRepository:        tenantRepo,
        TenantLicenceRepository: tenantLicenceRepo,
        RegistrationService:     services.NewUserRegistrationService(userRepo, tenantRepo),
        LoginService:            services.NewLoginService(userRepo, tenantRepo),
        TenantService:           services.NewTenantService(tenantRepo),
        UserService:             services.NewUserService(userRepo),
        TenantLicenceService:    services.NewTenantLicenceService(tenantLicenceRepo),
    }, nil
}
```

### Tenant Registration

Register a new tenant (organization) with an admin user:

```go
import "github.com/geekible-ltd/auth-server/src/internal/dto"

func RegisterNewTenant(app *AuthServerApp) error {
    tenantDTO := &dto.TenantRegistrationDTO{
        Name:    "Tech Startup Inc",
        Email:   "contact@techstartup.com",
        Phone:   "+1-555-0100",
        Address: "456 Innovation Drive, Silicon Valley, CA",
    }
    
    // Set admin user details
    tenantDTO.User.FirstName = "Alice"
    tenantDTO.User.LastName = "Smith"
    tenantDTO.User.Email = "alice.smith@techstartup.com"
    tenantDTO.User.Password = "StrongPassword123!"
    
    err := app.RegistrationService.RegisterTenant(tenantDTO)
    if err != nil {
        return fmt.Errorf("tenant registration failed: %w", err)
    }
    
    fmt.Println("Tenant and admin user registered successfully!")
    return nil
}
```

**Key Points:**
- The admin user is automatically assigned the `tenant_admin` role
- The email domain from the tenant email is used for tenant identification
- Passwords are automatically hashed using BCrypt before storage
- Both tenant and user are created in a single transaction

### User Registration

Add additional users to an existing tenant:

```go
func RegisterNewUser(app *AuthServerApp, tenantID uint) error {
    userDTO := &dto.UserRegistrationDTO{
        TenantID:  tenantID,
        FirstName: "Bob",
        LastName:  "Johnson",
        Email:     "bob.johnson@techstartup.com",
        Password:  "SecurePass456!",
    }
    
    err := app.RegistrationService.RegisterUser(tenantID, userDTO)
    if err != nil {
        return fmt.Errorf("user registration failed: %w", err)
    }
    
    fmt.Println("User registered successfully!")
    return nil
}
```

**Key Points:**
- New users are assigned the `tenant_user` role by default
- Email domain must match the tenant's domain
- Users are automatically marked as active and email unverified

### User Login

Authenticate users and track login information:

```go
func LoginUser(app *AuthServerApp, email, password, ipAddress string) (*dto.LoginResponseDTO, error) {
    loginDTO := dto.LoginDTO{
        Email:    email,
        Password: password,
    }
    
    loginResponse, err := app.LoginService.Login(loginDTO, ipAddress)
    if err != nil {
        return nil, fmt.Errorf("login failed: %w", err)
    }
    
    fmt.Printf("Login successful!\n")
    fmt.Printf("User ID: %d\n", loginResponse.UserID)
    fmt.Printf("Tenant ID: %d\n", loginResponse.TenantID)
    fmt.Printf("Email: %s\n", loginResponse.Email)
    fmt.Printf("Role: %s\n", loginResponse.Role)
    
    return &loginResponse, nil
}
```

**Security Features:**
- Failed login attempts are tracked and incremented on wrong password
- After 3 failed attempts (configurable), the account is automatically deactivated
- Last login time and IP address are recorded
- Failed login counter is reset on successful login
- BCrypt is used for secure password comparison

### Tenant Management

#### Get Tenant by ID

```go
func GetTenantDetails(app *AuthServerApp, tenantID uint) error {
    tenant, err := app.TenantService.GetTenantByID(tenantID)
    if err != nil {
        return fmt.Errorf("failed to get tenant: %w", err)
    }
    
    fmt.Printf("Tenant: %s (%s)\n", tenant.Name, tenant.Email)
    fmt.Printf("Phone: %s\n", tenant.Phone)
    fmt.Printf("Address: %s\n", tenant.Address)
    
    return nil
}
```

#### Get All Tenants

```go
func ListAllTenants(app *AuthServerApp) error {
    tenants, err := app.TenantService.GetAllTenants()
    if err != nil {
        return fmt.Errorf("failed to get tenants: %w", err)
    }
    
    fmt.Printf("Total Tenants: %d\n", len(tenants))
    for _, tenant := range tenants {
        fmt.Printf("- %s (%s)\n", tenant.Name, tenant.Email)
    }
    
    return nil
}
```

#### Update Tenant

```go
func UpdateTenantInfo(app *AuthServerApp, tenantID uint) error {
    updateDTO := dto.TenantRequestDTO{
        Name:    "Tech Startup Inc (Updated)",
        Email:   "info@techstartup.com",
        Phone:   "+1-555-0200",
        Address: "789 New Address, Silicon Valley, CA",
    }
    
    err := app.TenantService.UpdateTenant(tenantID, updateDTO)
    if err != nil {
        return fmt.Errorf("failed to update tenant: %w", err)
    }
    
    fmt.Println("Tenant updated successfully!")
    return nil
}
```

#### Delete Tenant (Soft Delete)

```go
func DeactivateTenant(app *AuthServerApp, tenantID uint) error {
    err := app.TenantService.DeleteTenant(tenantID)
    if err != nil {
        return fmt.Errorf("failed to delete tenant: %w", err)
    }
    
    fmt.Println("Tenant deactivated successfully!")
    return nil
}
```

**Note:** Delete operations are soft deletes - the tenant is marked as inactive and a deleted timestamp is set.

### User Management

#### Get User by ID

```go
func GetUserDetails(app *AuthServerApp, tenantID, userID uint) error {
    user, err := app.UserService.GetUserByID(tenantID, userID)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }
    
    fmt.Printf("User: %s %s (%s)\n", user.FirstName, user.LastName, user.Email)
    fmt.Printf("Role: %s\n", user.Role)
    fmt.Printf("Active: %v\n", user.IsActive)
    if user.LastLoginAt != nil {
        fmt.Printf("Last Login: %s\n", user.LastLoginAt.Format("2006-01-02 15:04:05"))
    }
    
    return nil
}
```

#### Get All Users for a Tenant

```go
func ListTenantUsers(app *AuthServerApp, tenantID uint) error {
    users, err := app.UserService.GetAllUsers(tenantID)
    if err != nil {
        return fmt.Errorf("failed to get users: %w", err)
    }
    
    fmt.Printf("Total Users: %d\n", len(users))
    for _, user := range users {
        fmt.Printf("- %s %s (%s) - Role: %s\n", 
            user.FirstName, user.LastName, user.Email, user.Role)
    }
    
    return nil
}
```

#### Update User

```go
func UpdateUserInfo(app *AuthServerApp, tenantID, userID uint) error {
    updateDTO := dto.UserUpdateRequestDTO{
        FirstName: "Bob",
        LastName:  "Johnson Jr.",
        Email:     "bob.johnson@techstartup.com",
        Role:      "tenant_admin", // Promote to admin
    }
    
    err := app.UserService.UpdateUser(tenantID, userID, updateDTO)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    fmt.Println("User updated successfully!")
    return nil
}
```

#### Delete User (Soft Delete)

```go
func DeactivateUser(app *AuthServerApp, tenantID, userID uint) error {
    err := app.UserService.DeleteUser(tenantID, userID)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    fmt.Println("User deactivated successfully!")
    return nil
}
```

### Licence Management

The package includes a comprehensive licence management system to control tenant access based on software licenses. Each tenant can have a licence with seat limits and expiry dates.

#### Get Licence by Tenant ID

```go
func GetTenantLicence(app *AuthServerApp, tenantID uint) error {
    licence, err := app.TenantLicenceService.GetTenantLicenceByTenantID(tenantID)
    if err != nil {
        return fmt.Errorf("failed to get licence: %w", err)
    }
    
    fmt.Printf("Licence Key: %s\n", licence.LicenceKey)
    fmt.Printf("Licensed Seats: %d\n", licence.LicencedSeats)
    fmt.Printf("Used Seats: %d\n", licence.UsedSeats)
    fmt.Printf("Available Seats: %d\n", licence.LicencedSeats - licence.UsedSeats)
    
    if licence.ExpiryDate != nil {
        fmt.Printf("Expiry Date: %s\n", licence.ExpiryDate.Format("2006-01-02"))
    } else {
        fmt.Println("Expiry Date: No expiry")
    }
    
    return nil
}
```

#### Get Licence by Licence Key

```go
func ValidateLicenceKey(app *AuthServerApp, licenceKey string) error {
    licence, err := app.TenantLicenceService.GetTenantLicenceByLicenceKey(licenceKey)
    if err != nil {
        return fmt.Errorf("invalid licence key: %w", err)
    }
    
    fmt.Printf("Licence ID: %d\n", licence.ID)
    fmt.Println("Licence key is valid!")
    
    return nil
}
```

#### Get All Licences

```go
func ListAllLicences(app *AuthServerApp) error {
    licences, err := app.TenantLicenceService.GetAllTenantLicences()
    if err != nil {
        return fmt.Errorf("failed to get licences: %w", err)
    }
    
    fmt.Printf("Total Licences: %d\n", len(licences))
    for _, licence := range licences {
        status := "Active"
        if licence.ExpiryDate != nil && licence.ExpiryDate.Before(time.Now()) {
            status = "Expired"
        }
        
        fmt.Printf("- Tenant ID: %d, Key: %s, Seats: %d/%d, Status: %s\n",
            licence.TenantID,
            licence.LicenceKey,
            licence.UsedSeats,
            licence.LicencedSeats,
            status,
        )
    }
    
    return nil
}
```

#### Update Licence

```go
func UpdateLicence(app *AuthServerApp, tenantID uint) error {
    expiryDate := time.Now().AddDate(1, 0, 0) // 1 year from now
    
    updateDTO := &dto.TenantLicenceUpdateRequestDTO{
        LicenceKey:    "NEW-LICENCE-KEY-2024",
        LicencedSeats: 50, // Increase from 10 to 50 seats
        ExpiryDate:    &expiryDate,
    }
    
    err := app.TenantLicenceService.UpdateTenantLicence(tenantID, updateDTO)
    if err != nil {
        return fmt.Errorf("failed to update licence: %w", err)
    }
    
    fmt.Println("Licence updated successfully!")
    return nil
}
```

#### Check Licence Validity

```go
func CheckLicenceValidity(app *AuthServerApp, tenantID uint) (bool, error) {
    licence, err := app.TenantLicenceService.GetTenantLicenceByTenantID(tenantID)
    if err != nil {
        return false, err
    }
    
    // Check if expired
    if licence.ExpiryDate != nil && licence.ExpiryDate.Before(time.Now()) {
        return false, config.ErrTenantLicenceExpired
    }
    
    // Check if seats exceeded
    if licence.UsedSeats >= licence.LicencedSeats {
        return false, config.ErrTenantLicenceExceeded
    }
    
    fmt.Printf("Licence is valid. %d/%d seats used.\n", 
        licence.UsedSeats, licence.LicencedSeats)
    return true, nil
}
```

**Key Features:**
- **Seat Management**: Track licensed vs. used seats to control user limits per tenant
- **Expiry Tracking**: Set expiration dates for time-limited licenses
- **Licence Keys**: Unique identifiers for each tenant's license
- **Validation**: Built-in checks for expired licenses and seat limits

**Use Cases:**
- SaaS pricing tiers (e.g., 10 users for Basic, 50 for Pro, unlimited for Enterprise)
- Time-limited trials with automatic expiry
- License key validation for on-premise deployments
- Seat-based billing and access control

## API Reference

### Services

#### UserRegistrationService

```go
type UserRegistrationService struct {
    // ...
}

// Register a new tenant with an admin user
func (s *UserRegistrationService) RegisterTenant(tenantDTO *dto.TenantRegistrationDTO) error

// Register a new user under an existing tenant
func (s *UserRegistrationService) RegisterUser(tenantId uint, userDTO *dto.UserRegistrationDTO) error
```

#### LoginService

```go
type LoginService struct {
    // ...
}

// Authenticate a user and return login response
func (s *LoginService) Login(loginRequest dto.LoginDTO, ipAddress string) (dto.LoginResponseDTO, error)
```

#### TenantService

```go
type TenantService struct {
    // ...
}

// Get tenant by ID
func (s *TenantService) GetTenantByID(tenantId uint) (dto.TenantResponseDTO, error)

// Get all tenants
func (s *TenantService) GetAllTenants() ([]dto.TenantResponseDTO, error)

// Update tenant information
func (s *TenantService) UpdateTenant(tenantId uint, tenantDTO dto.TenantRequestDTO) error

// Soft delete tenant
func (s *TenantService) DeleteTenant(tenantId uint) error
```

#### UserService

```go
type UserService struct {
    // ...
}

// Get user by ID and tenant ID
func (s *UserService) GetUserByID(tenantId, userId uint) (dto.UserResponseDTO, error)

// Get all users for a tenant
func (s *UserService) GetAllUsers(tenantId uint) ([]dto.UserResponseDTO, error)

// Update user information
func (s *UserService) UpdateUser(tenantId, userId uint, userDTO dto.UserUpdateRequestDTO) error

// Soft delete user
func (s *UserService) DeleteUser(tenantId, userId uint) error
```

#### TenantLicenceService

```go
type TenantLicenceService struct {
    // ...
}

// Get tenant licence by tenant ID
func (s *TenantLicenceService) GetTenantLicenceByID(tenantID uint) (*dto.TenantLicenceResponseDTO, error)

// Get all tenant licences
func (s *TenantLicenceService) GetAllTenantLicences() ([]entities.TenantLicence, error)

// Get tenant licence by licence key
func (s *TenantLicenceService) GetTenantLicenceByLicenceKey(licenceKey string) (*dto.TenantLicenceResponseDTO, error)

// Get tenant licence by tenant ID
func (s *TenantLicenceService) GetTenantLicenceByTenantID(tenantID uint) (*dto.TenantLicenceResponseDTO, error)

// Update tenant licence
func (s *TenantLicenceService) UpdateTenantLicence(tenantID uint, tenantLicence *dto.TenantLicenceUpdateRequestDTO) error
```

### Data Transfer Objects (DTOs)

#### TenantRegistrationDTO
```go
type TenantRegistrationDTO struct {
    Name    string `json:"name"`
    Email   string `json:"email"`
    Phone   string `json:"phone"`
    Address string `json:"address"`
    User    struct {
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
        Email     string `json:"email"`
        Password  string `json:"password"`
    }
}
```

#### UserRegistrationDTO
```go
type UserRegistrationDTO struct {
    TenantID  uint   `json:"tenant_id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}
```

#### LoginDTO
```go
type LoginDTO struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

#### LoginResponseDTO
```go
type LoginResponseDTO struct {
    TenantID uint   `json:"tenant_id"`
    UserID   uint   `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}
```

#### TenantResponseDTO
```go
type TenantResponseDTO struct {
    ID      uint   `json:"id"`
    Name    string `json:"name"`
    Email   string `json:"email"`
    Phone   string `json:"phone"`
    Address string `json:"address"`
}
```

#### UserResponseDTO
```go
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
```

#### TenantLicenceResponseDTO
```go
type TenantLicenceResponseDTO struct {
    ID            uint       `json:"id"`
    TenantID      uint       `json:"tenant_id"`
    LicenceKey    string     `json:"licence_key"`
    LicencedSeats int        `json:"licenced_seats"`
    UsedSeats     int        `json:"used_seats"`
    ExpiryDate    *time.Time `json:"expiry_date"`
}
```

#### TenantLicenceUpdateRequestDTO
```go
type TenantLicenceUpdateRequestDTO struct {
    LicenceKey    string     `json:"licence_key"`
    LicencedSeats int        `json:"licenced_seats"`
    ExpiryDate    *time.Time `json:"expiry_date"`
}
```

## Error Handling

The package provides predefined error constants for common scenarios:

```go
import "github.com/geekible-ltd/auth-server/src/internal/config"

// Available error constants
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
```

### Error Handling Example

```go
import (
    "errors"
    "github.com/geekible-ltd/auth-server/src/internal/config"
)

func HandleLogin(app *AuthServerApp, email, password string) {
    loginDTO := dto.LoginDTO{Email: email, Password: password}
    
    response, err := app.LoginService.Login(loginDTO, "192.168.1.1")
    if err != nil {
        switch {
        case errors.Is(err, config.ErrUserNotFound):
            fmt.Println("User not found. Please check your email.")
        case errors.Is(err, config.ErrInvalidPassword):
            fmt.Println("Invalid password. Account may be locked after 3 failed attempts.")
        default:
            fmt.Printf("Login error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Welcome, %s!\n", response.Email)
}

func HandleLicenceValidation(app *AuthServerApp, tenantID uint) {
    licence, err := app.TenantLicenceService.GetTenantLicenceByTenantID(tenantID)
    if err != nil {
        switch {
        case errors.Is(err, config.ErrTenantLicenceNotFound):
            fmt.Println("No licence found for this tenant.")
        default:
            fmt.Printf("Licence error: %v\n", err)
        }
        return
    }
    
    // Check expiry
    if licence.ExpiryDate != nil && licence.ExpiryDate.Before(time.Now()) {
        fmt.Println("Licence has expired. Please renew your subscription.")
        return
    }
    
    // Check seat availability
    if licence.UsedSeats >= licence.LicencedSeats {
        fmt.Printf("All %d seats are in use. Please upgrade your plan.\n", licence.LicencedSeats)
        return
    }
    
    fmt.Printf("Licence valid. %d of %d seats available.\n", 
        licence.LicencedSeats - licence.UsedSeats, licence.LicencedSeats)
}
```

## Configuration

### User Roles

The package provides four pre-defined roles:

```go
import "github.com/geekible-ltd/auth-server/src/internal/config"

const (
    UserRoleSuperAdmin  = "super_admin"   // System-wide administrator
    UserRoleAdmin       = "admin"          // Platform administrator
    UserRoleTenantAdmin = "tenant_admin"   // Tenant administrator
    UserRoleTenantUser  = "tenant_user"    // Regular tenant user
)
```

### Security Configuration

```go
import "github.com/geekible-ltd/auth-server/src/internal/config"

const MaxFailedLoginAttempts = 3  // Account locked after 3 failed attempts
```

You can modify this constant in your fork if you need different security settings.

## Complete Example

Here's a complete, production-ready example with a RESTful API using Gin:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/geekible-ltd/auth-server/src/auth-service-registration"
    "github.com/geekible-ltd/auth-server/src/internal/config"
    "github.com/geekible-ltd/auth-server/src/internal/dto"
    "github.com/geekible-ltd/auth-server/src/internal/repositories"
    "github.com/geekible-ltd/auth-server/src/internal/services"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type App struct {
    DB                   *gorm.DB
    RegistrationService  *services.UserRegistrationService
    LoginService         *services.LoginService
    TenantService        *services.TenantService
    UserService          *services.UserService
    TenantLicenceService *services.TenantLicenceService
}

func main() {
    // Database connection
    dsn := "host=localhost user=postgres password=postgres dbname=authdb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Initialize auth server
    authServer := authserviceregistration.NewAuthServerRegistration(db)
    if err := authServer.MigrateDBModels(); err != nil {
        log.Fatal("Migration failed:", err)
    }
    
    // Initialize repositories and services
    userRepo := repositories.NewUserRepository(db)
    tenantRepo := repositories.NewTenantRepository(db)
    tenantLicenceRepo := repositories.NewTenantLicenceRepository(db)
    
    app := &App{
        DB:                   db,
        RegistrationService:  services.NewUserRegistrationService(userRepo, tenantRepo),
        LoginService:         services.NewLoginService(userRepo, tenantRepo),
        TenantService:        services.NewTenantService(tenantRepo),
        UserService:          services.NewUserService(userRepo),
        TenantLicenceService: services.NewTenantLicenceService(tenantLicenceRepo),
    }
    
    // Setup routes
    router := gin.Default()
    
    // Public routes
    router.POST("/register/tenant", app.RegisterTenantHandler)
    router.POST("/login", app.LoginHandler)
    
    // Protected routes (add your auth middleware here)
    protected := router.Group("/api")
    // protected.Use(YourAuthMiddleware())
    {
        protected.POST("/users", app.RegisterUserHandler)
        protected.GET("/tenants/:id", app.GetTenantHandler)
        protected.GET("/tenants", app.ListTenantsHandler)
        protected.PUT("/tenants/:id", app.UpdateTenantHandler)
        protected.DELETE("/tenants/:id", app.DeleteTenantHandler)
        protected.GET("/tenants/:tenantId/users/:userId", app.GetUserHandler)
        protected.GET("/tenants/:tenantId/users", app.ListUsersHandler)
        protected.PUT("/tenants/:tenantId/users/:userId", app.UpdateUserHandler)
        protected.DELETE("/tenants/:tenantId/users/:userId", app.DeleteUserHandler)
        protected.GET("/tenants/:tenantId/licence", app.GetTenantLicenceHandler)
        protected.PUT("/tenants/:tenantId/licence", app.UpdateTenantLicenceHandler)
        protected.GET("/licences", app.ListAllLicencesHandler)
    }
    
    log.Println("Server starting on :8080")
    router.Run(":8080")
}

// Handler: Register Tenant
func (app *App) RegisterTenantHandler(c *gin.Context) {
    var tenantDTO dto.TenantRegistrationDTO
    if err := c.ShouldBindJSON(&tenantDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := app.RegistrationService.RegisterTenant(&tenantDTO); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{"message": "Tenant registered successfully"})
}

// Handler: Login
func (app *App) LoginHandler(c *gin.Context) {
    var loginDTO dto.LoginDTO
    if err := c.ShouldBindJSON(&loginDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    ipAddress := c.ClientIP()
    response, err := app.LoginService.Login(loginDTO, ipAddress)
    if err != nil {
        switch err {
        case config.ErrUserNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        case config.ErrInvalidPassword:
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    
    c.JSON(http.StatusOK, response)
}

// Handler: Register User
func (app *App) RegisterUserHandler(c *gin.Context) {
    var userDTO dto.UserRegistrationDTO
    if err := c.ShouldBindJSON(&userDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := app.RegistrationService.RegisterUser(userDTO.TenantID, &userDTO); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Handler: Get Tenant
func (app *App) GetTenantHandler(c *gin.Context) {
    var tenantID uint
    if _, err := fmt.Sscanf(c.Param("id"), "%d", &tenantID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
        return
    }
    
    tenant, err := app.TenantService.GetTenantByID(tenantID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, tenant)
}

// Handler: List Tenants
func (app *App) ListTenantsHandler(c *gin.Context) {
    tenants, err := app.TenantService.GetAllTenants()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, tenants)
}

// Handler: Update Tenant
func (app *App) UpdateTenantHandler(c *gin.Context) {
    var tenantID uint
    if _, err := fmt.Sscanf(c.Param("id"), "%d", &tenantID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
        return
    }
    
    var tenantDTO dto.TenantRequestDTO
    if err := c.ShouldBindJSON(&tenantDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := app.TenantService.UpdateTenant(tenantID, tenantDTO); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Tenant updated successfully"})
}

// Handler: Delete Tenant
func (app *App) DeleteTenantHandler(c *gin.Context) {
    var tenantID uint
    if _, err := fmt.Sscanf(c.Param("id"), "%d", &tenantID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
        return
    }
    
    if err := app.TenantService.DeleteTenant(tenantID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}

// Handler: Get User
func (app *App) GetUserHandler(c *gin.Context) {
    var tenantID, userID uint
    fmt.Sscanf(c.Param("tenantId"), "%d", &tenantID)
    fmt.Sscanf(c.Param("userId"), "%d", &userID)
    
    user, err := app.UserService.GetUserByID(tenantID, userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// Handler: List Users
func (app *App) ListUsersHandler(c *gin.Context) {
    var tenantID uint
    fmt.Sscanf(c.Param("tenantId"), "%d", &tenantID)
    
    users, err := app.UserService.GetAllUsers(tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, users)
}

// Handler: Update User
func (app *App) UpdateUserHandler(c *gin.Context) {
    var tenantID, userID uint
    fmt.Sscanf(c.Param("tenantId"), "%d", &tenantID)
    fmt.Sscanf(c.Param("userId"), "%d", &userID)
    
    var userDTO dto.UserUpdateRequestDTO
    if err := c.ShouldBindJSON(&userDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := app.UserService.UpdateUser(tenantID, userID, userDTO); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Handler: Delete User
func (app *App) DeleteUserHandler(c *gin.Context) {
    var tenantID, userID uint
    fmt.Sscanf(c.Param("tenantId"), "%d", &tenantID)
    fmt.Sscanf(c.Param("userId"), "%d", &userID)
    
    if err := app.UserService.DeleteUser(tenantID, userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Handler: Get Tenant Licence
func (app *App) GetTenantLicenceHandler(c *gin.Context) {
    var tenantID uint
    fmt.Sscanf(c.Param("tenantId"), "%d", &tenantID)
    
    licence, err := app.TenantLicenceService.GetTenantLicenceByTenantID(tenantID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, licence)
}

// Handler: Update Tenant Licence
func (app *App) UpdateTenantLicenceHandler(c *gin.Context) {
    var tenantID uint
    fmt.Sscanf(c.Param("tenantId"), "%d", &tenantID)
    
    var licenceDTO dto.TenantLicenceUpdateRequestDTO
    if err := c.ShouldBindJSON(&licenceDTO); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := app.TenantLicenceService.UpdateTenantLicence(tenantID, &licenceDTO); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Licence updated successfully"})
}

// Handler: List All Licences
func (app *App) ListAllLicencesHandler(c *gin.Context) {
    licences, err := app.TenantLicenceService.GetAllTenantLicences()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, licences)
}
```

## Security Best Practices

1. **Use HTTPS in Production**: Always serve your application over HTTPS to protect credentials in transit.

2. **Strong Passwords**: Implement password strength validation before accepting user passwords:
   ```go
   func ValidatePassword(password string) bool {
       // Minimum 8 characters, at least one uppercase, one lowercase, one digit
       return len(password) >= 8 && 
              regexp.MustCompile(`[A-Z]`).MatchString(password) &&
              regexp.MustCompile(`[a-z]`).MatchString(password) &&
              regexp.MustCompile(`[0-9]`).MatchString(password)
   }
   ```

3. **Rate Limiting**: Implement rate limiting on login endpoints to prevent brute force attacks.

4. **JWT Tokens**: Extend the `LoginResponseDTO` to include JWT tokens for stateless authentication:
   ```go
   // After successful login
   token := GenerateJWT(loginResponse.UserID, loginResponse.TenantID, loginResponse.Role)
   ```

5. **Input Validation**: Always validate and sanitize user inputs before processing.

6. **Audit Logging**: Log all authentication and authorization events for security auditing.

7. **Regular Updates**: Keep dependencies up to date to patch security vulnerabilities.

8. **Environment Variables**: Store database credentials and secrets in environment variables, never in code:
   ```go
   dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
       os.Getenv("DB_HOST"),
       os.Getenv("DB_USER"),
       os.Getenv("DB_PASSWORD"),
       os.Getenv("DB_NAME"),
       os.Getenv("DB_PORT"),
   )
   ```

## Database Schema

### Users Table
- `id` - Primary key
- `tenant_id` - Foreign key to tenants
- `first_name` - User's first name
- `last_name` - User's last name
- `email` - Unique email address
- `password_hash` - BCrypt hashed password
- `failed_login_attempts` - Counter for failed logins
- `is_active` - Account status
- `role` - User role (super_admin, admin, tenant_admin, tenant_user)
- `last_login_at` - Timestamp of last successful login
- `last_login_ip` - IP address of last login
- `reset_password_token` - Token for password reset (future feature)
- `reset_password_token_expires_at` - Expiry for reset token
- `is_email_verified` - Email verification status
- `email_verification_token` - Token for email verification
- `email_verification_token_expires_at` - Expiry for verification token
- `created_at` - Record creation timestamp
- `updated_at` - Record update timestamp
- `deleted_at` - Soft delete timestamp

### Tenants Table
- `id` - Primary key
- `name` - Tenant/organization name
- `email` - Contact email
- `phone` - Contact phone
- `address` - Physical address
- `is_active` - Tenant status
- `created_at` - Record creation timestamp
- `updated_at` - Record update timestamp
- `deleted_at` - Soft delete timestamp

### Tenant Licences Table
- `id` - Primary key
- `tenant_id` - Foreign key to tenants (one-to-one relationship)
- `licence_key` - Unique licence key identifier
- `licenced_seats` - Maximum number of user seats allowed
- `used_seats` - Current number of seats in use
- `expiry_date` - Licence expiration date (nullable)
- `created_at` - Record creation timestamp
- `updated_at` - Record update timestamp
- `deleted_at` - Soft delete timestamp

**Relationship**: Each tenant can have one licence. The licence controls access limits and expiry for the tenant's subscription.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the terms specified in the LICENSE file.

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/geekible-ltd/auth-server).
