package model

import (
	"time"

	"github.com/google/uuid"
)

type AuthProvider string

const (
	ProviderEmail AuthProvider = "EMAIL"
	ProviderGmail AuthProvider = "GMAIL"
)

type AuthAccount struct {
	ID             uuid.UUID    `json:"id"`
	UserID         uuid.UUID    `json:"user_id"`
	Provider       AuthProvider `json:"provider"`
	ProviderUserID *string      `json:"provider_user_id"`
	PasswordHash   *string      `json:"password_hash"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
