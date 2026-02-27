-- +goose Up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_categories_name_lower_unique ON categories (LOWER(name));

-- +goose Down
DROP TABLE categories;