package repository

import (
	"github.com/geekible-ltd/auth-server/internal/models"
	"gorm.io/gorm"
)

type TenantLicenceRepository struct {
	db *gorm.DB
}

func NewTenantLicenceRepository(db *gorm.DB) *TenantLicenceRepository {
	return &TenantLicenceRepository{db: db}
}

func (r *TenantLicenceRepository) Create(tenantLicence *models.TenantLicence) error {
	return r.db.Create(tenantLicence).Error
}

func (r *TenantLicenceRepository) GetByID(tenantID uint) (*models.TenantLicence, error) {
	var tenantLicence models.TenantLicence
	if err := r.db.First(&tenantLicence, "tenant_id = ?", tenantID).Error; err != nil {
		return nil, err
	}
	return &tenantLicence, nil
}

func (r *TenantLicenceRepository) Update(tenantLicence *models.TenantLicence) error {
	return r.db.Save(tenantLicence).Error
}

func (r *TenantLicenceRepository) Delete(tenantLicence *models.TenantLicence) error {
	return r.db.Delete(tenantLicence).Error
}

func (r *TenantLicenceRepository) GetAll() ([]models.TenantLicence, error) {
	var tenantLicences []models.TenantLicence
	if err := r.db.Find(&tenantLicences).Error; err != nil {
		return nil, err
	}
	return tenantLicences, nil
}

func (r *TenantLicenceRepository) GetByLicenceKey(licenceKey string) (*models.TenantLicence, error) {
	var tenantLicence models.TenantLicence
	if err := r.db.First(&tenantLicence, "licence_key = ?", licenceKey).Error; err != nil {
		return nil, err
	}
	return &tenantLicence, nil
}

func (r *TenantLicenceRepository) GetByTenantID(tenantID uint) (*models.TenantLicence, error) {
	var tenantLicence models.TenantLicence
	if err := r.db.First(&tenantLicence, "tenant_id = ?", tenantID).Error; err != nil {
		return nil, err
	}
	return &tenantLicence, nil
}
