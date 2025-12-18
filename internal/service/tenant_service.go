package service

import (
	"time"

	"github.com/geekible-ltd/auth-server/dto"
	"github.com/geekible-ltd/auth-server/internal/config"
	"github.com/geekible-ltd/auth-server/internal/repository"
	"gorm.io/gorm"
)

type TenantService struct {
	tenantRepository *repository.TenantRepository
}

func NewTenantService(tenantRepository *repository.TenantRepository) *TenantService {
	return &TenantService{tenantRepository: tenantRepository}
}

func (s *TenantService) GetTenantByID(tenantId uint) (dto.TenantResponseDTO, error) {
	tenant, err := s.tenantRepository.GetByID(tenantId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return dto.TenantResponseDTO{}, config.ErrTenantNotFound
	} else if err != nil {
		return dto.TenantResponseDTO{}, err
	}
	return dto.TenantResponseDTO{
		ID:      tenant.ID,
		Name:    tenant.Name,
		Email:   tenant.Email,
		Phone:   tenant.Phone,
		Address: tenant.Address,
	}, nil
}

func (s *TenantService) GetAllTenants() ([]dto.TenantResponseDTO, error) {
	tenantsDTO := []dto.TenantResponseDTO{}

	tenants, err := s.tenantRepository.GetAll()
	if err != nil {
		return nil, err
	}
	for _, tenant := range tenants {
		tenantsDTO = append(tenantsDTO, dto.TenantResponseDTO{
			ID:      tenant.ID,
			Name:    tenant.Name,
			Email:   tenant.Email,
			Phone:   tenant.Phone,
			Address: tenant.Address,
		})
	}
	return tenantsDTO, nil
}

func (s *TenantService) UpdateTenant(tenantId uint, tenantDTO dto.TenantRequestDTO) error {
	tenant, err := s.tenantRepository.GetByID(tenantId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return config.ErrTenantNotFound
	} else if err != nil {
		return err
	}

	tenant.Name = tenantDTO.Name
	tenant.Email = tenantDTO.Email
	tenant.Phone = tenantDTO.Phone
	tenant.Address = tenantDTO.Address
	tenant.UpdatedAt = time.Now()

	return s.tenantRepository.Update(tenant)
}

func (s *TenantService) DeleteTenant(tenantId uint) error {
	tenant, err := s.tenantRepository.GetByID(tenantId)
	if err != nil && err == gorm.ErrRecordNotFound {
		return config.ErrTenantNotFound
	} else if err != nil {
		return err
	}

	tenant.IsActive = false
	tenant.UpdatedAt = time.Now()
	tenant.DeletedAt = time.Now()

	return s.tenantRepository.Delete(tenant)
}
