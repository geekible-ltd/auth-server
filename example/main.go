package main

import (
	"fmt"
	"log"

	"github.com/geekible-ltd/auth-server"
	"github.com/geekible-ltd/auth-server/dto"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Connect to database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize auth server with all services
	authServer := authserver.NewAuthServer(db)
	if err := authServer.MigrateDB(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Register a tenant - use services from authServer
	tenantDTO := &dto.TenantRegistrationDTO{
		Name:    "Example Corp",
		Email:   "contact@example.com",
		Phone:   "+1234567890",
		Address: "123 Main St",
	}
	tenantDTO.User.FirstName = "Admin"
	tenantDTO.User.LastName = "User"
	tenantDTO.User.Email = "admin@example.com"
	tenantDTO.User.Password = "SecurePass123!"

	if err := authServer.RegistrationService.RegisterTenant(tenantDTO); err != nil {
		log.Fatal("Failed to register tenant:", err)
	}

	fmt.Println("✓ Tenant registered successfully")

	// Login
	loginDTO := dto.LoginDTO{
		Email:    "admin@example.com",
		Password: "SecurePass123!",
	}

	response, err := authServer.LoginService.Login(loginDTO, "127.0.0.1")
	if err != nil {
		log.Fatal("Login failed:", err)
	}

	fmt.Printf("✓ Login successful! User ID: %d, Role: %s\n", response.UserID, response.Role)
}

