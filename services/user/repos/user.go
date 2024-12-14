package repos

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/sweetloveinyourheart/planning-poker/pkg/db"
	"github.com/sweetloveinyourheart/planning-poker/services/user/models"
)

type UserRepository struct {
	Tx db.DbOrTx
}

func NewUserRepository(tx db.DbOrTx) IUserRepository {
	return &UserRepository{
		Tx: tx,
	}
}

func (ur *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, bool, error) {
	var user models.User

	query := `
		SELECT user_id, first_name, last_name, status, created_at, updated_at
		FROM users
		WHERE user_id = $1;
	`
	err := ur.Tx.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, false, nil
		}

		return models.User{}, false, errors.WithStack(err)
	}

	return user, true, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO users (
			user_id,
			first_name,
			last_name,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := ur.Tx.Exec(ctx,
		query,
		user.UserID,
		user.FirstName,
		user.LastName,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return errors.WithStack(err)
}

func (ur *UserRepository) UpdateUserData(ctx context.Context, user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, status = $3, updated_at = $4
		WHERE user_id = $5;
	`
	_, err := ur.Tx.Exec(ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Status,
		user.UpdatedAt,
		user.UserID,
	)

	return errors.WithStack(err)
}
