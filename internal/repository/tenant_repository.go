package repository

import (
	"github.com/geekible-ltd/auth-server/internal/models"
	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(tenant *models.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *TenantRepository) GetByID(tenantId uint) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.db.First(&tenant, "id = ?", tenantId).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) Update(tenant *models.Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *TenantRepository) Delete(tenant *models.Tenant) error {
	return r.db.Delete(tenant).Error
}

func (r *TenantRepository) GetAll() ([]models.Tenant, error) {
	var tenants []models.Tenant
	if err := r.db.Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *TenantRepository) GetAllWithUsers(tenantId uint) ([]models.Tenant, error) {
	var tenants []models.Tenant
	if err := r.db.Preload("Users").Find(&tenants, "tenant_id = ?", tenantId).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *TenantRepository) GetByEmailDomain(emailDomain string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.db.Where("email LIKE ?", "%"+emailDomain).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}
