package models

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
)

type UserSession struct {
	SessionID         int64      `json:"session_id"`
	UserID            uuid.UUID  `json:"user_id"`
	Token             string     `json:"token"`
	SessionStart      time.Time  `json:"session_start"`
	LastUpdated       time.Time  `json:"last_updated"`
	SessionExpiration *time.Time `json:"session_expiration,omitempty"` // Pointer to handle null
	SessionEnd        *time.Time `json:"session_end,omitempty"`
}

func (us UserSession) Validate() error {
	if us.SessionID == 0 || us.SessionID < 0 {
		return errors.New("SessionID: invalid session")
	}

	if us.UserID == uuid.Nil {
		return errors.New("UserID: nil")
	}

	if us.Token == "" {
		return errors.New("Token: blank")
	}

	if us.SessionStart.IsZero() {
		return errors.New("SessionStart: blank")
	}

	if us.LastUpdated.IsZero() {
		return errors.New("LastUpdated: blank")
	}

	if us.SessionExpiration != nil && us.SessionExpiration.IsZero() {
		return errors.New("SessionExpiration: blank")
	}

	if us.SessionEnd != nil && us.SessionEnd.IsZero() {
		return errors.New("SessionEnd: blank")
	}

	return nil
}
