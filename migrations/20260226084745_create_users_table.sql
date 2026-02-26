-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(50),
    email VARCHAR(20) UNIQUE,
    password VARCHAR(50),
    avatar_url TEXT DEFAULT NULL,
    phone VARCHAR(15) DEFAULT NULL,
    role VARCHAR(10), --STAFF | IT | ADMIN
    division_id INTEGER,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_division ON users(division_id);

-- +goose Down
DROP TABLE users;
