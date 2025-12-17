package dto

import "time"

type TenantLicenceResponseDTO struct {
	ID            uint       `json:"id"`
	TenantID      uint       `json:"tenant_id"`
	LicenceKey    string     `json:"licence_key"`
	LicencedSeats int        `json:"licenced_seats"`
	UsedSeats     int        `json:"used_seats"`
	ExpiryDate    *time.Time `json:"expiry_date"`
}

type TenantLicenceUpdateRequestDTO struct {
	LicenceKey    string     `json:"licence_key"`
	LicencedSeats int        `json:"licenced_seats"`
	ExpiryDate    *time.Time `json:"expiry_date"`
}
