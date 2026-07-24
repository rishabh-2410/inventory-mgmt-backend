ALTER TABLE categories
ADD CONSTRAINT categories_business_name_unique
UNIQUE (business_id, name)
