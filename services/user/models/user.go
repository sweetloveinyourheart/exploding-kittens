package models

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) Validate() error {
	if u.UserID == uuid.Nil {
		return errors.New("UserID: nil")
	}

	if u.FirstName == "" {
		return errors.New("FirstName: blank")
	}

	if u.LastName == "" {
		return errors.New("LastName: blank")
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
