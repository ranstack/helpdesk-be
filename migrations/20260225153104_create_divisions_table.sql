-- +goose Up
CREATE TABLE divisions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE
);

-- +goose Down
DROP TABLE divisions;
