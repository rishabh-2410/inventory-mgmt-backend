package repository

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (tr *TokenRepo) SaveRefreshToken(userID uuid.UUID, refreshToken string) error {
	slog.Debug("token repository save refresh token", "user_id", userID, "refres_token", refreshToken)
	result, err := tr.db.Exec(
		`INSERT INTO refresh_tokens(
    	user_id,
   	 	token_hash,
		expires_at
	) VALUES (
	 	$1,
		$2,
		NOW() + INTERVAL '30 days'
	) 
	`, userID, refreshToken)

	if err != nil {
		slog.Error("token repository save refresh token failed user_id=%s err=%v", userID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("token repository save refresh token rows failed", "user_id", userID, "err", err)
		return err
	}

	if rowsAffected == 0 {
		slog.Error("token repository save refresh token no rows affected", "user_id", userID)
		return errors.New("user already is already logged in")
	}

	slog.Info("token repository save refresh token success", "user_id", userID)

	return nil
}

func (tr *TokenRepo) MarkTokenExpired(refreshToken string) (uuid.UUID, error) {
	slog.Debug("token repository mark token expired", "refresh_token", refreshToken)
	var userID uuid.UUID
	err := tr.db.QueryRow(`
		UPDATE refresh_tokens 
		SET revoked_at=NOW()
		WHERE token_hash = $1
			AND expires_at > NOW()
  			AND revoked_at IS NULL
		RETURNING user_id
		`, refreshToken).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("token repository mark token expired not found", "refresh_token", refreshToken)
			return uuid.UUID{}, errors.New("refresh token not found")
		}
		slog.Error("token repository mark token expired failed", "err", err)
		return uuid.UUID{}, err
	}

	slog.Info("token repository mark token expired success", "user_id", userID)

	return userID, nil
}
