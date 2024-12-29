package models

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
)

const (
	USER_STATUS_ENABLED  = 1
	USER_STATUS_DISABLED = 2
)

type User struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) Validate() error {
	if u.UserID == uuid.Nil {
		return errors.New("UserID: nil")
	}

	if u.Username == "" {
		return errors.New("Username: blank")
	}

	if u.FullName == "" {
		return errors.New("FullName: blank")
	}

	if u.Status < 0 || u.Status > 2 {
		return errors.New("Status: not a valid status")
	}

	if u.CreatedAt.IsZero() {
		return errors.New("CreatedAt: zero")
	}

	if u.UpdatedAt.IsZero() {
		return errors.New("UpdatedAt: zero")
	}

	return nil
}
