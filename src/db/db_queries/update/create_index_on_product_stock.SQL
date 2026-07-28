CREATE INDEX idx_product_stock_business
ON product_stock(business_id);


CREATE INDEX idx_product_stock_business_product
ON product_stock(business_id, product_id);