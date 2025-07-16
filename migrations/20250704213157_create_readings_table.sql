-- +goose Up
CREATE TABLE readings (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    finished BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE readings;