package repositories

import (
	"github.com/geekible-ltd/auth-server/src/internal/entities"
	"gorm.io/gorm"
)

type TenantLicenceRepository struct {
	db *gorm.DB
}

func NewTenantLicenceRepository(db *gorm.DB) *TenantLicenceRepository {
	return &TenantLicenceRepository{db: db}
}

func (r *TenantLicenceRepository) Create(tenantLicence *entities.TenantLicence) error {
	return r.db.Create(tenantLicence).Error
}

func (r *TenantLicenceRepository) GetByID(tenantID uint) (*entities.TenantLicence, error) {
	var tenantLicence entities.TenantLicence
	if err := r.db.First(&tenantLicence, "tenant_id = ?", tenantID).Error; err != nil {
		return nil, err
	}
	return &tenantLicence, nil
}

func (r *TenantLicenceRepository) Update(tenantLicence *entities.TenantLicence) error {
	return r.db.Save(tenantLicence).Error
}

func (r *TenantLicenceRepository) Delete(tenantLicence *entities.TenantLicence) error {
	return r.db.Delete(tenantLicence).Error
}

func (r *TenantLicenceRepository) GetAll() ([]entities.TenantLicence, error) {
	var tenantLicences []entities.TenantLicence
	if err := r.db.Find(&tenantLicences).Error; err != nil {
		return nil, err
	}
	return tenantLicences, nil
}

func (r *TenantLicenceRepository) GetByLicenceKey(licenceKey string) (*entities.TenantLicence, error) {
	var tenantLicence entities.TenantLicence
	if err := r.db.First(&tenantLicence, "licence_key = ?", licenceKey).Error; err != nil {
		return nil, err
	}
	return &tenantLicence, nil
}

func (r *TenantLicenceRepository) GetByTenantID(tenantID uint) (*entities.TenantLicence, error) {
	var tenantLicence entities.TenantLicence
	if err := r.db.First(&tenantLicence, "tenant_id = ?", tenantID).Error; err != nil {
		return nil, err
	}
	return &tenantLicence, nil
}
