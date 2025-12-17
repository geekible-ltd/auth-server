package repositories

import (
	"github.com/geekible-ltd/auth-server/src/internal/entities"
	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(tenant *entities.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *TenantRepository) GetByID(tenantId uint) (*entities.Tenant, error) {
	var tenant entities.Tenant
	if err := r.db.First(&tenant, "id = ?", tenantId).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepository) Update(tenant *entities.Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *TenantRepository) Delete(tenant *entities.Tenant) error {
	return r.db.Delete(tenant).Error
}

func (r *TenantRepository) GetAll() ([]entities.Tenant, error) {
	var tenants []entities.Tenant
	if err := r.db.Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *TenantRepository) GetAllWithUsers(tenantId uint) ([]entities.Tenant, error) {
	var tenants []entities.Tenant
	if err := r.db.Preload("Users").Find(&tenants, "tenant_id = ?", tenantId).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *TenantRepository) GetByEmailDomain(emailDomain string) (*entities.Tenant, error) {
	var tenant entities.Tenant
	if err := r.db.Where("email LIKE ?", "%"+emailDomain).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}
