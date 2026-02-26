-- +goose Up
CREATE TABLE ticket_attachments (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    uploaded_by INT NOT NULL REFERENCES users(id),
    file_url TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attach_ticket ON ticket_attachments(ticket_id);
CREATE INDEX idx_attach_type ON ticket_attachments(type);

-- +goose Down
DROP TABLE ticket_attachments;