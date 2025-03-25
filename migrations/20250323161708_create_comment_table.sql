-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    idPost INTEGER NOT NULL,
    parentIdcomment INTEGER,
    text TEXT NOT NULL,
    CONSTRAINT fk_post FOREIGN KEY (idPost) REFERENCES posts(id),
    CONSTRAINT fk_parent_comment FOREIGN KEY (parentIdcomment) REFERENCES comments(id)
);

-- +goose Down
DROP TABLE IF EXISTS comments;
