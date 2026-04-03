package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	RefreshTokenHash  string    `json:"refresh_token_hash"`
	ExpiresAt         time.Time `json:"expires_at"`
	Revoked           bool      `json:"revoked"`
	DeviceInfo        *string   `json:"device_info"`
	IPAddr            *string   `json:"ip_addr"`
	CreatedAt         time.Time `json:"created_at"`
}
