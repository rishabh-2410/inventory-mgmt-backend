CREATE TABLE stock_movement(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    product_stock_id uuid NOT NULL,
    user_id uuid NOT NULL,
    product_id uuid NOT NULL,
    warehouse_id uuid NOT NULL,
    business_id uuid NOT NULL,
    movement_type TEXT NOT NULL CHECK ( movement_type IN ('SALE', 'RETURN', 'DAMAGE', 'RECEIVE')),
    movement_quantity NUMERIC(10, 2) NOT NULL CHECK ( movement_quantity > 0),
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,


    FOREIGN KEY(product_stock_id) references product_stock(id),
    FOREIGN KEY(user_id) references users(id),
    FOREIGN KEY(warehouse_id) references warehouses(id),
    FOREIGN KEY(product_id) references products(id),
    FOREIGN KEY(business_id) references business(id)
)