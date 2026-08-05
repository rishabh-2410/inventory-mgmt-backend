package productstock

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	dbmodels "inventory-backend/src/models/db_object"
	reqmodels "inventory-backend/src/models/req_dto"

	"github.com/google/uuid"
)

type ProductStockRepo struct {
	db *sql.DB
}

func NewProductStockRepo(db *sql.DB) *ProductStockRepo {
	return &ProductStockRepo{db: db}
}

func (r *ProductStockRepo) GetInventoryList(businessID uuid.UUID, requestWithParams *reqmodels.GetInventoryListRequestDTO, offset int) ([]*dbmodels.ProductStock, error) {

	args := []any{businessID}
	argPos := 2
	conditions := []string{
		" ps.business_id = $1 ",
	}
	limit := "" // default

	if requestWithParams.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf(" ps.warehouse_id = $%d", argPos))
		args = append(args, *requestWithParams.WarehouseID)
		argPos++
	}

	if requestWithParams.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf(" ps.product_id = $%d", argPos))
		args = append(args, *requestWithParams.ProductID)
		argPos++
	}

	if requestWithParams.Search != nil {

		conditions = append(conditions, fmt.Sprintf(`
		(
    		LOWER(p.name) LIKE LOWER($%d)
    		OR LOWER(p.sku) LIKE LOWER($%d)
		)
			`, argPos, argPos))
		args = append(args, "%"+*requestWithParams.Search+"%")
		argPos++
	}

	if requestWithParams.Limit != nil {
		limit = fmt.Sprintf(" LIMIT $%d OFFSET $%d ", argPos, argPos+1)
		args = append(args, *requestWithParams.Limit, offset)
		argPos += 2
	}

	conditionString := strings.Join(conditions, " AND ")
	slog.Info("condition string", "conditionString", conditionString)

	query := `
		SELECT ps.id, p.name, p.sku, c.name, w.name, ps.current_stock, ps.updated_at
		FROM product_stock ps
		JOIN products p ON ps.product_id = p.id
		JOIN categories c ON p.category_id = c.id
		JOIN warehouses w ON ps.warehouse_id = w.id
		WHERE ` + strings.Join(conditions, " AND ") +
		limit + ";"
	slog.Debug("query", "query", query)

	productStockList := make([]*dbmodels.ProductStock, 0)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var productStock dbmodels.ProductStock
		err = rows.Scan(
			&productStock.ID,
			&productStock.ProductName,
			&productStock.Sku,
			&productStock.CategoryName,
			&productStock.WarehouseName,
			&productStock.CurrentStock,
			&productStock.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		productStockList = append(productStockList, &productStock)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to fetch product stock list", "error", err)
		return nil, err
	}

	slog.Info("fetched product stock list successfully for", "business_id", businessID)
	return productStockList, nil
}
