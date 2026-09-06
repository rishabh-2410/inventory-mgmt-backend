package repository

import (
	"database/sql"
	"errors"
	dbmodel "inventory-backend/src/models/db_object"
	enum "inventory-backend/src/models/enums"

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

func (ur *UserRepo) AddEmployee(employee *dbmodel.User) (*dbmodel.User, error) {
	newEmployee := &dbmodel.User{}
	slog.Debug("adding employee...", "email", employee.Email, "name", employee.Name)
	err := ur.db.QueryRow(`INSERT INTO users(
		name,
		email,
		business_id,
		password_hash,
		role,
		is_active
	) VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id, name, email, business_id, password_hash, role, is_active, created_at, updated_at
	`, employee.Name, employee.Email, employee.BusinessID, employee.PasswordHash, employee.Role, employee.IsActive).Scan(
		&newEmployee.ID,
		&newEmployee.Name,
		&newEmployee.Email,
		&newEmployee.BusinessID,
		&newEmployee.PasswordHash,
		&newEmployee.Role,
		&newEmployee.IsActive,
		&newEmployee.CreatedAt,
		&newEmployee.UpdatedAt,
	)
	if err != nil {
		slog.Error("failed to add new employee", "email", employee.Email, "err", err)
		return nil, err
	}

	return newEmployee, nil
}

func (ur *UserRepo) FindByBusinessIDAndRole(businessID uuid.UUID, role enum.UserRole) ([]dbmodel.User, error) {
	slog.Debug("finding users by business id and role", "business_id", businessID, "role", role)

	rows, err := ur.db.Query(`
	SELECT id, name, email, business_id, password_hash, role, is_active, created_at, updated_at
	FROM users 
	WHERE business_id=$1 AND role=$2
	`, businessID, role)

	if err != nil {
		slog.Error("failed to find users by business id and role", "business_id", businessID, "role", role, "err", err)
		return nil, err
	}

	defer rows.Close()

	usersList := make([]dbmodel.User, 0)

	for rows.Next() {
		var user dbmodel.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.BusinessID, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			slog.Error("failed to scan user", "err", err)
			return nil, err
		}
		usersList = append(usersList, user)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to fetch users", "error", err)
		return nil, err
	}

	slog.Debug("users found by business id and role", "business_id", businessID, "role", role)

	return usersList, nil
}

func (ur *UserRepo) DeleteOwnerAndBusiness(businessID uuid.UUID) error {
	slog.Debug("deleting owner and business...", "business_id", businessID)

	tx, err := ur.db.Begin()
	if err != nil {
		slog.Error("failed to begin transaction", "err", err)
		return err
	}
	defer tx.Rollback()

	deleteBusinessQuery := `DELETE FROM businesses WHERE id=$1`
	_, err = tx.Exec(deleteBusinessQuery, businessID)
	if err != nil {
		slog.Error("failed to delete business", "business_id", businessID, "err", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("failed to commit transaction", "err", err)
		return err
	}

	slog.Info("owner and business deleted successfully", "business_id", businessID)

	return nil
}

func (ur *UserRepo) DeleteEmployee(employeeID string) error {
	slog.Debug("deleting employee...", "employee_id", employeeID)

	_, err := ur.db.Exec(`UPDATE users SET is_active = false WHERE id=$1`, employeeID)
	if err != nil {
		slog.Error("failed to deactivate employee", "employee_id", employeeID, "err", err)
		return err
	}

	slog.Info("employee deactivated successfully", "employee_id", employeeID)

	return nil
}
