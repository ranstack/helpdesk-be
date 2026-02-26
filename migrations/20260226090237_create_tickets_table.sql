-- +goose Up
CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories(id),
    priority VARCHAR(50) NOT NULL, --LOW|MEDIUM|URGENT
    status VARCHAR(20) NOT NULL, --OPEN|INPROGRESS|RESOLVED|CLOSED
    created_by INTEGER NOT NULL REFERENCES users(id),
    assigned_to INTEGER DEFAULT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP DEFAULT NULL,
    closed_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_tickets_title ON tickets(title);
CREATE INDEX idx_tickets_category ON tickets(category_id);
CREATE INDEX idx_tickets_priority ON tickets(priority);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_created_by ON tickets(created_by);
CREATE INDEX idx_tickets_assign_to ON tickets(assigned_to);

-- +goose Down
DROP TABLE tickets;
