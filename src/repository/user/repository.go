package repository

import (
	"database/sql"
	dbmodel "inventory-backend/src/models/db_object"

	"log/slog"
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

func (ur *UserRepo) AddUser(user *dbmodel.User) error {
	slog.Debug("adding user...", "email", user.Email, "name", user.Name)
	err := ur.db.QueryRow(`INSERT INTO users(
    name,
    email,
    password_hash,
	role
	) VALUES ($1,$2,$3, $4)
	`, user.Name, user.Email, user.PasswordHash, user.Role).Err()

	if err != nil {
		slog.Error("failed to add new user", "email", user.Email, "err", err)
		return err
	}

	slog.Info("user added successfully", "email", user.Email)

	return nil

}
