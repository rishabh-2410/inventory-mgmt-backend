package repository

import (
	"database/sql"
	"errors"
	dbmodel "inventory-backend/src/models/db_object"

	"log/slog"

	"github.com/google/uuid"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) FindByEmail(email string) (*dbmodel.User, error) {
	slog.Debug("finding user in db...", "email", email)
	var user dbmodel.User
	err := ur.db.QueryRow(`
		SELECT * FROM users WHERE email=$1
	`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.BusinessID,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		slog.Error("unable to find user by email", "email", email, "err", err)
		return nil, err
	}

	slog.Debug("user found by email", "email", email)

	return &user, nil
}

func (ur *UserRepo) FindByID(userID uuid.UUID) (*dbmodel.User, error) {
	slog.Debug("finding user by id", "id", userID)
	var user dbmodel.User
	err := ur.db.QueryRow(`
		SELECT * FROM users WHERE id=$1
	`, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.BusinessID,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		slog.Error("unable to find user by id", "id", userID, "err", err)
		return nil, err
	}

	slog.Debug("user found by id", "id", userID)

	return &user, nil
}

func (ur *UserRepo) AddUser(user *dbmodel.User) (uuid.UUID, error) {
	var userID uuid.UUID
	slog.Debug("adding user...", "email", user.Email, "name", user.Name)
	err := ur.db.QueryRow(`INSERT INTO users(
    name,
    email,
    password_hash,
	role
	) VALUES ($1,$2,$3, $4)
	RETURNING id
	`, user.Name, user.Email, user.PasswordHash, user.Role).Scan(&userID)
	if err != nil {
		slog.Error("failed to add new user", "email", user.Email, "err", err)
		return uuid.Nil, err
	}

	return userID, nil

}

func (ur *UserRepo) RevokeUser(hashedToken string) error {
	slog.Debug("user repository revoke user", "hashed_refresh_token", hashedToken)

	result, err := ur.db.Exec(
		`UPDATE refresh_tokens
			SET revoked_at = NOW()
		WHERE token_hash = $1 
			AND revoked_at IS NULL 
			AND expires_at > NOW()
		`, hashedToken)
	if err != nil {
		slog.Error("user repository revoke user exec failed err=%v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("user repository revoke user rows failed err=%v", err)
		return err
	}

	if rowsAffected == 0 {
		slog.Error("user repository revoke user no rows affected")
		return errors.New("refresh token not found or already revoked")
	}

	slog.Info("user repository revoke user success")

	return nil
}
