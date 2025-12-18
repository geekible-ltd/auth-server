package authhandlers

import (
	"net/http"

	"github.com/geekible-ltd/auth-server/dto"
	"github.com/geekible-ltd/auth-server/internal/service"
	ginmiddleware "github.com/geekible-ltd/gin-middleware"
	authmodels "github.com/geekible-ltd/gin-middleware/auth-models"
	responseutils "github.com/geekible-ltd/response-utils"
	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	jwtSecret            string
	ginEngine            *gin.Engine
	LoginService         *service.LoginService
	RegistrationService  *service.UserRegistrationService
	TenantService        *service.TenantService
	UserService          *service.UserService
	TenantLicenceService *service.TenantLicenceService
}

func NewAuthHandlers(jwtSecret string, ginEngine *gin.Engine, loginService *service.LoginService, registrationService *service.UserRegistrationService, tenantService *service.TenantService, userService *service.UserService, tenantLicenceService *service.TenantLicenceService) *AuthHandlers {
	return &AuthHandlers{
		jwtSecret:            jwtSecret,
		ginEngine:            ginEngine,
		LoginService:         loginService,
		RegistrationService:  registrationService,
		TenantService:        tenantService,
		UserService:          userService,
		TenantLicenceService: tenantLicenceService,
	}
}

func (h *AuthHandlers) RegisterRoutes() {
	h.registerRegisterRoutes()
	h.registerLoginRoutes()
}

func (h *AuthHandlers) registerRegisterRoutes() {
	authGroup := h.ginEngine.Group("/register")
	{
		authGroup.POST("/new-tenant", func(ctx *gin.Context) {
			var tenantDTO dto.TenantRegistrationDTO
			if err := ctx.ShouldBindJSON(&tenantDTO); err != nil {
				responseutils.ErrorResponse(ctx, responseutils.BadRequest("Invalid request body"))
				return
			}
			if err := h.RegistrationService.RegisterTenant(tenantDTO); err != nil {
				responseutils.ErrorResponse(ctx, responseutils.InternalServerError("Failed to register tenant"))
				return
			}

			responseutils.SuccessResponse(ctx, http.StatusCreated, nil, "Tenant registered successfully")
		})

		authGroupProtected := authGroup.Group("/user-management")
		authGroupProtected.Use(ginmiddleware.BearerAuthMiddleware(h.jwtSecret))
		{
			authGroupProtected.POST("/new-user", func(ctx *gin.Context) {
				var userDTO dto.UserRegistrationDTO
				if err := ctx.ShouldBindJSON(&userDTO); err != nil {
					responseutils.ErrorResponse(ctx, responseutils.BadRequest("Invalid request body"))
					return
				}

				tokenData, exists := ctx.Get(ginmiddleware.TokenKey)
				if !exists {
					responseutils.ErrorResponse(ctx, responseutils.Unauthorized("Unauthorized"))
					return
				}

				token := tokenData.(authmodels.TokenDTO)

				if err := h.RegistrationService.RegisterUser(token.CompanyID.(uint), userDTO); err != nil {
					responseutils.ErrorResponse(ctx, responseutils.InternalServerError("Failed to register user"))
					return
				}

				responseutils.SuccessResponse(ctx, http.StatusCreated, nil, "User registered successfully")
			})
		}
	}
}

func (h *AuthHandlers) registerLoginRoutes() {
	authGroup := h.ginEngine.Group("/auth")
	{
		authGroup.POST("/login", func(ctx *gin.Context) {
			var loginDTO dto.LoginDTO
			if err := ctx.ShouldBindJSON(&loginDTO); err != nil {
				responseutils.ErrorResponse(ctx, responseutils.BadRequest("Invalid request body"))
				return
			}
			loginResponse, err := h.LoginService.Login(loginDTO, ctx.ClientIP())
			if err != nil {
				responseutils.ErrorResponse(ctx, responseutils.InternalServerError("Failed to login"))
				return
			}
			responseutils.SuccessResponse(ctx, http.StatusOK, loginResponse, "Login successful")
		})
	}
}
