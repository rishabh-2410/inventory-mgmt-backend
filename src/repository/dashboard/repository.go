package repository

import (
	"database/sql"

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
