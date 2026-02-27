-- +goose Up
CREATE TABLE ticket_resolutions (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL UNIQUE REFERENCES tickets(id) ON DELETE CASCADE,
    resolved_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resolution_note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE ticket_resolutions;
