CREATE TABLE product_stock(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id uuid NOT NULL,
    product_id uuid NOT NULL,
    warehouse_id uuid NOT NULL,
    current_stock NUMERIC(10,2) NOT NULL CHECK ( current_stock >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE (warehouse_id, product_id),

    FOREIGN KEY(product_id) references products(id),
    FOREIGN KEY(warehouse_id) references warehouses(id)
 )