package repos

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/sweetloveinyourheart/planning-poker/pkg/db"
	"github.com/sweetloveinyourheart/planning-poker/services/user/models"
)

type UserCredentialRepository struct {
	Tx db.DbOrTx
}

func NewUserCredentialRepository(tx db.DbOrTx) IUserCredentialRepository {
	return &UserCredentialRepository{
		Tx: tx,
	}
}

func (repo *UserCredentialRepository) CreateCredential(ctx context.Context, userCredential *models.UserCredential) error {
	if err := userCredential.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO user_credentials (
			user_id, 
			auth_provider, 
			meta,
			created_at,
			updated_at
		) 
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err := repo.Tx.Exec(ctx,
		query,
		userCredential.UserID,
		userCredential.AuthProvider,
		userCredential.Meta,
		userCredential.CreatedAt,
		userCredential.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *UserCredentialRepository) GetUserCredentials(ctx context.Context, userID uuid.UUID) ([]models.UserCredential, error) {
	query := `
	SELECT 
		user_id, 
		auth_provider, 
		meta, 
		created_at, 
		updated_at 
	FROM 
		user_credentials 
	WHERE 
		user_id = $1;
`

	rows, err := repo.Tx.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userCredentials []models.UserCredential
	for rows.Next() {
		var userCredential models.UserCredential
		if err := rows.Scan(
			&userCredential.UserID,
			&userCredential.AuthProvider,
			&userCredential.Meta,
			&userCredential.CreatedAt,
			&userCredential.UpdatedAt,
		); err != nil {
			return nil, err
		}
		userCredentials = append(userCredentials, userCredential)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userCredentials, nil
}
