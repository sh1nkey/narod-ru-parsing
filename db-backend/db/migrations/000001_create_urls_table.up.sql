CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    is_empty BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_is_empty_not_true_when_description_not_empty
    CHECK (is_empty = FALSE OR description IS NULL OR description = '')
);
