package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	dbmodels "inventory-backend/src/models/db_object"
	reqmodels "inventory-backend/src/models/req_dto"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// #region agent log
func debugLog(hypothesisId, location, message string, data map[string]interface{}) {
	payload := map[string]interface{}{
		"sessionId":    "d4a16d",
		"runId":        "pre-fix",
		"hypothesisId": hypothesisId,
		"location":     location,
		"message":      message,
		"data":         data,
		"timestamp":    time.Now().UnixMilli(),
	}
	b, _ := json.Marshal(payload)
	f, err := os.OpenFile("/Users/nikolai/Dev/work/inventory management/.cursor/debug-d4a16d.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(append(b, '\n'))
}

// #endregion

type StockMovementRepo struct {
	db *sql.DB
}

func NewStockMovementRepo(db *sql.DB) *StockMovementRepo {
	return &StockMovementRepo{db: db}
}

func (r *StockMovementRepo) UpdateStockMovement(userID uuid.UUID, businessID uuid.UUID, stockMovementReq *reqmodels.StockMovementRequestDTO) (float64, *dbmodels.StockMovement, error) {
	currentStock := float64(0)
	stockMovement := &dbmodels.StockMovement{}
	// Check if product stock exists for given product and warehouse
	productStockExists, err := r.validateProductStockExists(stockMovementReq.ProductID, stockMovementReq.WarehouseID)
	if err != nil {
		return float64(0), nil, err
	}

	if !productStockExists {
		if stockMovementReq.MovementType == "RECEIVE" {
			currentStock, stockMovement, err = r.createProductStockAndUpdateStockMovement(userID, businessID, stockMovementReq)
			if err != nil {
				return float64(0), nil, err
			}

		} else {
			return float64(0), nil, errors.New("product stock does not exist")
		}
	} else {
		currentStock, stockMovement, err = r.updateProductStockAndStockMovement(userID, businessID, stockMovementReq)
		if err != nil {
			return float64(0), nil, err
		}
	}

	return currentStock, stockMovement, nil
}

func (r *StockMovementRepo) GetStockMovementHistory(businessID uuid.UUID, requestWithParams *reqmodels.GetStockMovementHistoryRequestDTO, offset int) ([]*dbmodels.StockMovementHistory, error) {

	args := []any{businessID}
	argPos := 2
	conditions := []string{
		" sm.business_id = $1 ",
	}
	orderBy := " ORDER BY created_at DESC " // default
	limit := ""                             // default

	if requestWithParams.MovementType != nil {
		conditions = append(conditions, fmt.Sprintf(" movement_type = $%d", argPos))
		args = append(args, *requestWithParams.MovementType)
		argPos++
	}

	if requestWithParams.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf(" warehouse_id = $%d", argPos))
		args = append(args, *requestWithParams.WarehouseID)
		argPos++
	}

	if requestWithParams.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf(" product_id = $%d", argPos))
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

	if requestWithParams.From != nil {
		conditions = append(conditions, fmt.Sprintf(" created_at >= $%d", argPos))
		args = append(args, *requestWithParams.From)
		argPos++
	}

	if requestWithParams.To != nil {
		conditions = append(conditions, fmt.Sprintf(" created_at <= $%d", argPos))
		args = append(args, *requestWithParams.To)
		argPos++
	}

	if requestWithParams.SortBy != nil {
		switch *requestWithParams.SortBy {
		case "created_at_desc":
			orderBy = " ORDER BY created_at ASC "
		case "created_at_asc":
			orderBy = " ORDER BY created_at DESC "
		case "movement_quantity_desc":
			orderBy = " ORDER BY movement_quantity DESC "
		case "movement_quantity_asc":
			orderBy = " ORDER BY movement_quantity ASC "
		default:
			orderBy = " ORDER BY created_at DESC "
		}
	}

	if requestWithParams.Limit != nil {
		limit = fmt.Sprintf(" LIMIT $%d OFFSET $%d ", argPos, argPos+1)
		args = append(args, *requestWithParams.Limit, offset)
		argPos += 2
	}

	conditionString := strings.Join(conditions, " AND ")
	slog.Info("condition string", "conditionString", conditionString)

	query := `
		SELECT sm.user_id, sm.product_id, sm.movement_type, sm.movement_quantity, sm.reason, sm.created_at, p.name, p.sku, w.name
		FROM stock_movement sm 
		JOIN products p ON sm.product_id = p.id
		JOIN warehouses w ON sm.warehouse_id = w.id
		WHERE ` + strings.Join(conditions, " AND ") +
		orderBy +
		limit + ";"

	slog.Debug("query", "query", query)

	stockMovementsHistory := []*dbmodels.StockMovementHistory{}
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var stockMovementHistory dbmodels.StockMovementHistory
		err = rows.Scan(
			&stockMovementHistory.UserID,
			&stockMovementHistory.ProductID,
			&stockMovementHistory.MovementType,
			&stockMovementHistory.MovementQuantity,
			&stockMovementHistory.Reason,
			&stockMovementHistory.CreatedAt,
			&stockMovementHistory.ProductName,
			&stockMovementHistory.Sku,
			&stockMovementHistory.WarehouseName,
		)
		if err != nil {
			return nil, err
		}
		stockMovementsHistory = append(stockMovementsHistory, &stockMovementHistory)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to fetch stock movement history", "error", err)
		return nil, err
	}

	slog.Info("fetched stock movement history successfully for", "business_id", businessID)
	return stockMovementsHistory, nil
}

