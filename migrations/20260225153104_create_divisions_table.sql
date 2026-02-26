-- +goose Up
CREATE TABLE divisions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE divisions;
