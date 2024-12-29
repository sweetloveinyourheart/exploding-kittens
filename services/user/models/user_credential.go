package models

import (
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
)

const (
	USER_AUTH_PROVIDER_GUEST  = "USER_AUTH_PROVIDER_GUEST"
	USER_AUTH_PROVIDER_GOOGLE = "USER_AUTH_PROVIDER_GOOGLE"
)

type UserCredential struct {
	UserID       uuid.UUID `json:"user_id"`
	AuthProvider string    `json:"auth_provider"`
	Meta         []byte    `json:"meta"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (uc UserCredential) Validate() error {
	if uc.UserID == uuid.Nil {
		return errors.New("UserID: nil")
	}

	if uc.AuthProvider == "" {
		return errors.New("AuthProvider: blank")
	}

	if len(uc.Meta) > 0 && !json.Valid(uc.Meta) {
		return errors.New("Meta: invalid json")
	}

	if uc.CreatedAt.IsZero() {
		return errors.New("CreatedAt: zero")
	}

	if uc.UpdatedAt.IsZero() {
		return errors.New("UpdatedAt: zero")
	}

	return nil
}
