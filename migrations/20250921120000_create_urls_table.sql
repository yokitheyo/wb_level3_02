-- +goose Up
CREATE TABLE IF NOT EXISTS urls (
                                    id SERIAL PRIMARY KEY,
                                    short VARCHAR(10) NOT NULL UNIQUE,
    original TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    expires_at TIMESTAMP NOT NULL,
    visits INT DEFAULT 0,
    is_disabled BOOLEAN DEFAULT FALSE
    );

-- +goose Down
DROP TABLE IF EXISTS urls;
