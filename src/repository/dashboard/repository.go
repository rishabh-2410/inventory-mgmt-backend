package repository

import (
	"database/sql"
	"log/slog"

	resmodels "inventory-backend/src/models/res_dto"

	"github.com/google/uuid"
)

type DashboardRepo struct {
	db *sql.DB
}

func NewDashboardRepo(db *sql.DB) *DashboardRepo {
	return &DashboardRepo{db: db}
}

func (d *DashboardRepo) GetDashboardData(businessID uuid.UUID) (*resmodels.DashboardData, error) {
	result := &resmodels.DashboardData{}
	query := `

		WITH product_totals AS (
		    SELECT
			    product_id,
				SUM(current_stock) AS total_stock
			FROM product_stock
			WHERE business_id = $1
			GROUP BY product_id
		)
		SELECT 
			COALESCE(SUM(total_stock), 0) AS current_stock_value,

			COUNT(DISTINCT CASE
				WHEN total_stock BETWEEN 1 AND 9 THEN product_id
			END) AS low_stock_products,


			COUNT(DISTINCT CASE
				WHEN total_stock = 0 THEN product_id
			END) AS out_of_stock_products
		FROM product_totals
	`
	err := d.db.QueryRow(query, businessID).Scan(
		&result.CurrentInventoryValue,
		&result.LowStockProducts,
		&result.OutOfStockProducts,
	)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (d *DashboardRepo) GetInventoryStats(businessID uuid.UUID) (*resmodels.InventoryStats, error) {
	slog.Info("dashboard repository get inventory stats", "business_id", businessID)
	stats := &resmodels.InventoryStats{}

	query := `
	SELECT
		COALESCE(SUM(movement_quantity) FILTER (
			WHERE movement_type = 'SALE'
			AND created_at >= NOW() - INTERVAL '30 days'
		), 0) AS units_sold,


		COALESCE(SUM(movement_quantity) FILTER(
			WHERE movement_type = 'RECEIVE'
			AND created_at >= NOW() - INTERVAL '30 days'
		), 0) AS units_received,

		COALESCE(SUM(movement_quantity) FILTER (
			WHERE movement_type = 'DAMAGE'
			AND created_at >= NOW() - INTERVAL '30 days'
		), 0) AS units_damaged,

		COALESCE(SUM(movement_quantity) FILTER (
			WHERE movement_type = 'RETURN'
			AND created_at >= NOW() - INTERVAL '30 days'
		), 0) AS units_returned

	FROM stock_movement
	WHERE business_id = $1
	`

	err := d.db.QueryRow(
		query,
		businessID,
	).Scan(
		&stats.UnitsSold,
		&stats.UnitsReceived,
		&stats.UnitsDamaged,
		&stats.UnitsReturned,
	)

	if err != nil {
		slog.Error("dashboard repository get inventory stats failed business_id=%s err=%v", businessID, err)
		return nil, err
	}

	slog.Info("dashboard repository get inventory stats success business_id=%s units_sold=%d units_received=%d units_damaged=%d units_returned=%d", businessID, stats.UnitsSold, stats.UnitsReceived, stats.UnitsDamaged, stats.UnitsReturned)
	return stats, nil
}
