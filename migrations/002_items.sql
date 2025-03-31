-- +goose Up
CREATE TABLE items(
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    sizes TEXT[] NOT NULL,
    prices FLOAT[] NOT NULL,
    photo BYTEA NOT NULL
);

-- +goose Down
DROP TABLE items;
