package repositories

import (
	"github.com/geekible-ltd/auth-server/src/internal/entities"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(userId, tenantId uint) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, "user_id = ? AND tenant_id = ?", userId, tenantId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(user *entities.User) error {
	return r.db.Delete(user).Error
}

func (r *UserRepository) GetAll(tenantId uint) ([]entities.User, error) {
	var users []entities.User
	if err := r.db.Find(&users, "tenant_id = ?", tenantId).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetAllWithTenant(tenantId uint) ([]entities.User, error) {
	var users []entities.User
	if err := r.db.Preload("Tenant").Find(&users, "tenant_id = ?", tenantId).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetByEmailDomain(emailDomain string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email LIKE ?", "%"+emailDomain).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
