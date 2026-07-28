package repository

import (
	"database/sql"
	"fmt"
	dbmodel "inventory-backend/src/models/db_object"
	reqmodels "inventory-backend/src/models/req_dto"
	"log/slog"

	"github.com/google/uuid"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (pr *ProductRepo) AddProduct(businessID uuid.UUID, product *reqmodels.AddProductRequestDTO) (*dbmodel.Product, error) {
	result := &dbmodel.Product{}

	err := pr.db.QueryRow(
		`INSERT INTO products(
			name,
			category_id,
			business_id,
			selling_price,
			cost_price,
			image_url,
			sku
		) 
		 VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		) RETURNING id, name, category_id, business_id, selling_price, cost_price, image_url, created_at, updated_at, sku
		`, product.Name, product.CategoryID, businessID, product.SellingPrice, product.CostPrice, product.ImageURL, product.SKU).Scan(
		&result.ID,
		&result.Name,
		&result.CategoryID,
		&result.BusinessID,
		&result.SellingPrice,
		&result.CostPrice,
		&result.ImageURL,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.SKU,
	)

	if err != nil {
		slog.Error("failed to add product", "error", err)
		return nil, err
	}

	return result, nil

}

func (pr *ProductRepo) GetProducts(businessID uuid.UUID) ([]dbmodel.Product, error) {
	result := make([]dbmodel.Product, 0)

	rows, err := pr.db.Query(
		`SELECT * FROM products
		WHERE business_id=$1	
		`, businessID)

	if err != nil {
		slog.Error("failed to fetch products", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product dbmodel.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.CategoryID,
			&product.BusinessID,
			&product.SellingPrice,
			&product.CostPrice,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.SKU,
		)

		if err != nil {
			slog.Error("failed to fetch products", "error", err)
			return nil, err
		}

		result = append(result, product)
	}

	if err := rows.Err(); err != nil {
		slog.Error("failed to fetch products", "error", err)
		return nil, err
	}

	slog.Info("fetched all products successfully for", "business_id", businessID)

	return result, nil

}

func (pr *ProductRepo) UpdateProduct(productID uuid.UUID, businessID uuid.UUID, product *reqmodels.EditProductRequestDTO) (*dbmodel.Product, error) {
	result := &dbmodel.Product{}

	err := pr.db.QueryRow(
		`UPDATE products
		 SET 
		 	name=COALESCE($1, name),
			sku=COALESCE($2, sku),
			selling_price=COALESCE($3, selling_price),
			cost_price=COALESCE($4, cost_price),
			category_id=COALESCE($5, category_id),
			image_url=COALESCE($6, image_url)
		 WHERE id=$7 AND business_id=$8
		 RETURNING id, business_id, name, sku, selling_price, cost_price, category_id, image_url, created_at, updated_at
		`, product.Name, product.SKU, product.SellingPrice, product.CostPrice, product.CategoryID, product.ImageURL, productID, businessID).Scan(
		&result.ID,
		&result.BusinessID,
		&result.Name,
		&result.SKU,
		&result.SellingPrice,
		&result.CostPrice,
		&result.CategoryID,
		&result.ImageURL,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		slog.Error("failed to update warehouse", "error", err)
		return nil, err
	}

	return result, nil

}

func (pr *ProductRepo) DeleteProduct(productID uuid.UUID, businessID uuid.UUID) error {

	result, err := pr.db.Exec(
		`DELETE FROM products
		WHERE id=$1 AND business_id=$2
		`, productID, businessID)

	if err != nil {
		slog.Error("failed to delete product", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to delete product", "error", err)
		return err
	}

	if rowsAffected == 0 {
		slog.Error("unable to delete product: no rows affected")
		return fmt.Errorf("unable to delete product: no rows affected")
	}

	return nil

}

func (pr *ProductRepo) ValidateCategory(categoryID uuid.UUID, businessID uuid.UUID) (bool, error) {
	var exists bool

	err := pr.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1 AND business_id=$2)`, categoryID, businessID).Scan(&exists)

	if err != nil {
		slog.Error("failed to check if category exists", "error", err)
		return false, err
	}

	return exists, nil
}

func (pr *ProductRepo) ListProductsWithStock(businessID uuid.UUID) ([]dbmodel.ProductWithStock, error) {
	query := `
		SELECT p.*,
			COALESCE(SUM(ps.current_stock), 0) AS total_stock
		FROM products p
		LEFT JOIN product_stock ps ON p.id = ps.product_id
		WHERE p.business_id = $1
		GROUP BY p.id
	`
	productsWithStock := make([]dbmodel.ProductWithStock, 0)
	rows, err := pr.db.Query(query, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productWithStock dbmodel.ProductWithStock
		err = rows.Scan(&productWithStock.ID, &productWithStock.Name, &productWithStock.CategoryID, &productWithStock.BusinessID, &productWithStock.SellingPrice, &productWithStock.CostPrice, &productWithStock.ImageURL, &productWithStock.CreatedAt, &productWithStock.UpdatedAt, &productWithStock.SKU, &productWithStock.TotalStock)
		if err != nil {
			return nil, err
		}
		productsWithStock = append(productsWithStock, productWithStock)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	slog.Info("listed products with stock successfully for", "business_id", businessID)
	return productsWithStock, nil
}
