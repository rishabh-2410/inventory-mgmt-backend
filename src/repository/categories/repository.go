package repository

import (
	"database/sql"
	"fmt"
	dbmodel "inventory-backend/src/models/db_object"
	reqmodels "inventory-backend/src/models/req_dto"
	"log/slog"

	"github.com/google/uuid"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

func (cr *CategoryRepo) AddCategory(businessID uuid.UUID, category *reqmodels.AddCategoryRequestDTO) (*dbmodel.Category, error) {
	result := &dbmodel.Category{}

	err := cr.db.QueryRow(
		`INSERT INTO categories(
			name,
			business_id,
			description
		) 
		 VALUES(
			$1,
			$2,
			$3
		 ) 
		 RETURNING id, name, business_id, description, created_at, updated_at
		`, category.Name, businessID, category.Description).Scan(
		&result.ID,
		&result.Name,
		&result.BusinessID,
		&result.Description,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		slog.Error("failed to add category", "error", err)
		return nil, err
	}

	return result, nil

}

func (cr *CategoryRepo) GetCategories(businessID uuid.UUID) ([]dbmodel.Category, error) {
	result := make([]dbmodel.Category, 0)

	rows, err := cr.db.Query(
		`SELECT * FROM categories
		WHERE business_id=$1	
		`, businessID)

	if err != nil {
		slog.Error("failed to fetch categories", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var category dbmodel.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.BusinessID,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			slog.Error("failed to fetch categories", "error", err)
			return nil, err
		}

		result = append(result, category)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to fetch categories", "error", err)
		return nil, err
	}

	slog.Info("fetched all categories successfully for", "business_id", businessID)

	return result, nil

}

func (cr *CategoryRepo) UpdateCategory(categoryID uuid.UUID, businessID uuid.UUID, category *reqmodels.UpdateCategoryRequestDTO) (*dbmodel.Category, error) {
	result := &dbmodel.Category{}

	err := cr.db.QueryRow(
		`UPDATE categories
		 SET 
		 	name=COALESCE($1, name),
			description=COALESCE($2, description)
		 WHERE id=$3 AND business_id=$4
		 RETURNING id, name, business_id, description, created_at, updated_at
		`, category.Name, category.Description, categoryID, businessID).Scan(
		&result.ID,
		&result.Name,
		&result.BusinessID,
		&result.Description,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		slog.Error("failed to add category", "error", err)
		return nil, err
	}

	return result, nil

}

func (cr *CategoryRepo) DeleteCategory(categoryID uuid.UUID, businessID uuid.UUID) error {

	result, err := cr.db.Exec(
		`DELETE FROM categories
		WHERE id=$1 AND business_id=$2
		`, categoryID, businessID)

	if err != nil {
		slog.Error("failed to add category", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to delete category", "error", err)
		return err
	}

	if rowsAffected == 0 {
		slog.Error("unable to delete category: no rows affected")
		return fmt.Errorf("unable to delete category: no rows affected")
	}

	return nil

}
