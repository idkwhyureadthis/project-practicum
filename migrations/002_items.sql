-- +goose Up
CREATE TABLE items(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    sizes TEXT[] NOT NULL,
    prices FLOAT[] NOT NULL,
    photos TEXT[] NOT NULL
);

-- +goose Down
DROP TABLE items;
