-- +goose Up
CREATE TABLE restaurants(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coordinates POINT NOT NULL,
    name TEXT NOT NULL,
    open_time TIME NOT NULL,
    close_time TIME NOT NULL
);

-- +goose Down
DROP TABLE restaurants;