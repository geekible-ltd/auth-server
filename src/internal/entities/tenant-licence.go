package entities

import "time"

type TenantLicence struct {
	ID            uint       `json:"id"`
	TenantID      uint       `json:"tenant_id"`
	LicenceKey    string     `json:"licence_key"`
	LicencedSeats int        `json:"licenced_seats"`
	UsedSeats     int        `json:"used_seats"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`

	Tenant Tenant `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
