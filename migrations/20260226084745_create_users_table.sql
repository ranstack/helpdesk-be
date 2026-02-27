-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar_url TEXT DEFAULT NULL,
    phone VARCHAR(15) DEFAULT NULL,
    role VARCHAR(10) NOT NULL,
    division_id INT NOT NULL REFERENCES divisions(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_email_lower_unique ON users (LOWER(email));
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_division ON users(division_id);

-- +goose Down
DROP TABLE users;
