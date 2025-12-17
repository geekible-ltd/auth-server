package services

import (
	"time"

	"github.com/geekible-ltd/auth-server/src/internal/config"
	"github.com/geekible-ltd/auth-server/src/internal/dto"
	"github.com/geekible-ltd/auth-server/src/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginService struct {
	userRepository   *repositories.UserRepository
	tenantRepository *repositories.TenantRepository
}

func NewLoginService(userRepository *repositories.UserRepository, tenantRepository *repositories.TenantRepository) *LoginService {
	return &LoginService{userRepository: userRepository, tenantRepository: tenantRepository}
}

func (s *LoginService) Login(loginRequest dto.LoginDTO, ipAddress string) (dto.LoginResponseDTO, error) {
	user, err := s.userRepository.GetByEmail(loginRequest.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return dto.LoginResponseDTO{}, config.ErrUserNotFound
	} else if err != nil {
		return dto.LoginResponseDTO{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
		return dto.LoginResponseDTO{}, config.ErrInvalidPassword
	}

	tenant, err := s.tenantRepository.GetByID(user.TenantID)
	if err != nil && err == gorm.ErrRecordNotFound {
		return dto.LoginResponseDTO{}, config.ErrTenantNotFound
	} else if err != nil {
		return dto.LoginResponseDTO{}, err
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = ipAddress
	if err := s.userRepository.Update(user); err != nil {
		return dto.LoginResponseDTO{}, err
	}

	return dto.LoginResponseDTO{
		TenantID: tenant.ID,
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}
