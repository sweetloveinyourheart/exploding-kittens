package repos

import (
	"context"

	"github.com/sweetloveinyourheart/planning-poker/pkg/db"
	"github.com/sweetloveinyourheart/planning-poker/services/user/models"
)

type UserSessionRepository struct {
	Tx db.DbOrTx
}

func NewUserSessionRepository(tx db.DbOrTx) IUserSessionRepository {
	return &UserSessionRepository{
		Tx: tx,
	}
}

func (repo *UserSessionRepository) CreateSession(ctx context.Context, userSession *models.UserSession) error {
	if err := userSession.Validate(); err != nil {
		return err
	}

	query := `
        INSERT INTO user_sessions (
            user_id,
            token,
            session_start,
            last_updated,
            session_expiration,
            session_end
        )
        VALUES ($1, $2, $3, $4, $5, $6);
    `
	_, err := repo.Tx.Exec(ctx,
		query,
		userSession.UserID,
		userSession.Token,
		userSession.SessionStart,
		userSession.LastUpdated,
		userSession.SessionExpiration,
		userSession.SessionEnd,
	)

	return err
}

func (repo *UserSessionRepository) GetUserSessionByToken(ctx context.Context, token string) (models.UserSession, error) {
	var userSession models.UserSession

	query := `
        SELECT session_id, user_id, token, session_start, last_updated, session_expiration, session_end
        FROM user_sessions
        WHERE token = $1;
    `
	err := repo.Tx.QueryRow(ctx, query, token).Scan(
		&userSession.SessionID,
		&userSession.UserID,
		&userSession.Token,
		&userSession.SessionStart,
		&userSession.LastUpdated,
		&userSession.SessionExpiration,
		&userSession.SessionEnd,
	)

	if err != nil {
		return models.UserSession{}, err
	}

	return userSession, nil
}
