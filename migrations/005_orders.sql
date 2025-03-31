-- +goose Up
CREATE TABLE orders(
    id UUID PRIMARY KEY,
    displayed_id INTEGER NOT NULL,
    restaurant_id UUID REFERENCES restaurants(id),
    total_price FLOAT,
    status TEXT
);


-- +goose Down
DROP TABLE orders;