func (r *StockMovementRepo) GetStockMovementHistoryByWarehouse(businessID uuid.UUID, warehouseID uuid.UUID) ([]*dbmodels.StockMovementHistory, error) {
	query := `
		SELECT sm.user_id, sm.product_id, sm.movement_type, sm.movement_quantity, sm.reason, sm.created_at, p.name, p.sku, w.warehouse_name
		FROM stock_movement sm 
		JOIN products p ON sm.product_id = p.id
		JOIN warehouses w ON sm.warehouse_id = w.id
		WHERE sm.business_id = $1 AND sm.warehouse_id = $2
	`
	stockMovementsHistory := []*dbmodels.StockMovementHistory{}
	rows, err := r.db.Query(query, businessID, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var stockMovementHistory dbmodels.StockMovementHistory
		err = rows.Scan(
			&stockMovementHistory.UserID,
			&stockMovementHistory.ProductID,
			&stockMovementHistory.MovementType,
			&stockMovementHistory.MovementQuantity,
			&stockMovementHistory.Reason,
			&stockMovementHistory.CreatedAt,
			&stockMovementHistory.ProductName,
			&stockMovementHistory.Sku,
			&stockMovementHistory.WarehouseName)
		if err != nil {
			return nil, err
		}

		stockMovementsHistory = append(stockMovementsHistory, &stockMovementHistory)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to fetch stock movement history by product id", "error", err)
		return nil, err
	}

	slog.Info("fetched stock movement history by warehouse successfully for", "business_id", businessID)
	return stockMovementsHistory, nil
}

