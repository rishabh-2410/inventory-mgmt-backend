CREATE TABLE users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    business_id uuid,
    password_hash TEXT NOT NULL,
    role TEXT CHECK (role IN ('owner', 'employee')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    FOREIGN KEY(business_id) references business(id) ON DELETE CASCADE
)