package repository

import (
	"database/sql"
	"fmt"
	dbmodel "inventory-backend/src/models/db_object"
	reqmodel "inventory-backend/src/models/req_dto"
	"log/slog"

	"github.com/google/uuid"
)

type BusinessRepo struct {
	db *sql.DB
}

func NewBusinessRepo(db *sql.DB) *BusinessRepo {
	return &BusinessRepo{
		db: db,
	}
}

func (br *BusinessRepo) FindBusinessByName(businessName string) (*dbmodel.Businesss, error) {
	result := &dbmodel.Businesss{}

	err := br.db.QueryRow(
		`SELECT * FROM business
			WHERE name=$1
		`, businessName).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Phone,
		&result.Address,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		slog.Warn("unable to find business in db", "name", businessName, "err", err)
		return nil, err
	}

	return result, nil
}

func (br *BusinessRepo) CreateBusiness(userID uuid.UUID, business *reqmodel.CreateBusinessRequestDTO) (*dbmodel.Businesss, error) {

	tx, err := br.db.Begin()
	if err != nil {
		slog.Error("unable to begin transaction", "err", err)
		return nil, err
	}

	defer tx.Rollback()

	result := &dbmodel.Businesss{}
	err = tx.QueryRow(
		`INSERT INTO business(
			name,
			email,
			address, 
			phone
		) VALUES(
				$1,
				$2,
				$3,
				$4
		)
		 RETURNING
			id,
			name,
			email,
			phone,
			address,
			created_at,
			updated_at
		`,
		business.Name,
		business.Email,
		business.Address,
		business.Phone,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Phone,
		&result.Address,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		slog.Error("unable to create business", "name", business.Name, "err", err)
		return nil, err
	}

	updateResult, err := tx.Exec(
		`UPDATE users
		 SET 
		   business_id=$1
		 WHERE 
		   id=$2 
		   AND business_id IS NULL 
		`, result.ID, userID)
	if err != nil {
		slog.Error("unable to assign business to user", "user_id", userID, "err", err)
		return nil, err
	}

	rowsAffected, err := updateResult.RowsAffected()
	if err != nil {
		slog.Error("unable to get rows affected", "err", err)
	}

	if rowsAffected != 1 {
		return nil, fmt.Errorf("user not found or already belongs to a business")
	}

	if err := tx.Commit(); err != nil {
		slog.Error("unable to commit transaction", "err", err)
		return nil, err
	}

	return result, nil
}
