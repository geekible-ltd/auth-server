package services

import (
	"strings"
	"time"

	"github.com/geekible-ltd/auth-server/src/internal/config"
	"github.com/geekible-ltd/auth-server/src/internal/dto"
	"github.com/geekible-ltd/auth-server/src/internal/entities"
	"github.com/geekible-ltd/auth-server/src/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRegistrationService struct {
	userRepository   *repositories.UserRepository
	tenantRepository *repositories.TenantRepository
}

func NewUserRegistrationService(userRepository *repositories.UserRepository, tenantRepository *repositories.TenantRepository) *UserRegistrationService {
	return &UserRegistrationService{
		userRepository:   userRepository,
		tenantRepository: tenantRepository,
	}
}

func (s *UserRegistrationService) RegisterTenant(tenantDTO *dto.TenantRegistrationDTO) error {
	emailDomain := strings.Split(tenantDTO.Email, "@")[1]
	_, err := s.tenantRepository.GetByEmailDomain(emailDomain)

	if err != nil && err == gorm.ErrRecordNotFound {
		return config.ErrTenantAlreadyExists
	} else if err != nil {
		return err
	}

	tenant := &entities.Tenant{
		Name:      tenantDTO.Name,
		Email:     tenantDTO.Email,
		Phone:     tenantDTO.Phone,
		Address:   tenantDTO.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.tenantRepository.Create(tenant); err != nil {
		return config.ErrFailedToCreateTenant
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(tenantDTO.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return config.ErrFailedToHashPassword
	}

	user := &entities.User{
		TenantID:                        tenant.ID,
		FirstName:                       tenantDTO.User.FirstName,
		LastName:                        tenantDTO.User.LastName,
		Email:                           tenantDTO.User.Email,
		PasswordHash:                    string(passwordHash),
		FailedLoginAttempts:             0,
		IsActive:                        true,
		Role:                            config.UserRoleTenantAdmin,
		LastLoginAt:                     nil,
		LastLoginIP:                     "",
		ResetPasswordToken:              "",
		ResetPasswordTokenExpiresAt:     nil,
		IsEmailVerified:                 false,
		EmailVerificationToken:          "",
		EmailVerificationTokenExpiresAt: nil,
		CreatedAt:                       time.Now(),
		UpdatedAt:                       time.Now(),
	}
	if err := s.userRepository.Create(user); err != nil {
		return config.ErrFailedToCreateUser
	}

	return nil
}

func (s *UserRegistrationService) RegisterUser(tenantId uint, userDTO *dto.UserRegistrationDTO) error {
	emailDomain := strings.Split(userDTO.Email, "@")[1]
	_, err := s.userRepository.GetByEmailDomain(emailDomain)

	if err != nil && err == gorm.ErrRecordNotFound {
		return config.ErrUserAlreadyExists
	} else if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return config.ErrFailedToHashPassword
	}

	user := &entities.User{
		TenantID:                        tenantId,
		FirstName:                       userDTO.FirstName,
		LastName:                        userDTO.LastName,
		Email:                           userDTO.Email,
		PasswordHash:                    string(passwordHash),
		FailedLoginAttempts:             0,
		IsActive:                        true,
		Role:                            config.UserRoleTenantUser,
		LastLoginAt:                     nil,
		LastLoginIP:                     "",
		ResetPasswordToken:              "",
		ResetPasswordTokenExpiresAt:     nil,
		IsEmailVerified:                 false,
		EmailVerificationToken:          "",
		EmailVerificationTokenExpiresAt: nil,
		CreatedAt:                       time.Now(),
		UpdatedAt:                       time.Now(),
	}
	if err := s.userRepository.Create(user); err != nil {
		return config.ErrFailedToCreateUser
	}

	return nil
}
