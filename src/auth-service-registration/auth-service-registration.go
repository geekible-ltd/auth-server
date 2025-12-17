package authserviceregistration

import (
	"github.com/geekible-ltd/auth-server/src/internal/entities"
	"gorm.io/gorm"
)

type AuthServerRegistration struct {
	db *gorm.DB
}

func NewAuthServerRegistration(db *gorm.DB) *AuthServerRegistration {
	return &AuthServerRegistration{db: db}
}

func (r *AuthServerRegistration) MigrateDBModels() error {
	err := r.db.AutoMigrate(&entities.User{}, &entities.Tenant{}, &entities.TenantLicence{})
	if err != nil {
		return err
	}
	return nil
}
