package main

import (
	"log"

	"github.com/geekible-ltd/auth-server"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Connect to database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize auth server with JWT secret
	jwtSecret := "your-secret-key-change-in-production"
	authServer := authserver.NewAuthServer(db, jwtSecret)
	if err := authServer.MigrateDB(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Create Gin router
	router := gin.Default()

	// Register auth routes
	authServer.RegisterRoutes(router)

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	log.Println("Server starting on :8080")
	log.Println("Try: POST http://localhost:8080/register/new-tenant")
	log.Println("Try: POST http://localhost:8080/auth/login")
	router.Run(":8080")
}

