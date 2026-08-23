package repository

import (
	"database/sql"
	"fmt"
	dbmodel "inventory-backend/src/models/db_object"
	reqmodels "inventory-backend/src/models/req_dto"
	"log/slog"

	"github.com/google/uuid"
)

type WarehouseRepo struct {
	db *sql.DB
}

func NewWarehouseRepo(db *sql.DB) *WarehouseRepo {
	return &WarehouseRepo{
		db: db,
	}
}

func (wr *WarehouseRepo) AddWarehouse(businessID uuid.UUID, warehouse *reqmodels.AddWarehouseRequestDTO) (*dbmodel.Warehouse, error) {
	result := &dbmodel.Warehouse{}

	err := wr.db.QueryRow(
		`INSERT INTO warehouses(
			name,
			business_id,
			address
		) 
		 VALUES(
			$1,
			$2,
			$3
		 ) 
		 RETURNING id, business_id, name, address, created_at, updated_at
		`, warehouse.Name, businessID, warehouse.Address).Scan(
		&result.ID,
		&result.BusinessID,
		&result.Name,
		&result.Address,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		slog.Error("failed to add warehouse", "error", err)
		return nil, err
	}

	return result, nil

}

func (wr *WarehouseRepo) GetWarehouses(businessID uuid.UUID) ([]dbmodel.Warehouse, error) {
	result := make([]dbmodel.Warehouse, 0)

	rows, err := wr.db.Query(
		`SELECT id, business_id, name, address, created_at, updated_at 
		FROM warehouses
		WHERE business_id=$1	
		`, businessID)

	if err != nil {
		slog.Error("failed to fetch warehouses", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var warehouse dbmodel.Warehouse

		err := rows.Scan(
			&warehouse.ID,
			&warehouse.BusinessID,
			&warehouse.Name,
			&warehouse.Address,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
		)

		if err != nil {
			slog.Error("failed to fetch warehouses", "error", err)
			return nil, err
		}

		result = append(result, warehouse)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to fetch warehouses", "error", err)
		return nil, err
	}

	slog.Info("fetched all warehouses successfully for", "business_id", businessID)

	return result, nil

}

func (wr *WarehouseRepo) UpdateWarehouse(warehouseID uuid.UUID, businessID uuid.UUID, warehouse *reqmodels.EditWarehouseRequestDTO) (*dbmodel.Warehouse, error) {
	result := &dbmodel.Warehouse{}

	err := wr.db.QueryRow(
		`UPDATE warehouses
		 SET 
		 	name=COALESCE($1, name),
			address=COALESCE($2, address)
		 WHERE id=$3 AND business_id=$4
		 RETURNING id, business_id, name, address, created_at, updated_at
		`, warehouse.Name, warehouse.Address, warehouseID, businessID).Scan(
		&result.ID,
		&result.BusinessID,
		&result.Name,
		&result.Address,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		slog.Error("failed to update warehouse", "error", err)
		return nil, err
	}

	return result, nil

}

func (wr *WarehouseRepo) DeleteWarehouse(warehouseID uuid.UUID, businessID uuid.UUID) error {

	result, err := wr.db.Exec(
		`DELETE FROM warehouses
		WHERE id=$1 AND business_id=$2
		`, warehouseID, businessID)

	if err != nil {
		slog.Error("failed to delete warehouse", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to delete warehouse", "error", err)
		return err
	}

	if rowsAffected == 0 {
		slog.Error("unable to delete warehouse: no rows affected")
		return fmt.Errorf("unable to delete warehouse: no rows affected")
	}

	return nil

}
