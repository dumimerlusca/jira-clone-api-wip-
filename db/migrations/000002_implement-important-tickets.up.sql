CREATE TABLE IF NOT EXISTS important_tickets(
    user_id VARCHAR(50) NOT NULL,
    ticket_id VARCHAR(50) NOT NULL,
    project_id VARCHAR(50) NOT NULL,
    CONSTRAINT duplicate_user_ticket_pair UNIQUE(user_id, ticket_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (ticket_id) REFERENCES tickets(id)
)