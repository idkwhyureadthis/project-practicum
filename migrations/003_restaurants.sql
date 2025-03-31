-- +goose Up
CREATE TABLE restaurants(
    id UUID PRIMARY KEY,
    coordinates POINT NOT NULL,
    address TEXT NOT NULL,
    open_time TIME NOT NULL,
    close_time TIME NOT NULL
);

-- +goose Down
DROP TABLE restaurants;