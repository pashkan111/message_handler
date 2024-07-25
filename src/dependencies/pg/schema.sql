CREATE TABLE message (
    id SERIAL PRIMARY KEY,
    content JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_processed BOOLEAN NOT NULL DEFAULT false
);
