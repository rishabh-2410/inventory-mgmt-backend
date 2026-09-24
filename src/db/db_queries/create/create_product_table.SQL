CREATE TABLE products(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    category_id uuid NOT NULL,
    business_id uuid NOT NULL,
    selling_price NUMERIC(10, 2) NOT NULL CHECK ( selling_price >=0),
    cost_price NUMERIC(10,2) NOT NULL CHECK ( cost_price >= 0),
    image_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    sku TEXT NOT NULL,

    FOREIGN KEY ( category_id ) references categories(id),
    FOREIGN KEY ( business_id) references business(id) ON DELETE CASCADE
)