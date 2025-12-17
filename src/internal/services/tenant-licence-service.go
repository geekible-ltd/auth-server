package services

import (
	"time"

	"github.com/geekible-ltd/auth-server/src/internal/config"
	"github.com/geekible-ltd/auth-server/src/internal/dto"
	"github.com/geekible-ltd/auth-server/src/internal/entities"
	"github.com/geekible-ltd/auth-server/src/internal/repositories"
	"gorm.io/gorm"
)

type TenantLicenceService struct {
	tenantLicenceRepository *repositories.TenantLicenceRepository
}

func NewTenantLicenceService(tenantLicenceRepository *repositories.TenantLicenceRepository) *TenantLicenceService {
	return &TenantLicenceService{tenantLicenceRepository: tenantLicenceRepository}
}

func (s *TenantLicenceService) GetTenantLicenceByID(tenantID uint) (*dto.TenantLicenceResponseDTO, error) {
	tenantLicence, err := s.tenantLicenceRepository.GetByID(tenantID)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, config.ErrTenantLicenceNotFound
	} else if err != nil {
		return nil, err
	}

	return &dto.TenantLicenceResponseDTO{
		ID:            tenantLicence.ID,
		TenantID:      tenantLicence.TenantID,
		LicenceKey:    tenantLicence.LicenceKey,
		LicencedSeats: tenantLicence.LicencedSeats,
		UsedSeats:     tenantLicence.UsedSeats,
		ExpiryDate:    tenantLicence.ExpiryDate,
	}, nil
}

func (s *TenantLicenceService) GetAllTenantLicences() ([]entities.TenantLicence, error) {
	return s.tenantLicenceRepository.GetAll()
}

func (s *TenantLicenceService) GetTenantLicenceByLicenceKey(licenceKey string) (*dto.TenantLicenceResponseDTO, error) {
	tenantLicence, err := s.tenantLicenceRepository.GetByLicenceKey(licenceKey)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, config.ErrTenantLicenceNotFound
	} else if err != nil {
		return nil, err
	}
	return &dto.TenantLicenceResponseDTO{
		ID: tenantLicence.ID,
	}, nil
}

func (s *TenantLicenceService) GetTenantLicenceByTenantID(tenantID uint) (*dto.TenantLicenceResponseDTO, error) {
	tenantLicence, err := s.tenantLicenceRepository.GetByTenantID(tenantID)
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, config.ErrTenantLicenceNotFound
	} else if err != nil {
		return nil, err
	}
	return &dto.TenantLicenceResponseDTO{
		ID: tenantLicence.ID,
	}, nil
}

func (s *TenantLicenceService) UpdateTenantLicence(tenantID uint, tenantLicence *dto.TenantLicenceUpdateRequestDTO) error {
	existingTenantLicence, err := s.tenantLicenceRepository.GetByID(tenantID)
	if err != nil && err == gorm.ErrRecordNotFound {
		return config.ErrTenantLicenceNotFound
	} else if err != nil {
		return err
	}

	existingTenantLicence.UpdatedAt = time.Now()
	existingTenantLicence.LicenceKey = tenantLicence.LicenceKey
	existingTenantLicence.LicencedSeats = tenantLicence.LicencedSeats
	existingTenantLicence.ExpiryDate = tenantLicence.ExpiryDate

	return s.tenantLicenceRepository.Update(existingTenantLicence)
}
