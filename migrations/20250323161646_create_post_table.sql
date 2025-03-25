-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    text TEXT NOT NULL
    readOnly BOOLEAN NOT NULL DEFAULT false
);

-- +goose Down
DROP TABLE IF EXISTS posts;
