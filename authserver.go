package authserver

import (
	authhandlers "github.com/geekible-ltd/auth-server/auth-handlers"
	"github.com/geekible-ltd/auth-server/internal/models"
	"github.com/geekible-ltd/auth-server/internal/repository"
	"github.com/geekible-ltd/auth-server/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthServer provides database migration and initialization for the auth server
type AuthServer struct {
	db                   *gorm.DB
	jwtSecret            string
	LoginService         *service.LoginService
	RegistrationService  *service.UserRegistrationService
	TenantService        *service.TenantService
	UserService          *service.UserService
	TenantLicenceService *service.TenantLicenceService
}

// New creates a new AuthServer instance
func NewAuthServer(db *gorm.DB, jwtSecret string) *AuthServer {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	tenantLicenceRepo := repository.NewTenantLicenceRepository(db)

	// Initialize services with repositories
	return &AuthServer{
		db:                   db,
		jwtSecret:            jwtSecret,
		LoginService:         service.NewLoginService(userRepo, tenantRepo),
		RegistrationService:  service.NewUserRegistrationService(userRepo, tenantRepo, tenantLicenceRepo),
		TenantService:        service.NewTenantService(tenantRepo),
		UserService:          service.NewUserService(userRepo),
		TenantLicenceService: service.NewTenantLicenceService(tenantLicenceRepo),
	}
}

// MigrateDB runs automatic database migrations for all auth server models
func (a *AuthServer) MigrateDB() error {
	return a.db.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.TenantLicence{},
	)
}

func (a *AuthServer) RegisterRoutes() *gin.Engine {
	authHandlers := authhandlers.NewAuthHandlers(a.jwtSecret, a.LoginService, a.RegistrationService, a.TenantService, a.UserService, a.TenantLicenceService)
	router := authHandlers.RegisterRoutes()
	return router
}
