package service

import (
	"time"

	"github.com/geekible-ltd/auth-server/dto"
	"github.com/geekible-ltd/auth-server/internal/config"
	"github.com/geekible-ltd/auth-server/internal/repository"
	"gorm.io/gorm"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetUserByID(tenantId, userId uint) (dto.UserResponseDTO, error) {
	user, err := s.userRepository.GetByID(tenantId, userId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return dto.UserResponseDTO{}, config.ErrUserNotFound
	} else if err != nil {
		return dto.UserResponseDTO{}, err
	}

	return dto.UserResponseDTO{
		ID:          user.ID,
		TenantID:    user.TenantID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Role:        user.Role,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *UserService) GetAllUsers(tenantId uint) ([]dto.UserResponseDTO, error) {
	users, err := s.userRepository.GetAll(tenantId)
	if err != nil {
		return nil, err
	}

	usersDTO := []dto.UserResponseDTO{}
	for _, user := range users {
		usersDTO = append(usersDTO, dto.UserResponseDTO{
			ID:          user.ID,
			TenantID:    user.TenantID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			Role:        user.Role,
			IsActive:    user.IsActive,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
		})
	}
	return usersDTO, nil
}

func (s *UserService) UpdateUser(tenantId, userId uint, userDTO dto.UserUpdateRequestDTO) error {
	user, err := s.userRepository.GetByID(tenantId, userId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return config.ErrUserNotFound
	} else if err != nil {
		return err
	}

	user.FirstName = userDTO.FirstName
	user.LastName = userDTO.LastName
	user.Email = userDTO.Email
	user.Role = userDTO.Role
	user.UpdatedAt = time.Now()

	return s.userRepository.Update(user)
}

func (s *UserService) DeleteUser(tenantId, userId uint) error {
	user, err := s.userRepository.GetByID(tenantId, userId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return config.ErrUserNotFound
	} else if err != nil {
		return err
	}

	now := time.Now()

	user.IsActive = false
	user.UpdatedAt = time.Now()
	user.DeletedAt = &now

	return s.userRepository.Delete(user)
}