// ----------------------------- HELPERS----------------------------------
func (r *StockMovementRepo) validateProductStockExists(productID uuid.UUID, warehouseID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*) FROM product_stock
		WHERE product_id = $1 AND warehouse_id = $2
	`
	var count int
	err := r.db.QueryRow(query, productID, warehouseID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *StockMovementRepo) createProductStockAndUpdateStockMovement(userID uuid.UUID, businessID uuid.UUID, stockMovementReq *reqmodels.StockMovementRequestDTO) (float64, *dbmodels.StockMovement, error) {
	tx, err := r.db.Begin()
	if err != nil {
		slog.Error("unable to begin transaction", "err", err)
		return float64(0), nil, err
	}

	defer tx.Rollback()
	query := `
		INSERT INTO product_stock (business_id, product_id, warehouse_id, current_stock)
		VALUES ($1, $2, $3, $4)
		RETURNING id, current_stock
	`
	var productStockID uuid.UUID
	var currentStock float64
	err = tx.QueryRow(
		query,
		businessID,
		stockMovementReq.ProductID,
		stockMovementReq.WarehouseID,
		stockMovementReq.MovementQuantity,
	).Scan(
		&productStockID,
		&currentStock,
	)
	if err != nil {
		slog.Error("unable to create product stock", "err", err)
		return float64(0), nil, err
	}

	query = `
		INSERT INTO stock_movement (product_stock_id, user_id, business_id, product_id, warehouse_id, movement_type, movement_quantity, reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, product_stock_id, user_id, product_id, warehouse_id, business_id,movement_type, movement_quantity, reason, created_at
	`
	stockMovement := &dbmodels.StockMovement{}
	err = tx.QueryRow(
		query,
		productStockID,
		userID,
		businessID,
		stockMovementReq.ProductID,
		stockMovementReq.WarehouseID,
		stockMovementReq.MovementType,
		stockMovementReq.MovementQuantity,
		stockMovementReq.Reason,
	).Scan(
		&stockMovement.ID,
		&stockMovement.ProductStockID,
		&stockMovement.UserID,
		&stockMovement.ProductID,
		&stockMovement.WarehouseID,
		&stockMovement.BusinessID,
		&stockMovement.MovementType,
		&stockMovement.MovementQuantity,
		&stockMovement.Reason,
		&stockMovement.CreatedAt,
	)
	if err != nil {
		return float64(0), nil, err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("unable to commit transaction", "err", err)
		return 0, nil, err
	}
	return currentStock, stockMovement, nil
}

func (r *StockMovementRepo) updateProductStockAndStockMovement(userID uuid.UUID, businessID uuid.UUID, stockMovementReq *reqmodels.StockMovementRequestDTO) (float64, *dbmodels.StockMovement, error) {
	tx, err := r.db.Begin()
	if err != nil {
		slog.Error("unable to begin transaction", "err", err)
		return float64(0), nil, err
	}

	defer tx.Rollback()

	var query string
	var currentStock float64
	switch stockMovementReq.MovementType {
	case "RECEIVE", "RETURN":
		query = `
			UPDATE product_stock SET current_stock = current_stock + $1 WHERE product_id = $2 AND warehouse_id = $3
			RETURNING id, current_stock
		`
	case "SALE", "DAMAGE":
		query = `
			UPDATE product_stock SET current_stock = current_stock - $1 WHERE product_id = $2 AND warehouse_id = $3 AND current_stock >= $1
			RETURNING id, current_stock
		`
	default:
		return float64(0), nil, errors.New("invalid movement type")
	}
	var productStockID uuid.UUID
	err = tx.QueryRow(
		query,
		stockMovementReq.MovementQuantity,
		stockMovementReq.ProductID,
		stockMovementReq.WarehouseID,
	).Scan(
		&productStockID,
		&currentStock,
	)
	if err != nil {
		slog.Error("unable to update product stock", "err", err)
		return float64(0), nil, err
	}

	query = `
		INSERT INTO stock_movement (product_stock_id, user_id, product_id, warehouse_id, business_id, movement_type, movement_quantity, reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, product_stock_id, user_id, product_id, warehouse_id, business_id,movement_type, movement_quantity, reason, created_at
	`
	stockMovement := &dbmodels.StockMovement{}
	err = tx.QueryRow(
		query,
		productStockID,
		userID,
		stockMovementReq.ProductID,
		stockMovementReq.WarehouseID,
		businessID,
		stockMovementReq.MovementType,
		stockMovementReq.MovementQuantity,
		stockMovementReq.Reason,
	).Scan(
		&stockMovement.ID,
		&stockMovement.ProductStockID,
		&stockMovement.UserID,
		&stockMovement.ProductID,
		&stockMovement.WarehouseID,
		&stockMovement.BusinessID,
		&stockMovement.MovementType,
		&stockMovement.MovementQuantity,
		&stockMovement.Reason,
		&stockMovement.CreatedAt,
	)
	if err != nil {
		slog.Error("unable to insert stock movement", "err", err)
		return float64(0), nil, err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("unable to commit transaction", "err", err)
		return float64(0), nil, err
	}
	return currentStock, stockMovement, nil
}
